package install

import (
	"kubeception.cloud/kubeception/pkg/apis/kubeception/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
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
