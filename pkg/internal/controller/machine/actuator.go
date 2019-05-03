package machine

import (
	"context"
	"fmt"
	"github.com/adracus/kubeception/pkg/apis/kubeception/v1alpha1"
	"github.com/adracus/kubeception/pkg/internal/controller"
	kubeceptioncluster "github.com/adracus/kubeception/pkg/internal/controller/cluster"
	"github.com/adracus/kubeception/pkg/internal/util"
	"github.com/adracus/kubeception/pkg/internal/util/pointers"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clusterv1alpha1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	"sigs.k8s.io/cluster-api/pkg/controller/machine"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
)

type actuator struct {
	ctx    context.Context
	client client.Client
	scheme *runtime.Scheme
}

func (a *actuator) InjectStopChannel(stopChan <-chan struct{}) error {
	a.ctx = util.ContextFromStopChannel(stopChan)
	return nil
}

func (a *actuator) InjectClient(client client.Client) error {
	a.client = client
	return nil
}

func (a *actuator) InjectScheme(scheme *runtime.Scheme) error {
	a.scheme = scheme
	return nil
}

func (a *actuator) InjectFunc(f inject.Func) error {
	return nil
}

func NewActuatorWithDeps(ctx context.Context, client client.Client, scheme *runtime.Scheme) machine.Actuator {
	return &actuator{
		ctx,
		client,
		scheme,
	}
}

func configsFromObjects(cluster *clusterv1alpha1.Cluster, machine *clusterv1alpha1.Machine) (*v1alpha1.ClusterConfig, *v1alpha1.MachineConfig, error) {
	clusterConfig, err := kubeceptioncluster.ConfigFromCluster(cluster)
	if err != nil {
		return nil, nil, err
	}

	machineConfig, err := ConfigFromMachine(machine)
	if err != nil {
		return nil, nil, err
	}

	return clusterConfig, machineConfig, nil
}

func mkMachineStatefulSet(machine *clusterv1alpha1.Machine) *appsv1.StatefulSet {
	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: machine.Namespace,
			Name:      machine.Name,
		},
	}
}

func (a *actuator) Create(ctx context.Context, cluster *clusterv1alpha1.Cluster, machine *clusterv1alpha1.Machine) error {
	config, _, err := configsFromObjects(cluster, machine)
	if err != nil {
		return err
	}

	labels := map[string]string{
		controller.MachineLabel: machine.Name,
	}

	statefulSet := mkMachineStatefulSet(machine)
	if _, err := controllerruntime.CreateOrUpdate(ctx, a.client, statefulSet, func(runtime.Object) error {
		statefulSet.Labels = labels
		statefulSet.Spec = appsv1.StatefulSetSpec{
			Replicas: pointers.Int32(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					AutomountServiceAccountToken: pointers.Bool(false),
					Containers: []corev1.Container{
						{
							Name:  "kubelet",
							Image: fmt.Sprintf("adracus/dind-kubelet:%s", config.KubernetesVersion),
							Env: []corev1.EnvVar{
								{Name: "KUBECONFIG", Value: "/etc/kubeconfig/kubeconfig"},
								{Name: "ADDITIONAL_DOCKERD_ARGS", Value: "--storage-driver=vfs"},
								{Name: "POD_IP", ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "status.podIP"}}},
							},
							Command: []string{
								"/entrypoint.sh",
								"--kubeconfig=/etc/kubeconfig/kubeconfig",
								"--fail-swap-on=false",
								"--port=20250",
								"--containerized",
								"--feature-gates=LocalStorageCapacityIsolation=false",
								"--cloud-provider=",
								"--hostname-override=$(POD_IP)",
							},
							SecurityContext: &corev1.SecurityContext{
								Privileged: pointers.Bool(true),
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "kubeconfig",
									MountPath: "/etc/kubeconfig",
								},
								{
									Name:      "rootfs",
									MountPath: "/rootfs",
									ReadOnly:  true,
								},
							},
							Ports: []corev1.ContainerPort{
								{
									HostPort:      20250,
									ContainerPort: 20250,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "kubeconfig",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: kubeceptioncluster.KubeconfigSecretName,
								},
							},
						},
						{
							Name: "rootfs",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/",
								},
							},
						},
					},
				},
			},
		}

		return controllerruntime.SetControllerReference(machine, statefulSet, a.scheme)
	}); err != nil {
		return err
	}

	return nil
}

func (a *actuator) Delete(ctx context.Context, cluster *clusterv1alpha1.Cluster, machine *clusterv1alpha1.Machine) error {
	return nil
}

func (a *actuator) Update(ctx context.Context, cluster *clusterv1alpha1.Cluster, machine *clusterv1alpha1.Machine) error {
	return a.Create(ctx, cluster, machine)
}

func (a *actuator) Exists(ctx context.Context, cluster *clusterv1alpha1.Cluster, machine *clusterv1alpha1.Machine) (bool, error) {
	_, _, err := configsFromObjects(cluster, machine)
	if err != nil {
		return false, err
	}

	statefulSet := mkMachineStatefulSet(machine)
	if err := a.client.Get(ctx, client.ObjectKey{Namespace: machine.Namespace, Name: machine.Name}, statefulSet); err != nil {
		if apierrors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
