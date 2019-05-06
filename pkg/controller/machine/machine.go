package machine

import (
	"context"
	machineinternal "github.com/adracus/kubeception/pkg/internal/controller/machine"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api/pkg/controller/machine"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewActuatorWithDeps instantiates a new machine actuator with the dependencies usually obtained via dependency injection.
// TODO: Remove this as soon as the cluster API supports proper dependency injection for its actuators.
func NewActuatorWithDeps(ctx context.Context, client client.Client, scheme *runtime.Scheme) machine.Actuator {
	return machineinternal.NewActuatorWithDeps(ctx, client, scheme)
}
