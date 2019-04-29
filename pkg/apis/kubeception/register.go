package kubeception

import (
	"github.com/adracus/kubeception/pkg/apis/kubeception/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	ProviderName = "kubeception"
)

var (
	localSchemeBuilder = runtime.NewSchemeBuilder(
		v1alpha1.AddToScheme,
	)

	AddToScheme = localSchemeBuilder.AddToScheme
)
