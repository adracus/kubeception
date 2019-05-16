//go:generate go run github.com/golang/mock/mockgen -package=manager -destination=zz_mocks.go sigs.k8s.io/controller-runtime/pkg/manager Manager

package manager
