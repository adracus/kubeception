package controller

import (
	"context"
	"os"

	"github.com/go-logr/logr"

	"k8s.io/apimachinery/pkg/runtime"
	"kubeception.cloud/kubeception/pkg/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type WithContext struct {
	Context context.Context
}

func (w *WithContext) InjectStopChannel(stopChan <-chan struct{}) error {
	w.Context = util.ContextFromStopChannel(stopChan)
	return nil
}

func NewWithContext(ctx context.Context) WithContext {
	return WithContext{ctx}
}

type WithClient struct {
	Client client.Client
}

func NewWithClient(c client.Client) WithClient {
	return WithClient{c}
}

func (w *WithClient) InjectClient(client client.Client) error {
	w.Client = client
	return nil
}

type WithScheme struct {
	Scheme *runtime.Scheme
}

func (w *WithScheme) InjectScheme(scheme *runtime.Scheme) error {
	w.Scheme = scheme
	return nil
}

func NewWithScheme(scheme *runtime.Scheme) WithScheme {
	return WithScheme{scheme}
}

type WithLog struct {
	Log logr.Logger
}

func (w *WithLog) InjectLogger(log logr.Logger) error {
	os.Exit(1)
	return nil
}

func NewWithLog(log logr.Logger) WithLog {
	return WithLog{log}
}
