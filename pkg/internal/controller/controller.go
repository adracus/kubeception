package controller

import "fmt"

const (
	LabelPrefix = "kubeception.io"
)

var (
	ControlPlaneComponentLabel = fmt.Sprintf("%s/control-plane-component", LabelPrefix)

	ETCDComponent      = "etcd"
	APIServerComponent = "kube-apiserver"
)
