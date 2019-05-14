package cluster

import (
	"kubeception.cloud/kubeception/pkg/apis/kubeception/v1alpha1"
	"kubeception.cloud/kubeception/pkg/internal/scheme"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	clusterv1alpha1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	"testing"
)

func TestCluster(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cluster")
}

var _ = Describe("Cluster Suite", func() {
	var (
		config    v1alpha1.ClusterConfig
		rawConfig []byte

		kubeConfig    clientcmdapi.Config
		rawKubeConfig []byte
	)
	BeforeEach(func() {
		config = v1alpha1.ClusterConfig{
			ControlPlane: v1alpha1.ControlPlane{},
		}

		var err error
		rawConfig, err = runtime.Encode(json.NewSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, false), &config)
		Expect(err).NotTo(HaveOccurred())

		kubeConfig = *clientcmdapi.NewConfig()
		rawKubeConfig, err = clientcmd.Write(kubeConfig)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("#ConfigFromCluster", func() {
		It("should correctly extract the configuration from the cluster", func() {
			actual, err := ConfigFromCluster(&clusterv1alpha1.Cluster{
				Spec: clusterv1alpha1.ClusterSpec{
					ProviderSpec: clusterv1alpha1.ProviderSpec{
						Value: &runtime.RawExtension{
							Raw: rawConfig,
						},
					},
				},
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(actual).To(Equal(&config))
		})
	})

	Describe("#ReadKubeconfigSecret", func() {
		It("should correctly read the kubeconfig secret", func() {
			actual, err := ReadKubeconfigSecret(&corev1.Secret{
				Data: map[string][]byte{
					KubeconfigField: rawKubeConfig,
				},
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(actual).To(Equal(&kubeConfig))
		})

		It("should error in case the secret does not have a data section", func() {
			_, err := ReadKubeconfigSecret(&corev1.Secret{})
			Expect(err).To(HaveOccurred())
		})

		It("should error if the secret data does not have a kubeconfig field", func() {
			_, err := ReadKubeconfigSecret(&corev1.Secret{Data: map[string][]byte{}})
			Expect(err).To(HaveOccurred())
		})

		It("should error if the kubeconfig data is empty", func() {
			_, err := ReadKubeconfigSecret(&corev1.Secret{Data: map[string][]byte{KubeconfigField: {}}})
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("#UpdateKubeconfigSecret", func() {
		It("should correctly update the kubeconfig secret", func() {
			secret := &corev1.Secret{}

			Expect(UpdateKubeconfigSecret(secret, &kubeConfig)).To(Succeed())

			Expect(secret.Data).To(Equal(map[string][]byte{
				KubeconfigField: rawKubeConfig,
			}))
		})
	})
})
