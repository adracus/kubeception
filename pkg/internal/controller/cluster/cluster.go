package cluster

import (
	"fmt"
	kubeceptionv1alpha1 "kubeception.cloud/kubeception/pkg/apis/kubeception/v1alpha1"
	"kubeception.cloud/kubeception/pkg/internal/scheme"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	clusterv1alpha1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

const (
	KubeconfigField = "kubeconfig"
)

// ConfigFromCluster tries to decode a ClusterConfig from the given cluster.
func ConfigFromCluster(cluster *clusterv1alpha1.Cluster) (*kubeceptionv1alpha1.ClusterConfig, error) {
	config := &kubeceptionv1alpha1.ClusterConfig{}
	if _, _, err := scheme.Decoder.Decode(cluster.Spec.ProviderSpec.Value.Raw, nil, config); err != nil {
		return nil, err
	}

	return config, nil
}

// ReadKubeconfigSecret reads the clientcmdapi.Config from the given secret.
func ReadKubeconfigSecret(secret *corev1.Secret) (*clientcmdapi.Config, error) {
	if secret.Data == nil {
		return nil, fmt.Errorf("secret does not contain data")
	}

	data, ok := secret.Data[KubeconfigField]
	if !ok {
		return nil, fmt.Errorf("secret does not contain kubeconfig at %q", KubeconfigField)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("kubeconfig is empty")
	}

	return clientcmd.Load(data)
}

// UpdateKubeconfigSecret updates the given secret to contain the given clientcmdapi.Config at the data KubeconfigField.
func UpdateKubeconfigSecret(secret *corev1.Secret, config *clientcmdapi.Config) error {
	data, err := clientcmd.Write(*config)
	if err != nil {
		return err
	}

	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}

	secret.Data[KubeconfigField] = data
	return nil
}
