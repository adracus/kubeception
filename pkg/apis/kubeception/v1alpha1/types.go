package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeception.cloud/kubeception/pkg/util"
)

var (
	ClusterConfigKind = util.MustTypeToKind(&ClusterConfig{})

	MachineConfigKind = util.MustTypeToKind(&MachineConfig{})
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterConfig is the kubeception cluster configuration.
type ClusterConfig struct {
	metav1.TypeMeta `json:",inline"`

	ControlPlane      ControlPlane `json:"controlPlane"`
	KubernetesVersion string       `json:"kubernetesVersion"`
}

// ControlPlane is the specification of a cluster control plane.
type ControlPlane struct {
	ETCD              ETCD               `json:"etcd"`
	APIServer         APIServer          `json:"apiServer"`
	ControllerManager *ControllerManager `json:"controllerManager,omitempty"`
	Scheduler         *Scheduler         `json:"scheduler,omitempty"`
}

// ETCD carries etcd configuration.
type ETCD struct {
}

// APIServer carries Kubernetes API server configuration.
type APIServer struct {
}

// ControllerManager carries Kubernetes controller manager configuration.
type ControllerManager struct {
}

// ControllerManager carries Kubernetes scheduler configuration.
type Scheduler struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterConfig is the kubeception machine configuration.
type MachineConfig struct {
	metav1.TypeMeta `json:",inline"`
}
