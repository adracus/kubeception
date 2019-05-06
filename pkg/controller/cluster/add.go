package cluster

import (
	"context"
	"sigs.k8s.io/cluster-api/pkg/controller/cluster"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// AddToManager adds the cluster controller with the kubeception actuator to the cluster.
func AddToManager(mgr manager.Manager) error {
	ctx := context.TODO()
	return cluster.AddWithActuator(mgr, NewActuatorWithDeps(ctx, mgr.GetClient(), mgr.GetScheme()))
}
