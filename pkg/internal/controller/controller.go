package controller

import "fmt"

const (
	LabelPrefix = "kubeception.io"
)

var (
	ControlPlaneComponentLabel = fmt.Sprintf("%s/control-plane-component", LabelPrefix)

	ETCDComponent              = "etcd"
	APIServerComponent         = "kube-apiserver"
	ControllerManagerComponent = "controller-manager"
	SchedulerComponent         = "scheduler"

	MachineLabel = fmt.Sprintf("%s/machine", LabelPrefix)
)
