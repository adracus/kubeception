//go:generate mockgen -package=util -destination=zz_funcs.go github.com/adracus/kubeception/pkg/internal/mock/kubeception/intern/util AddToManager

package util

import "sigs.k8s.io/controller-runtime/pkg/manager"

// AddToManager allows mocking `func(manager.Manager) error`.
type AddToManager interface {
	Do(manager.Manager) error
}
