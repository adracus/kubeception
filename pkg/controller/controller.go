package controller

import (
	"kubeception.cloud/kubeception/pkg/controller/cluster"
	"kubeception.cloud/kubeception/pkg/controller/machine"
	"kubeception.cloud/kubeception/pkg/internal/util"
)

var (
	addToManagerBuilder = util.NewAddToManagerBuilder(
		cluster.AddToManager,
		machine.AddToManager,
	)

	// AddToManager adds all kubeception controllers to the given manager.
	AddToManager = addToManagerBuilder.AddToManager
)
