package certificate

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"math/big"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeception.cloud/kubeception/pkg/apitypes"

	corev1 "k8s.io/api/core/v1"
	controllerruntime "sigs.k8s.io/controller-runtime"

	"github.com/go-logr/logr"
	"k8s.io/client-go/tools/record"
	"kubeception.cloud/kubeception/pkg/apis/certificate/v1alpha1"
	"kubeception.cloud/kubeception/pkg/controller/certificate/keypair"
	"kubeception.cloud/kubeception/pkg/util"
	"kubeception.cloud/kubeception/pkg/util/controller"
	"kubeception.cloud/kubeception/pkg/util/finalizer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

const (
	FinalizerName = "kubeception.cloud/certificate"
)

var logger = log.Log.WithName("certificate")

type reconciler struct {
	recorder record.EventRecorder
	controller.WithClient
	controller.WithScheme
	controller.WithContext
	controller.WithLog
}

func NewReconciler(recorder record.EventRecorder) reconcile.Reconciler {
	return &reconciler{recorder: recorder, WithLog: controller.NewWithLog(logger)}
}

func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	log := r.Log.WithValues("certificate", req.String())
	cert := &v1alpha1.Certificate{}
	if err := r.Client.Get(r.Context, req.NamespacedName, cert); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	return reconcile.Result{}, finalizer.Handle(r.Context, r.Client, FinalizerName, cert, &finalizer.Funcs{
		ReconcileFunc: func() error {
			return r.reconcile(r.Context, log, cert)
		},
	})
}

func (r *reconciler) setCertDefaults(ctx context.Context, log logr.Logger, cert *v1alpha1.Certificate) error {
	needsUpdate := false
	defer func() {
		if needsUpdate {
			if err := r.Client.Update(ctx, cert); err != nil {
				r.Log.Error(err, "Could not update outdated cert")
			}
		}
	}()

	if cert.Spec.KeyPair == nil {
		keyPair, err := r.createCertKeyPair(ctx, log, cert)
		if err != nil {
			return err
		}

		ref := util.LocalObjectReferenceToObject(keyPair)
		cert.Spec.KeyPair = &ref
		needsUpdate = true
	}

	if cert.Spec.Info.SerialNumber == nil {
		bigInt, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
		if err != nil {
			return err
		}

		serialNumber := apitypes.NewBigInt(bigInt)
		cert.Spec.Info.SerialNumber = &serialNumber
		needsUpdate = true
	}

	if cert.Spec.Info.NotBefore == nil {
		notBefore := metav1.NewTime(time.Now())
		cert.Spec.Info.NotBefore = &notBefore
		needsUpdate = true
	}

	if cert.Spec.Info.NotAfter == nil {
		notAfter := metav1.NewTime(cert.Spec.Info.NotBefore.Time.AddDate(10, 0, 0))
		cert.Spec.Info.NotAfter = &notAfter
		needsUpdate = true
	}

	return nil
}

func (r *reconciler) createCertKeyPair(ctx context.Context, log logr.Logger, cert *v1alpha1.Certificate) (*v1alpha1.KeyPair, error) {
	keyPair := &v1alpha1.KeyPair{
		ObjectMeta: controllerruntime.ObjectMeta{
			Namespace:    cert.Namespace,
			GenerateName: fmt.Sprintf("%s-keypair-", cert.Name),
		},
	}
	if err := controllerruntime.SetControllerReference(cert, keyPair, r.Scheme); err != nil {
		return nil, err
	}

	if err := r.Client.Create(ctx, keyPair); err != nil {
		return nil, err
	}

	return keyPair, nil
}

func (r *reconciler) getParentSigner(ctx context.Context, log logr.Logger, cert *v1alpha1.Certificate) (*x509.Certificate, *rsa.PrivateKey, error) {
	parentKey := util.Key(cert.Namespace, cert.Spec.Parent.Name)
	parent := &v1alpha1.Certificate{}
	if err := r.Client.Get(ctx, parentKey, parent); err != nil {
		return nil, nil, err
	}

	if parent.Spec.KeyPair == nil {
		return nil, nil, fmt.Errorf("parent certificate does not have a key pair linked to it")
	}

	signerKey, err := keypair.GetKeyPairFromSecret(ctx, r.Client, util.Key(parent.Namespace, parent.Spec.KeyPair.Name))
	if err != nil {
		return nil, nil, err
	}

	parentData, err := GetCertificateFromSecret(ctx, r.Client, parentKey)
	if err != nil {
		return nil, nil, err
	}

	return parentData, signerKey, nil
}

func (r *reconciler) generateCertificate(ctx context.Context, log logr.Logger, cert *v1alpha1.Certificate) ([]byte, error) {
	template, err := TemplateForCertificate(cert)
	if err != nil {
		return nil, err
	}

	privateKey, err := keypair.GetKeyPairFromSecret(ctx, r.Client, util.Key(cert.Namespace, cert.Spec.KeyPair.Name))
	if err != nil {
		return nil, err
	}

	parent := template
	signerKey := privateKey
	if cert.Spec.Parent != nil {
		parent, signerKey, err = r.getParentSigner(ctx, log, cert)
		if err != nil {
			return nil, err
		}
	}

	return x509.CreateCertificate(rand.Reader, template, parent, &privateKey.PublicKey, signerKey)
}

func (r *reconciler) reconcile(ctx context.Context, log logr.Logger, cert *v1alpha1.Certificate) error {
	var certData []byte
	secret := &corev1.Secret{ObjectMeta: util.ObjectMeta(cert.Namespace, cert.Name)}
	_, err := controllerruntime.CreateOrUpdate(ctx, r.Client, secret, func() error {
		if err := controllerruntime.SetControllerReference(cert, secret, r.Scheme); err != nil {
			return err
		}

		if err := r.setCertDefaults(ctx, log, cert); err != nil {
			return err
		}

		var err error
		r.recorder.Event(cert, corev1.EventTypeNormal, v1alpha1.EventGeneratingCertificate, "Generating Certificate")
		certData, err = r.generateCertificate(ctx, log, cert)
		if err != nil {
			r.recorder.Eventf(cert, corev1.EventTypeWarning, v1alpha1.EventErrorGenerateCertificate, "Could not generate certificate: %v", err)
			return err
		}

		UpdateSecret(secret, certData)
		return nil
	})
	if err != nil {
		return err
	}

	checksum := ComputeChecksum(certData)
	withoutChecksum := cert.DeepCopy()
	UpdateChecksum(cert, checksum)
	return r.Client.Patch(ctx, cert, client.MergeFrom(withoutChecksum))
}
