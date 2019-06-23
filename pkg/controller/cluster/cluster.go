package cluster

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	KubeconfigField = "kubeconfig"
)

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
