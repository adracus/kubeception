package finalizer

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func containsString(items []string, s string) bool {
	for _, item := range items {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(items []string, s string) (bool, []string) {
	var (
		ok  bool
		out []string
	)

	for _, item := range items {
		if item == s {
			ok = true
			continue
		}

		out = append(out, item)
	}

	return ok, out
}

// Has checks if the given object has a finalizer with the given name.
func Has(finalizer string, obj runtime.Object) (bool, error) {
	acc, err := meta.Accessor(obj)
	if err != nil {
		return false, err
	}

	return containsString(acc.GetFinalizers(), finalizer), nil
}

// Add adds the finalizer to the object. It returns a boolean whether the finalizer
// was added (i.e. was not yet present in the finalizer list) or not.
func Add(finalizer string, obj runtime.Object) (bool, error) {
	acc, err := meta.Accessor(obj)
	if err != nil {
		return false, err
	}

	finalizers := acc.GetFinalizers()
	if !containsString(finalizers, finalizer) {
		finalizers = append(finalizers, finalizer)
		acc.SetFinalizers(finalizers)
		return true, nil
	}
	return false, nil
}

// Remove removes the finalizer from the object. It returns a boolean whether the finalizer
// was removed (i.e. was present in the finalizer list) or not.
func Remove(finalizer string, obj runtime.Object) (bool, error) {
	acc, err := meta.Accessor(obj)
	if err != nil {
		return false, err
	}

	ok, finalizers := removeString(acc.GetFinalizers(), finalizer)
	acc.SetFinalizers(finalizers)
	return ok, nil
}

// Handler are functions that handle the finalizing reconcile process.
//
// Reconcile is called if an object is not being deleted (i.e. deletionTimestamp is zero).
// Finalize is called if an object is being deleted (i.e. deletionTimestamp is non-zero).
type Handler interface {
	// Reconcile is called if an object is not being deleted.
	// It may be assumed that the remote object already has a finalizer present.
	Reconcile() error
	// Finalize is called if an object is being deleted.
	// If Finalize returns with no error, the finalizer shall be removed.
	Finalize() error
}

// Funcs are functions that implement the Handler interface.
// If a function is nil, it won't be called and no error is returned.
type Funcs struct {
	// ReconcileFunc is a function that is called upon reconciliation of an object.
	ReconcileFunc func() error
	// FinalizeFunc is a function that is called upon finalization of an object.
	FinalizeFunc func() error
}

// Reconcile implements Handler.
func (f Funcs) Reconcile() error {
	if f.ReconcileFunc != nil {
		return f.ReconcileFunc()
	}
	return nil
}

// Finalize implements Handler.
func (f Funcs) Finalize() error {
	if f.FinalizeFunc != nil {
		return f.FinalizeFunc()
	}
	return nil
}

// Finalize finalizes the object with the function if the finalizer is present on the object.
// If the finalization function has returned with no error, the finalizer is removed and the object is patched.
func Finalize(ctx context.Context, c client.Client, finalizer string, obj runtime.Object, f func() error) error {
	if ok, err := Has(finalizer, obj); err != nil || !ok {
		return err
	}

	if err := f(); err != nil {
		return err
	}

	withFinalizer := obj.DeepCopyObject()
	if ok, err := Remove(finalizer, obj); err != nil || !ok {
		return err
	}

	return c.Patch(ctx, obj, client.MergeFrom(withFinalizer))
}

// Reconcile reconciles the object with the function, adding the finalizer before and patching the object if not present.
func Reconcile(ctx context.Context, c client.Client, finalizer string, obj runtime.Object, f func() error) error {
	withoutFinalizer := obj.DeepCopyObject()
	ok, err := Add(finalizer, obj)
	if err != nil {
		return err
	}

	if ok {
		if err := c.Patch(ctx, obj, client.MergeFrom(withoutFinalizer)); err != nil {
			return err
		}
	}

	return f()
}

// Handle either Reconciles or Finalizes the object depending on the deletionTimestamp of the object.
func Handle(ctx context.Context, c client.Client, finalizer string, obj runtime.Object, handler Handler) error {
	acc, err := meta.Accessor(obj)
	if err != nil {
		return err
	}

	if acc.GetDeletionTimestamp().IsZero() {
		return Reconcile(ctx, c, finalizer, obj, handler.Reconcile)
	}
	return Finalize(ctx, c, finalizer, obj, handler.Finalize)
}
