package machine

import (
	"github.com/adracus/kubeception/pkg/apis/kubeception/v1alpha1"
	"github.com/adracus/kubeception/pkg/internal/scheme"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	clusterv1alpha1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	"testing"
)

func TestMachine(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Machine")
}

var _ = Describe("Machine Suite", func() {
	var (
		config    v1alpha1.MachineConfig
		rawConfig []byte
	)
	BeforeEach(func() {
		config = v1alpha1.MachineConfig{}

		var err error
		rawConfig, err = runtime.Encode(json.NewSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, false), &config)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("#ConfigFromMachine", func() {
		It("should correctly extract the configuration from the machine", func() {
			actual, err := ConfigFromMachine(&clusterv1alpha1.Machine{
				Spec: clusterv1alpha1.MachineSpec{
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
})
