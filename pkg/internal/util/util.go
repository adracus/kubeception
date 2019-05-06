package util

import (
	"context"
	"fmt"
	"github.com/adracus/kubeception/pkg/apis/kubeception/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	HyperkubeRepository = "k8s.gcr.io/hyperkube"
)

func finalizersAndAccessorOf(obj runtime.Object) (sets.String, metav1.Object, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return nil, nil, err
	}

	return sets.NewString(accessor.GetFinalizers()...), accessor, nil
}

// HyperkubeImageForConfig returns the proper hyperkube image for the given cluster configuration.
func HyperkubeImageForConfig(config *v1alpha1.ClusterConfig) string {
	return fmt.Sprintf("%s:%s", HyperkubeRepository, config.KubernetesVersion)
}

// HasFinalizer checks if the given object has a finalizer with the given name.
func HasFinalizer(obj runtime.Object, finalizerName string) (bool, error) {
	finalizers, _, err := finalizersAndAccessorOf(obj)
	if err != nil {
		return false, err
	}

	return finalizers.Has(finalizerName), nil
}

// EnsureFinalizer ensures that a finalizer of the given name is set on the given object.
// If the finalizer is not set, it adds it to the list of finalizers and updates the remote object.
func EnsureFinalizer(ctx context.Context, client client.Client, finalizerName string, obj runtime.Object) error {
	finalizers, accessor, err := finalizersAndAccessorOf(obj)
	if err != nil {
		return err
	}

	if finalizers.Has(finalizerName) {
		return nil
	}

	finalizers.Insert(finalizerName)
	accessor.SetFinalizers(finalizers.UnsortedList())

	return client.Update(ctx, obj)
}

// DeleteFinalizer ensures that the given finalizer is not present anymore in the given object.
// If it is set, it removes it and issues an update.
func DeleteFinalizer(ctx context.Context, client client.Client, finalizerName string, obj runtime.Object) error {
	finalizers, accessor, err := finalizersAndAccessorOf(obj)
	if err != nil {
		return err
	}

	if !finalizers.Has(finalizerName) {
		return nil
	}

	finalizers.Delete(finalizerName)
	accessor.SetFinalizers(finalizers.UnsortedList())

	return client.Update(ctx, obj)
}

// ContextFromStopChannel instantiates a context that is open as long as the stopCh is open.
// It will return context.ErrCanceled as soon as the stopCh is closed.
func ContextFromStopChannel(stopCh <-chan struct{}) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		<-stopCh
	}()
	return ctx
}

// SetMetaDataLabel sets the given key value pair on the given metav1.Object.
func SetMetaDataLabel(obj metav1.Object, key, value string) {
	labels := obj.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}

	labels[key] = value
	obj.SetLabels(labels)
}

// SetMetaDataLabels sets all new labels on the given metav1.Object.
func SetMetaDataLabels(obj metav1.Object, newLabels map[string]string) {
	labels := obj.GetLabels()
	if labels == nil {
		labels = newLabels
	} else {
		for k, v := range newLabels {
			labels[k] = v
		}
	}
	obj.SetLabels(labels)
}

// IgnoreNotFound ignores `apierrors.IsNotFound` errors and returns `nil` if it encounters them.
func IgnoreNotFound(err error) error {
	if apierrors.IsNotFound(err) {
		return nil
	}
	return err
}

// AddToManagerBuilder aggregates various AddToManager functions.
type AddToManagerBuilder []func(manager.Manager) error

// NewAddToManagerBuilder creates a new AddToManagerBuilder and registers the given functions.
func NewAddToManagerBuilder(funcs ...func(manager.Manager) error) AddToManagerBuilder {
	var builder AddToManagerBuilder
	builder.Register(funcs...)
	return builder
}

// Register registers the given functions in this builder.
func (a *AddToManagerBuilder) Register(funcs ...func(manager.Manager) error) {
	*a = append(*a, funcs...)
}

// AddToManager traverses over all AddToManager-functions of this builder, sequentially applying
// them. It exits on the first error and returns it.
func (a *AddToManagerBuilder) AddToManager(m manager.Manager) error {
	for _, f := range *a {
		if err := f(m); err != nil {
			return err
		}
	}
	return nil
}
