package keypair

import (
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"

	"kubeception.cloud/kubeception/pkg/util"

	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
	"kubeception.cloud/kubeception/pkg/apis/certificate/v1alpha1"
)

func ReadSecret(secret *corev1.Secret) (*rsa.PrivateKey, error) {
	privateKeyData, ok := secret.Data[v1alpha1.PrivateKeyDataKey]
	if !ok {
		return nil, fmt.Errorf("private key data missing")
	}

	privateKey, err := DecodePrivateKey(privateKeyData)
	if err != nil {
		return nil, err
	}

	publicKeyData, ok := secret.Data[v1alpha1.PublicKeyDataKey]
	if !ok {
		return nil, fmt.Errorf("public key data missing")
	}

	publicKey, err := DecodePublicKey(publicKeyData)
	if err != nil {
		return nil, err
	}

	privateKey.PublicKey = *publicKey
	return privateKey, nil
}

func DecodePrivateKey(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil || block.Type != RSAPrivateKeyBlockType {
		return nil, fmt.Errorf("invalid PEM encoded RSA private key")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func EncodePrivateKey(privateKey *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  RSAPrivateKeyBlockType,
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
}

func DecodePublicKey(data []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil || block.Type != RSAPublicKeyBlockType {
		return nil, fmt.Errorf("invalid PEM encoded RSA public key")
	}

	return x509.ParsePKCS1PublicKey(block.Bytes)
}

func EncodePublicKey(publicKey *rsa.PublicKey) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  RSAPublicKeyBlockType,
		Bytes: x509.MarshalPKCS1PublicKey(publicKey),
	})
}

func EncodeKeyPair(pair *rsa.PrivateKey) (privateKeyData []byte, publicKeyData []byte) {
	return EncodePrivateKey(pair), EncodePublicKey(&pair.PublicKey)
}

func UpdateSecret(secret *corev1.Secret, privateKeyData, publicKeyData []byte) {
	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}
	secret.Data[v1alpha1.PrivateKeyDataKey] = privateKeyData
	secret.Data[v1alpha1.PublicKeyDataKey] = publicKeyData
}

func GetKeyPairFromSecret(ctx context.Context, c client.Client, key client.ObjectKey) (*rsa.PrivateKey, error) {
	secret := &corev1.Secret{}
	if err := c.Get(ctx, key, secret); err != nil {
		return nil, err
	}

	return ReadSecret(secret)
}

func ComputeChecksum(privateKeyData, publicKeyData []byte) (string, error) {
	data, err := json.Marshal(struct {
		PrivateKey []byte `json:"privateKey"`
		PublicKey  []byte `json:"publicKey"`
	}{
		PrivateKey: privateKeyData,
		PublicKey:  publicKeyData,
	})
	if err != nil {
		return "", err
	}

	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func UpdateChecksum(keyPair *v1alpha1.KeyPair, checksum string) {
	util.SetMetaDataAnnotation(keyPair, v1alpha1.KeyPairChecksumKey, checksum)
}
