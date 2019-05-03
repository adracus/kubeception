package machine

import (
	"context"
	"sigs.k8s.io/cluster-api/pkg/controller/machine"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func AddToManager(mgr manager.Manager) error {
	ctx := context.TODO()
	return machine.AddWithActuator(mgr, NewActuatorWithDeps(ctx, mgr.GetClient(), mgr.GetScheme()))
}
