package install

import (
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"kubeception.cloud/kubeception/pkg/apis/certificate/v1alpha1"
)

var (
	localSchemeBuilder = runtime.NewSchemeBuilder(
		v1alpha1.AddToScheme,
	)

	// AddToScheme adds all kubecdeption APIs to the given scheme.
	AddToScheme = localSchemeBuilder.AddToScheme
)

// Install installs all kubeception APIs to the given scheme, panicking if it fails.
func Install(scheme *runtime.Scheme) {
	utilruntime.Must(AddToScheme(scheme))
}
