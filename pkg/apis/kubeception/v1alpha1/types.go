package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterConfig is the kubeception cluster configuration.
type ClusterConfig struct {
	metav1.TypeMeta `json:",inline"`

	ControlPlane ControlPlane `json:"controlPlane"`
}

// ControlPlane is the specification of a cluster control plane.
type ControlPlane struct {
	ETCD                  ETCD                  `json:"etcd"`
	APIServer             APIServer             `json:"apiServer"`
	KubeControllerManager KubeControllerManager `json:"kubeControllerManager"`
}

// ETCD carries etcd configuration.
type ETCD struct {
}

// APIServer carries Kubernetes API server configuration.
type APIServer struct {
}

// KubeControllerManager carries Kubernetes controller manager configuration.
type KubeControllerManager struct {
}
