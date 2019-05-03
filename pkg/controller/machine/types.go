package machine

import (
	"context"
	"github.com/adracus/kubeception/pkg/apis/kubeception/v1alpha1"
	clusterv1alpha1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

type Actuator interface {
	Create(context.Context, *clusterv1alpha1.Cluster, *v1alpha1.ClusterConfig, *clusterv1alpha1.Machine, *v1alpha1.MachineConfig) error
	Delete(context.Context, *clusterv1alpha1.Cluster, *v1alpha1.ClusterConfig, *clusterv1alpha1.Machine, *v1alpha1.MachineConfig) error
	Update(context.Context, *clusterv1alpha1.Cluster, *v1alpha1.ClusterConfig, *clusterv1alpha1.Machine, *v1alpha1.MachineConfig) error
	Exists(context.Context, *clusterv1alpha1.Cluster, *v1alpha1.ClusterConfig, *clusterv1alpha1.Machine, *v1alpha1.MachineConfig) (bool, error)
}
