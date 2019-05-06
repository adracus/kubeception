package util

import (
	"fmt"
	"github.com/adracus/kubeception/pkg/apis/kubeception/v1alpha1"
	mockmanager "github.com/adracus/kubeception/pkg/internal/mock/controller-runtime/manager"
	mockutil "github.com/adracus/kubeception/pkg/internal/mock/kubeception/intern/util"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"testing"
)

func TestMachine(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Util")
}

var _ = Describe("Utils Suite", func() {
	var (
		ctrl *gomock.Controller
	)
	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})
	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("#HyperkubeImageForConfig", func() {
		It("should correctly compute the kubernetes image for the given configuration", func() {
			Expect(HyperkubeImageForConfig(&v1alpha1.ClusterConfig{KubernetesVersion: "v1.13.5"})).
				To(Equal(fmt.Sprintf("%s:v1.13.5", HyperkubeRepository)))
		})
	})

	Describe("#SetMetaDataLabels", func() {
		It("should add the metadata section and set all given labels", func() {
			obj := &corev1.ConfigMap{}
			labels := map[string]string{
				"foo": "bar",
				"baz": "bang",
			}

			SetMetaDataLabels(obj, labels)

			Expect(obj.Labels).To(Equal(labels))
		})

		It("should merge the existing and given labels correctly", func() {
			obj := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"foo": "foo",
						"qux": "baz",
					},
				},
			}
			labels := map[string]string{
				"foo": "bar",
				"baz": "bang",
			}

			SetMetaDataLabels(obj, labels)

			Expect(obj.Labels).To(Equal(map[string]string{
				"foo": "bar",
				"qux": "baz",
				"baz": "bang",
			}))
		})
	})

	Describe("#IgnoreNotFound", func() {
		It("should ignore a NotFound error", func() {
			Expect(IgnoreNotFound(apierrors.NewNotFound(schema.GroupResource{}, ""))).To(BeNil())
		})

		It("should ignore nil errrors", func() {
			Expect(IgnoreNotFound(nil)).To(BeNil())
		})

		It("should not ignore non-not found errors", func() {
			err := fmt.Errorf("other error")
			Expect(IgnoreNotFound(err)).To(BeIdenticalTo(err))
		})
	})

	Context("AddToManagerBuilder", func() {
		Describe("#NewAddToManagerBuilder", func() {
			It("should create a new AddToManagerBuilder that contains the given functions", func() {
				mgr := mockmanager.NewMockManager(ctrl)
				f1 := mockutil.NewMockAddToManager(ctrl)
				f2 := mockutil.NewMockAddToManager(ctrl)

				gomock.InOrder(
					f1.EXPECT().Do(mgr),
					f2.EXPECT().Do(mgr),
				)

				addToManagerBuilder := NewAddToManagerBuilder(f1.Do, f2.Do)
				Expect(addToManagerBuilder.AddToManager(mgr)).To(Succeed())
			})
		})

		Describe("#Register", func() {
			It("should allow registering other functions once it was initialized", func() {
				mgr := mockmanager.NewMockManager(ctrl)
				f1 := mockutil.NewMockAddToManager(ctrl)
				f2 := mockutil.NewMockAddToManager(ctrl)
				f3 := mockutil.NewMockAddToManager(ctrl)

				gomock.InOrder(
					f1.EXPECT().Do(mgr),
					f2.EXPECT().Do(mgr),
					f3.EXPECT().Do(mgr),
				)

				addToManagerBuilder := NewAddToManagerBuilder(f1.Do, f2.Do)
				addToManagerBuilder.Register(f3.Do)

				Expect(addToManagerBuilder.AddToManager(mgr)).To(Succeed())
			})
		})

		Describe("#AddToManager", func() {
			It("should abort and immediately return the error if a function fails", func() {
				mgr := mockmanager.NewMockManager(ctrl)
				f1 := mockutil.NewMockAddToManager(ctrl)
				f2 := mockutil.NewMockAddToManager(ctrl)
				f3 := mockutil.NewMockAddToManager(ctrl)
				err := fmt.Errorf("expected error")

				gomock.InOrder(
					f1.EXPECT().Do(mgr),
					f2.EXPECT().Do(mgr).Return(err),
				)

				addToManagerBuilder := NewAddToManagerBuilder(f1.Do, f2.Do, f3.Do)

				Expect(addToManagerBuilder.AddToManager(mgr)).To(Equal(err))
			})
		})
	})
})
