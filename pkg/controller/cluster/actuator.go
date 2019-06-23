package cluster

import (
	"context"
	"fmt"

	"kubeception.cloud/kubeception/pkg/controller/common"

	"kubeception.cloud/kubeception/pkg/util/controller"

	"kubeception.cloud/kubeception/pkg/util"

	"kubeception.cloud/kubeception/pkg/apis/kubeception/helper"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"kubeception.cloud/kubeception/pkg/apis/kubeception/v1alpha1"
	"kubeception.cloud/kubeception/pkg/util/pointers"
	clusterv1alpha1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	"sigs.k8s.io/cluster-api/pkg/controller/cluster"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ETCDClientPort      = 2379
	ETCDServiceName     = "etcd"
	ETCDStatefulSetName = ETCDServiceName

	APIServerPort           = 443
	APIServerDeploymentName = "apiserver"
	APIServerServiceName    = APIServerDeploymentName

	KubeconfigSecretName = "kubeconfig"

	ControllerManagerDeploymentName = "controller-manager"

	SchedulerDeploymentName = "scheduler"
)

// NewActuatorWithDeps instantiates a new actuator with the dependencies that are usually injected.
// TODO: Remove this constructor as soon as the cluster api supports proper injection on the actuators.
func NewActuatorWithDeps(ctx context.Context, c client.Client, scheme *runtime.Scheme) cluster.Actuator {
	return &actuator{
		WithContext: controller.NewWithContext(ctx),
		WithClient:  controller.NewWithClient(c),
		WithScheme:  controller.NewWithScheme(scheme),
	}
}

type actuator struct {
	controller.WithContext
	controller.WithClient
	controller.WithScheme
}

func (a *actuator) Reconcile(cluster *clusterv1alpha1.Cluster) error {
	config, err := helper.LoadClusterConfig(cluster.Spec.ProviderSpec.Value.Raw)
	if err != nil {
		return err
	}

	return a.reconcile(a.Context, cluster, config)
}

func (a *actuator) reconcile(ctx context.Context, cluster *clusterv1alpha1.Cluster, config *v1alpha1.ClusterConfig) error {
	if err := a.reconcileETCD(ctx, cluster, &config.ControlPlane.ETCD); err != nil {
		return err
	}

	if err := a.reconcileAPIServer(ctx, cluster, config, &config.ControlPlane.APIServer); err != nil {
		return err
	}

	if config.ControlPlane.ControllerManager != nil {
		if err := a.reconcileControllerManager(ctx, cluster, config, config.ControlPlane.ControllerManager); err != nil {
			return err
		}
	} else {
		if err := a.deleteControllerManager(ctx, cluster); err != nil {
			return err
		}
	}

	if config.ControlPlane.Scheduler != nil {
		if err := a.reconcileScheduler(ctx, cluster, config, config.ControlPlane.Scheduler); err != nil {
			return err
		}
	} else {
		if err := a.deleteScheduler(ctx, cluster); err != nil {
			return err
		}
	}

	return nil
}

func (a *actuator) reconcileETCD(ctx context.Context, cluster *clusterv1alpha1.Cluster, etcd *v1alpha1.ETCD) error {
	etcdService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.Namespace,
			Name:      ETCDServiceName,
		},
	}
	if _, err := controllerruntime.CreateOrUpdate(ctx, a.Client, etcdService, func() error {
		etcdService.Spec.Selector = ETCDLabels
		etcdService.Spec.Type = corev1.ServiceTypeClusterIP
		etcdService.Spec.Ports = []corev1.ServicePort{
			{
				Port:       ETCDClientPort,
				TargetPort: intstr.FromInt(ETCDClientPort),
			},
		}
		util.SetMetaDataLabels(etcdService, ETCDLabels)
		return controllerruntime.SetControllerReference(cluster, etcdService, a.Scheme)
	}); err != nil {
		return err
	}

	etcdStatefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.Namespace,
			Name:      ETCDStatefulSetName,
		},
	}
	_, err := controllerruntime.CreateOrUpdate(ctx, a.Client, etcdStatefulSet, func() error {
		util.SetMetaDataLabels(etcdStatefulSet, ETCDLabels)
		etcdStatefulSet.Spec = appsv1.StatefulSetSpec{
			Replicas: pointers.Int32(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: ETCDLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ETCDLabels,
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
		return controllerruntime.SetControllerReference(cluster, etcdStatefulSet, a.Scheme)
	})
	return err
}

func (a *actuator) reconcileAPIServer(ctx context.Context, cluster *clusterv1alpha1.Cluster, config *v1alpha1.ClusterConfig, apiServer *v1alpha1.APIServer) error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.Namespace,
			Name:      APIServerServiceName,
		},
	}
	if _, err := controllerruntime.CreateOrUpdate(ctx, a.Client, service, func() error {
		util.SetMetaDataLabels(service, APIServerLabels)
		service.Spec.Type = corev1.ServiceTypeNodePort
		service.Spec.Selector = APIServerLabels
		service.Spec.Ports = []corev1.ServicePort{
			{
				Port:       APIServerPort,
				TargetPort: intstr.FromInt(APIServerPort),
			},
		}
		return nil
	}); err != nil {
		return err
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.Namespace,
			Name:      APIServerDeploymentName,
		},
	}
	if _, err := controllerruntime.CreateOrUpdate(ctx, a.Client, deployment, func() error {
		util.SetMetaDataLabels(deployment, APIServerLabels)
		deployment.Spec = appsv1.DeploymentSpec{
			Replicas: pointers.Int32(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: APIServerLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: APIServerLabels,
					Annotations: map[string]string{
						"basic-auth": `kubeception,kubeception,kubeception,"system:masters"`,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "kube-apiserver",
							Image: common.HyperkubeImageForConfig(config),
							Command: []string{
								"/hyperkube",
								"apiserver",
								fmt.Sprintf("--etcd-servers=http://%s:%d", ETCDServiceName, ETCDClientPort),
								fmt.Sprintf("--secure-port=%d", APIServerPort),
								"--basic-auth-file=/etc/basic-auth/basic-auth",
								"--authorization-mode=AlwaysAllow,RBAC,Node",
								"--disable-admission-plugins=ServiceAccount",
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "apiserver",
									ContainerPort: APIServerPort,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "basic-auth",
									MountPath: "/etc/basic-auth",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "basic-auth",
							VolumeSource: corev1.VolumeSource{
								DownwardAPI: &corev1.DownwardAPIVolumeSource{
									Items: []corev1.DownwardAPIVolumeFile{
										{Path: "basic-auth", FieldRef: &corev1.ObjectFieldSelector{FieldPath: "metadata.annotations['basic-auth']"}},
									},
								},
							},
						},
					},
				},
			},
		}
		return controllerruntime.SetControllerReference(cluster, deployment, a.Scheme)
	}); err != nil {
		return err
	}

	kubeConfigSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.Namespace,
			Name:      KubeconfigSecretName,
		},
	}
	if _, err := controllerruntime.CreateOrUpdate(ctx, a.Client, kubeConfigSecret, func() error {
		if err := UpdateKubeconfigSecret(kubeConfigSecret, &clientcmdapi.Config{
			APIVersion:  "v1",
			Kind:        "Config",
			Preferences: clientcmdapi.Preferences{},
			Clusters: map[string]*clientcmdapi.Cluster{
				"kubeception": {
					Server:                fmt.Sprintf("https://%s:%d", APIServerServiceName, APIServerPort),
					InsecureSkipTLSVerify: true,
				},
			},
			Contexts: map[string]*clientcmdapi.Context{
				"kubeception": {
					Cluster:  "kubeception",
					AuthInfo: "kubeception",
				},
			},
			CurrentContext: "kubeception",
			AuthInfos: map[string]*clientcmdapi.AuthInfo{
				"kubeception": {
					Username: "kubeception",
					Password: "kubeception",
				},
			},
		}); err != nil {
			return err
		}

		return controllerruntime.SetControllerReference(cluster, kubeConfigSecret, a.Scheme)
	}); err != nil {
		return err
	}

	return nil
}

func (a *actuator) reconcileControllerManager(ctx context.Context, cluster *clusterv1alpha1.Cluster, config *v1alpha1.ClusterConfig, controllerManager *v1alpha1.ControllerManager) error {
	deployment := &appsv1.Deployment{
		ObjectMeta: util.ObjectMeta(cluster.Namespace, ControllerManagerDeploymentName),
	}
	if _, err := controllerruntime.CreateOrUpdate(ctx, a.Client, deployment, func() error {
		util.SetMetaDataLabels(deployment, ControllerManagerLabels)
		deployment.Spec = appsv1.DeploymentSpec{
			Replicas: pointers.Int32(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: ControllerManagerLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ControllerManagerLabels,
				},
				Spec: corev1.PodSpec{
					AutomountServiceAccountToken: pointers.Bool(false),
					Containers: []corev1.Container{
						{
							Name:  "controller-manager",
							Image: common.HyperkubeImageForConfig(config),
							Command: []string{
								"/hyperkube",
								"controller-manager",
								fmt.Sprintf("--service-cluster-ip-range=%s", cluster.Spec.ClusterNetwork.Services.CIDRBlocks[0]),
								fmt.Sprintf("--cluster-cidr=%s", cluster.Spec.ClusterNetwork.Pods.CIDRBlocks[0]),
								"--kubeconfig=/etc/kubeconfig/kubeconfig",
								"--allocate-node-cidrs=true",
								"--cluster-name=kubeception",
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "kubeconfig",
									MountPath: "/etc/kubeconfig",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "kubeconfig",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: KubeconfigSecretName,
								},
							},
						},
					},
				},
			},
		}

		return controllerruntime.SetControllerReference(cluster, deployment, a.Scheme)
	}); err != nil {
		return err
	}

	return nil
}

func (a *actuator) deleteControllerManager(ctx context.Context, cluster *clusterv1alpha1.Cluster) error {
	controllerManagerDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.Namespace,
			Name:      ControllerManagerDeploymentName,
		},
	}
	return client.IgnoreNotFound(a.Client.Delete(ctx, controllerManagerDeployment))
}

func (a *actuator) reconcileScheduler(ctx context.Context, cluster *clusterv1alpha1.Cluster, config *v1alpha1.ClusterConfig, scheduler *v1alpha1.Scheduler) error {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.Namespace,
			Name:      SchedulerDeploymentName,
		},
	}
	if _, err := controllerruntime.CreateOrUpdate(ctx, a.Client, deployment, func() error {
		util.SetMetaDataLabels(deployment, SchedulerLabels)
		deployment.Spec = appsv1.DeploymentSpec{
			Replicas: pointers.Int32(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: SchedulerLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: SchedulerLabels,
				},
				Spec: corev1.PodSpec{
					AutomountServiceAccountToken: pointers.Bool(false),
					Containers: []corev1.Container{
						{
							Name:  "scheduler",
							Image: common.HyperkubeImageForConfig(config),
							Command: []string{
								"/hyperkube",
								"scheduler",
								"--kubeconfig=/etc/kubeconfig/kubeconfig",
								"--authorization-kubeconfig=/etc/kubeconfig/kubeconfig",
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "kubeconfig",
									MountPath: "/etc/kubeconfig",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "kubeconfig",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: KubeconfigSecretName,
								},
							},
						},
					},
				},
			},
		}

		return controllerruntime.SetControllerReference(cluster, deployment, a.Scheme)
	}); err != nil {
		return err
	}

	return nil
}

func (a *actuator) deleteScheduler(ctx context.Context, cluster *clusterv1alpha1.Cluster) error {
	schedulerDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.Namespace,
			Name:      SchedulerDeploymentName,
		},
	}
	return client.IgnoreNotFound(a.Client.Delete(ctx, schedulerDeployment))
}

func (a *actuator) Delete(cluster *clusterv1alpha1.Cluster) error {
	config, err := helper.LoadClusterConfig(cluster.Spec.ProviderSpec.Value.Raw)
	if err != nil {
		return err
	}

	return a.delete(a.Context, cluster, config)
}

func (a *actuator) delete(ctx context.Context, cluster *clusterv1alpha1.Cluster, config *v1alpha1.ClusterConfig) error {
	return nil
}
