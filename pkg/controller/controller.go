package controller

import (
	"kubeception.cloud/kubeception/pkg/controller/certificate"
	"kubeception.cloud/kubeception/pkg/controller/cluster"
	"kubeception.cloud/kubeception/pkg/controller/machine"
	"kubeception.cloud/kubeception/pkg/util"
)

var (
	addToManagerBuilder = util.NewAddToManagerBuilder(
		cluster.AddToManager,
		machine.AddToManager,
		certificate.AddToManager,
	)

	// AddToManager adds all kubeception controllers to the given manager.
	AddToManager = addToManagerBuilder.AddToManager
)
