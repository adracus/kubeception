package cluster

import (
	"context"
	"fmt"
	"github.com/adracus/kubeception/pkg/apis/kubeception/v1alpha1"
	"github.com/adracus/kubeception/pkg/internal/controller"
	"github.com/adracus/kubeception/pkg/internal/util"
	"github.com/adracus/kubeception/pkg/internal/util/pointers"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	clusterv1alpha1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	"sigs.k8s.io/cluster-api/pkg/controller/cluster"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ETCDClientPort      = 2379
	ETCDServiceName     = "etcd"
	ETCDStatefulSetName = ETCDServiceName
)

func NewActuator() cluster.Actuator {
	return &actuator{}
}

func NewActuatorWithDeps(ctx context.Context, client client.Client, scheme *runtime.Scheme) cluster.Actuator {
	return &actuator{
		ctx:    ctx,
		client: client,
		scheme: scheme,
	}
}

type actuator struct {
	ctx    context.Context
	client client.Client
	scheme *runtime.Scheme
}

func (a *actuator) InjectClient(c client.Client) error {
	a.client = c
	return nil
}

func (a *actuator) InjectStopChannel(stopCh <-chan struct{}) error {
	a.ctx = util.ContextFromStopChannel(stopCh)
	return nil
}

func (a *actuator) InjectScheme(scheme *runtime.Scheme) error {
	a.scheme = scheme
	return nil
}

func (a *actuator) Reconcile(cluster *clusterv1alpha1.Cluster) error {
	config, err := ConfigFromCluster(cluster)
	if err != nil {
		return err
	}

	return a.reconcile(a.ctx, cluster, config)
}

func (a *actuator) reconcile(ctx context.Context, cluster *clusterv1alpha1.Cluster, config *v1alpha1.ClusterConfig) error {
	if err := a.reconcileETCD(ctx, cluster, &config.ControlPlane.ETCD); err != nil {
		return err
	}

	if err := a.reconcileAPIServer(ctx, cluster, &config.ControlPlane.APIServer); err != nil {
		return err
	}

	return nil
}

func (a *actuator) reconcileETCD(ctx context.Context, cluster *clusterv1alpha1.Cluster, etcd *v1alpha1.ETCD) error {
	etcdLabels := map[string]string{
		controller.ControlPlaneComponentLabel: controller.ETCDComponent,
	}

	etcdService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.Namespace,
			Name:      ETCDServiceName,
		},
	}
	if _, err := controllerruntime.CreateOrUpdate(ctx, a.client, etcdService, func(runtime.Object) error {
		etcdService.Spec.Selector = etcdLabels
		etcdService.Spec.Type = corev1.ServiceTypeClusterIP
		etcdService.Spec.Ports = []corev1.ServicePort{
			{
				Port:       ETCDClientPort,
				TargetPort: intstr.FromInt(ETCDClientPort),
			},
		}
		util.SetMetaDataLabels(etcdService, etcdLabels)
		return controllerruntime.SetControllerReference(cluster, etcdService, a.scheme)
	}); err != nil {
		return err
	}

	etcdStatefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.Namespace,
			Name:      ETCDStatefulSetName,
		},
	}
	_, err := controllerruntime.CreateOrUpdate(ctx, a.client, etcdStatefulSet, func(existing runtime.Object) error {
		util.SetMetaDataLabels(etcdStatefulSet, etcdLabels)
		etcdStatefulSet.Spec = appsv1.StatefulSetSpec{
			Replicas: pointers.Int32(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: etcdLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: etcdLabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "etcd",
							Image: "quay.io/coreos/etcd:v3.3.12",
							Command: []string{
								"etcd",
								fmt.Sprintf("--advertise-client-urls=http://%s:%d", ETCDServiceName, ETCDClientPort),
								fmt.Sprintf("--listen-client-urls=http://0.0.0.0:%d", ETCDClientPort),
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: ETCDClientPort,
									Name:          "etcd",
								},
							},
						},
					},
				},
			},
		}
		return controllerruntime.SetControllerReference(cluster, etcdStatefulSet, a.scheme)
	})
	return err
}

func (a *actuator) reconcileAPIServer(ctx context.Context, cluster *clusterv1alpha1.Cluster, apiServer *v1alpha1.APIServer) error {
	apiServerLabels := map[string]string{
		controller.ControlPlaneComponentLabel: controller.APIServerComponent,
	}

	apiServerDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.Namespace,
			Name:      "kube-apiserver",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointers.Int32(1),
		},
	}
	if _, err := controllerruntime.CreateOrUpdate(ctx, a.client, apiServerDeployment, func(runtime.Object) error {
		apiServerDeployment.Spec = appsv1.DeploymentSpec{
			Replicas: pointers.Int32(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: apiServerLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: apiServerLabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "kube-apiserver",
							Image: "googlecontainer/kube-apiserver:v1.14.1",
							Command: []string{
								"kube-apiserver",
								fmt.Sprintf("--etcd-servers=http://%s:%d", ETCDServiceName, ETCDClientPort),
							},
						},
					},
				},
			},
		}
		return controllerruntime.SetControllerReference(cluster, apiServerDeployment, a.scheme)
	}); err != nil {
		return err
	}

	return nil
}

func (a *actuator) Delete(cluster *clusterv1alpha1.Cluster) error {
	config, err := ConfigFromCluster(cluster)
	if err != nil {
		return err
	}

	return a.delete(a.ctx, cluster, config)
}

func (a *actuator) delete(ctx context.Context, cluster *clusterv1alpha1.Cluster, config *v1alpha1.ClusterConfig) error {
	return nil
}
