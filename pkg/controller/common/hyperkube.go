package common

import (
	"fmt"

	"kubeception.cloud/kubeception/pkg/apis/kubeception/v1alpha1"
)

const (
	HyperkubeRepository = "k8s.gcr.io/hyperkube"
)

// HyperkubeImageForConfig returns the proper hyperkube image for the given cluster configuration.
func HyperkubeImageForConfig(config *v1alpha1.ClusterConfig) string {
	return fmt.Sprintf("%s:%s", HyperkubeRepository, config.KubernetesVersion)
}
