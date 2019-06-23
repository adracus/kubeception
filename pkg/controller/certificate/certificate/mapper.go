package certificate

import (
	"kubeception.cloud/kubeception/pkg/apis/certificate/v1alpha1"
	"kubeception.cloud/kubeception/pkg/util"
	"kubeception.cloud/kubeception/pkg/util/controller"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type keyPairToCertificateMapper struct {
	controller.WithClient
	controller.WithLog
	controller.WithContext
}

func NewKeyPairToCertificateMapper() handler.Mapper {
	return &keyPairToCertificateMapper{WithLog: controller.NewWithLog(logger.WithName("keypair-mapper"))}
}

func (k *keyPairToCertificateMapper) doMap(mapObject handler.MapObject) ([]reconcile.Request, error) {
	certList := &v1alpha1.CertificateList{}
	if err := k.Client.List(k.Context, certList, client.InNamespace(mapObject.Meta.GetNamespace())); err != nil {
		return nil, err
	}

	var requests []reconcile.Request
	for _, cert := range certList.Items {
		if keyPair := cert.Spec.KeyPair; keyPair != nil && keyPair.Name == mapObject.Meta.GetName() {
			requests = append(requests, util.RequestFromObject(&cert))
		}
	}
	return requests, nil
}

func (k *keyPairToCertificateMapper) Map(mapObject handler.MapObject) []reconcile.Request {
	requests, err := k.doMap(mapObject)
	if err != nil {
		k.Log.Error(err, "Could not map keypair", "keypair", util.KeyFromObject(mapObject.Meta).String())
		return nil
	}

	return requests
}

type certificateMapper struct {
	controller.WithClient
	controller.WithLog
	controller.WithContext
}

func NewCertificateMapper() handler.Mapper {
	return &certificateMapper{WithLog: controller.NewWithLog(logger.WithName("certificate-mapper"))}
}

func (k *certificateMapper) doMap(mapObject handler.MapObject) ([]reconcile.Request, error) {
	certList := &v1alpha1.CertificateList{}
	if err := k.Client.List(k.Context, certList, client.InNamespace(mapObject.Meta.GetNamespace())); err != nil {
		return nil, err
	}

	requests := []reconcile.Request{util.RequestFromObject(mapObject.Meta)}
	for _, cert := range certList.Items {
		if parent := cert.Spec.Parent; parent != nil && parent.Name == mapObject.Meta.GetName() {
			requests = append(requests, util.RequestFromObject(&cert))
		}
	}
	return requests, nil
}

func (k *certificateMapper) Map(mapObject handler.MapObject) []reconcile.Request {
	requests, err := k.doMap(mapObject)
	if err != nil {
		k.Log.Error(err, "Could not map certificate", "certificate", util.KeyFromObject(mapObject.Meta).String())
		return nil
	}

	return requests
}
