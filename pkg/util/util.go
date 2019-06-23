package util

import (
	"context"
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

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

func nameAndNamespace(namespaceOrName string, nameOpt ...string) (namespace, name string) {
	if len(nameOpt) > 1 {
		panic(fmt.Sprintf("more than name/namespace for key specified: %s/%v", namespaceOrName, nameOpt))
	}
	if len(nameOpt) == 0 {
		name = namespaceOrName
		return
	}
	namespace = namespaceOrName
	name = nameOpt[0]
	return
}

// Key creates a new client.ObjectKey from the given parameters.
// There are only two ways to call this function:
// - If only namespaceOrName is set, then a client.ObjectKey with name set to namespaceOrName is returned.
// - If namespaceOrName and one nameOpt is given, then a client.ObjectKey with namespace set to namespaceOrName
//   and name set to nameOpt[0] is returned.
// For all other cases, this method panics.
func Key(namespaceOrName string, nameOpt ...string) client.ObjectKey {
	namespace, name := nameAndNamespace(namespaceOrName, nameOpt...)
	return client.ObjectKey{Namespace: namespace, Name: name}
}

// Request creates a new reconcile.Request from the given parameters.
// There are only two ways to call this function:
// - If only namespaceOrName is set, then a client.ObjectKey with name set to namespaceOrName is returned.
// - If namespaceOrName and one nameOpt is given, then a client.ObjectKey with namespace set to namespaceOrName
//   and name set to nameOpt[0] is returned.
// For all other cases, this method panics.
func Request(namespaceOrName string, nameOpt ...string) reconcile.Request {
	return reconcile.Request{NamespacedName: Key(namespaceOrName, nameOpt...)}
}

func RequestFromObject(obj metav1.Object) reconcile.Request {
	return reconcile.Request{NamespacedName: Key(obj.GetNamespace(), obj.GetName())}
}

// ObjectMeta creates a new metav1.ObjectMeta from the given parameters.
// There are only two ways to call this function:
// - If only namespaceOrName is set, then a metav1.ObjectMeta with name set to namespaceOrName is returned.
// - If namespaceOrName and one nameOpt is given, then a metav1.ObjectMeta with namespace set to namespaceOrName
//   and name set to nameOpt[0] is returned.
// For all other cases, this method panics.
func ObjectMeta(namespaceOrName string, nameOpt ...string) metav1.ObjectMeta {
	namespace, name := nameAndNamespace(namespaceOrName, nameOpt...)
	return metav1.ObjectMeta{Namespace: namespace, Name: name}
}

// KeyFromObject obtains the ObjectKey from the given Object.
func KeyFromObject(obj metav1.Object) client.ObjectKey {
	return client.ObjectKey{Namespace: obj.GetNamespace(), Name: obj.GetName()}
}

func LocalObjectReferenceToObject(obj metav1.Object) corev1.LocalObjectReference {
	return corev1.LocalObjectReference{Name: obj.GetName()}
}

// LogErrorAndExit logs the error and exits with code 1.
func LogErrorAndExit(log logr.Logger, err error, msg string, keysAndValues ...interface{}) {
	log.Error(err, msg, keysAndValues...)
	os.Exit(1)
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

func SetMetaDataAnnotation(obj metav1.Object, key, value string) {
	annotations := obj.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}

	annotations[key] = value
	obj.SetAnnotations(annotations)
}

// SetMetaDataAnnotations sets all new annotations on the given metav1.Object.
func SetMetaDataAnnotations(obj metav1.Object, newAnnotations map[string]string) {
	annotations := obj.GetAnnotations()
	if annotations == nil {
		annotations = newAnnotations
	} else {
		for k, v := range newAnnotations {
			annotations[k] = v
		}
	}
	obj.SetAnnotations(annotations)
}

func TypeToKind(obj runtime.Object) (string, error) {
	v, err := conversion.EnforcePtr(obj)
	if err != nil {
		return "", err
	}

	return v.Type().Name(), nil
}

func MustTypeToKind(obj runtime.Object) string {
	kind, err := TypeToKind(obj)
	utilruntime.Must(err)
	return kind
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
