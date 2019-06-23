package machine

import (
	"fmt"

	"kubeception.cloud/kubeception/pkg/controller/common"
)

var (
	StatefulSetNameLabel = fmt.Sprintf("%s/machine", common.LabelPrefix)
)

func StatefulSetLabels(name string) map[string]string {
	return map[string]string{
		StatefulSetNameLabel: name,
	}
}
