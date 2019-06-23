package cluster

import (
	"fmt"

	"kubeception.cloud/kubeception/pkg/controller/common"
)

var (
	ControlPlaneComponentLabel = fmt.Sprintf("%s/control-plane-component", common.LabelPrefix)

	ETCDComponent              = "etcd"
	APIServerComponent         = "kube-apiserver"
	ControllerManagerComponent = "controller-manager"
	SchedulerComponent         = "scheduler"

	ETCDLabels = map[string]string{
		ControlPlaneComponentLabel: ETCDComponent,
	}

	APIServerLabels = map[string]string{
		ControlPlaneComponentLabel: APIServerComponent,
	}

	ControllerManagerLabels = map[string]string{
		ControlPlaneComponentLabel: ControllerManagerComponent,
	}

	SchedulerLabels = map[string]string{
		ControlPlaneComponentLabel: SchedulerComponent,
	}
)
