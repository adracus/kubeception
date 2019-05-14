package cluster

import (
	"context"
	clusterinternal "kubeception.cloud/kubeception/pkg/internal/controller/cluster"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api/pkg/controller/cluster"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewActuatorWithDeps instantiates a new actuator with the dependencies that are usually injected.
// TODO: Remove this constructor as soon as the cluster api supports proper injection on the actuators.
func NewActuatorWithDeps(ctx context.Context, client client.Client, scheme *runtime.Scheme) cluster.Actuator {
	return clusterinternal.NewActuatorWithDeps(ctx, client, scheme)
}
