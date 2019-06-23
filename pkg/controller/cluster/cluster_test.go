package cluster

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func TestCluster(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cluster")
}

var _ = Describe("Cluster Suite", func() {
	var (
		kubeConfig    clientcmdapi.Config
		rawKubeConfig []byte
	)
	BeforeEach(func() {
		var err error
		kubeConfig = *clientcmdapi.NewConfig()
		rawKubeConfig, err = clientcmd.Write(kubeConfig)
		Expect(err).NotTo(HaveOccurred())
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
