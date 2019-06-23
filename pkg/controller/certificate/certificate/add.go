package certificate

import (
	corev1 "k8s.io/api/core/v1"
	"kubeception.cloud/kubeception/pkg/apis/certificate/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	Name = "certificate"
)

type AddArgs struct {
	MaxConcurrentReconciles int
}

var DefaultArgs AddArgs

func AddToManager(mgr manager.Manager) error {
	return AddToManagerWithArgs(mgr, DefaultArgs)
}

func AddToManagerWithArgs(mgr manager.Manager, args AddArgs) error {
	ctrl, err := controller.New(Name, mgr, controller.Options{
		Reconciler:              NewReconciler(mgr.GetEventRecorderFor(Name)),
		MaxConcurrentReconciles: args.MaxConcurrentReconciles,
	})
	if err != nil {
		return err
	}

	if err := ctrl.Watch(&source.Kind{Type: &v1alpha1.Certificate{}}, &handler.EnqueueRequestsFromMapFunc{ToRequests: NewCertificateMapper()}); err != nil {
		return err
	}

	if err := ctrl.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{OwnerType: &v1alpha1.Certificate{}, IsController: true}); err != nil {
		return err
	}

	if err := ctrl.Watch(&source.Kind{Type: &v1alpha1.KeyPair{}}, &handler.EnqueueRequestsFromMapFunc{ToRequests: NewKeyPairToCertificateMapper()}); err != nil {
		return err
	}

	return nil
}
