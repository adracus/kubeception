package controller

import (
	"github.com/adracus/kubeception/pkg/controller/cluster"
	"github.com/adracus/kubeception/pkg/controller/machine"
	"github.com/adracus/kubeception/pkg/internal/util"
)

var (
	addToManagerBuilder = util.NewAddToManagerBuilder(
		cluster.AddToManager,
		machine.AddToManager,
	)

	// AddToManager adds all kubeception controllers to the given manager.
	AddToManager = addToManagerBuilder.AddToManager
)
