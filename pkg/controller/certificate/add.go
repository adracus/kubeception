package certificate

import (
	"kubeception.cloud/kubeception/pkg/controller/certificate/certificate"
	"kubeception.cloud/kubeception/pkg/controller/certificate/keypair"
	"kubeception.cloud/kubeception/pkg/util"
)

var (
	addToManagerBuilder = util.NewAddToManagerBuilder(
		certificate.AddToManager,
		keypair.AddToManager,
	)

	AddToManager = addToManagerBuilder.AddToManager
)
