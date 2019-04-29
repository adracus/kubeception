package cluster

import (
	"context"
	clusterinternal "github.com/adracus/kubeception/pkg/internal/controller/cluster"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api/pkg/controller/cluster"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewActuator() cluster.Actuator {
	return clusterinternal.NewActuator()
}

func NewActuatorWithDeps(ctx context.Context, client client.Client, scheme *runtime.Scheme) cluster.Actuator {
	return clusterinternal.NewActuatorWithDeps(ctx, client, scheme)
}
