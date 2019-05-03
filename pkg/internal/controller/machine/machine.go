package machine

import (
	kubeceptionv1alpha1 "github.com/adracus/kubeception/pkg/apis/kubeception/v1alpha1"
	"github.com/adracus/kubeception/pkg/internal/scheme"
	clusterv1alpha1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

func ConfigFromMachine(cluster *clusterv1alpha1.Machine) (*kubeceptionv1alpha1.MachineConfig, error) {
	config := &kubeceptionv1alpha1.MachineConfig{}
	if _, _, err := scheme.Decoder.Decode(cluster.Spec.ProviderSpec.Value.Raw, nil, config); err != nil {
		return nil, err
	}

	return config, nil
}
