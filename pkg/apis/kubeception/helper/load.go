package helper

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/runtime/serializer/versioning"
	"kubeception.cloud/kubeception/pkg/apis/kubeception/install"
	"kubeception.cloud/kubeception/pkg/apis/kubeception/v1alpha1"
)

const (
	Group = "kubeception.cloud"
)

var (
	Scheme *runtime.Scheme

	Codec runtime.Codec

	defaultClusterConfigGVK = v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.ClusterConfigKind)
	defaultMachineConfigGVK = v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.MachineConfigKind)
)

func init() {
	Scheme = runtime.NewScheme()
	install.Install(Scheme)

	yamlSerializer := json.NewYAMLSerializer(json.DefaultMetaFactory, Scheme, Scheme)
	Codec = versioning.NewDefaultingCodecForScheme(
		Scheme,
		yamlSerializer,
		yamlSerializer,
		v1alpha1.SchemeGroupVersion,
		v1alpha1.SchemeGroupVersion,
	)
}

func LoadClusterConfig(data []byte) (*v1alpha1.ClusterConfig, error) {
	config := &v1alpha1.ClusterConfig{}
	if _, _, err := Codec.Decode(data, &defaultClusterConfigGVK, config); err != nil {
		return nil, err
	}
	return config, nil
}

func LoadMachineConfig(data []byte) (*v1alpha1.MachineConfig, error) {
	config := &v1alpha1.MachineConfig{}
	if _, _, err := Codec.Decode(data, &defaultMachineConfigGVK, config); err != nil {
		return nil, err
	}
	return config, nil
}
