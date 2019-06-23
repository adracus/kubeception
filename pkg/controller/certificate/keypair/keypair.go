package keypair

import (
	"context"
	"crypto/rand"
	"crypto/rsa"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/go-logr/logr"

	"k8s.io/client-go/tools/record"

	"kubeception.cloud/kubeception/pkg/util/controller"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"kubeception.cloud/kubeception/pkg/apis/certificate/v1alpha1"
	"kubeception.cloud/kubeception/pkg/util"
	"kubeception.cloud/kubeception/pkg/util/finalizer"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	FinalizerName = "kubeception.cloud/keypair"

	RSAPrivateKeyBlockType = "RSA PRIVATE KEY"
	RSAPublicKeyBlockType  = "RSA PUBLIC KEY"
)

var logger = log.Log.WithName("keypair")

type reconciler struct {
	recorder record.EventRecorder
	controller.WithScheme
	controller.WithContext
	controller.WithClient
	controller.WithLog
}

func NewReconciler(recorder record.EventRecorder) reconcile.Reconciler {
	return &reconciler{recorder: recorder, WithLog: controller.NewWithLog(logger)}
}

func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	log := r.Log.WithValues("keypair", req.String())
	keyPair := &v1alpha1.KeyPair{}
	if err := r.Client.Get(r.Context, req.NamespacedName, keyPair); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, finalizer.Handle(r.Context, r.Client, FinalizerName, keyPair, finalizer.Funcs{
		ReconcileFunc: func() error {
			return r.reconcile(r.Context, log, keyPair)
		},
	})
}

func (r *reconciler) reconcile(ctx context.Context, log logr.Logger, keyPair *v1alpha1.KeyPair) error {
	if keyPair.Spec.Secrets == v1alpha1.SecretsSelfProvisioned {
		log.Info("Secret is self provisioned, skipping")
		return nil
	}

	var privateKeyData, publicKeyData []byte
	secret := &corev1.Secret{ObjectMeta: util.ObjectMeta(keyPair.Namespace, keyPair.Name)}
	_, err := controllerruntime.CreateOrUpdate(ctx, r.Client, secret, func() error {
		if err := controllerruntime.SetControllerReference(keyPair, secret, r.Scheme); err != nil {
			return err
		}

		privateKey, err := ReadSecret(secret)
		if err != nil {
			log.Info("Generating RSA Key")
			r.recorder.Event(keyPair, corev1.EventTypeNormal, v1alpha1.EventGeneratingKey, "Generating RSA Key")
			privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				return err
			}
		}

		privateKeyData, publicKeyData = EncodeKeyPair(privateKey)
		UpdateSecret(secret, privateKeyData, publicKeyData)
		return nil
	})
	if err != nil {
		return err
	}

	checksum, err := ComputeChecksum(privateKeyData, publicKeyData)
	if err != nil {
		return err
	}

	withoutChecksum := keyPair.DeepCopy()
	UpdateChecksum(keyPair, checksum)
	return r.Client.Patch(ctx, keyPair, client.MergeFrom(withoutChecksum))
}
