package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	GroupName = "kubeception.io"
	Version   = "v1alpha1"
)

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: Version}

func Resource(resource string) schema.GroupResource {
	return schema.GroupResource{Group: GroupName, Resource: resource}
}

var (
	localSchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)

	AddToScheme = localSchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&ClusterConfig{},
	)
	return nil
}
