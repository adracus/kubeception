package certificate

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"fmt"
	"net"

	"kubeception.cloud/kubeception/pkg/util"

	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"

	"kubeception.cloud/kubeception/pkg/apis/certificate/v1alpha1"
)

func TemplateForCertificate(cert *v1alpha1.Certificate) (*x509.Certificate, error) {
	var ipAddresses []net.IP
	for _, ipAddress := range cert.Spec.Info.IPAddresses {
		ipAddresses = append(ipAddresses, ipAddress.IP)
	}

	template := x509.Certificate{
		BasicConstraintsValid: true,
		IsCA:                  cert.Spec.Type == v1alpha1.CACert,
		SerialNumber:          &cert.Spec.Info.SerialNumber.BigInt,
		NotBefore:             cert.Spec.Info.NotBefore.Time,
		NotAfter:              cert.Spec.Info.NotAfter.Time,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		Subject: pkix.Name{
			CommonName:   cert.Spec.Info.Subject.CommonName,
			Organization: cert.Spec.Info.Subject.Organization,
		},
		DNSNames:    cert.Spec.Info.DNSNames,
		IPAddresses: ipAddresses,
	}

	switch cert.Spec.Type {
	case v1alpha1.CACert:
		template.KeyUsage |= x509.KeyUsageCertSign | x509.KeyUsageCRLSign
	case v1alpha1.ServerCert:
		template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	case v1alpha1.ClientCert:
		template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	case v1alpha1.ServerClientCert:
		template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}
	}

	return &template, nil
}

func ReadSecret(secret *corev1.Secret) (*x509.Certificate, error) {
	certData, ok := secret.Data[v1alpha1.CertificateDataKey]
	if !ok {
		return nil, fmt.Errorf("certificate data missing")
	}

	return x509.ParseCertificate(certData)
}

func UpdateSecret(secret *corev1.Secret, certData []byte) {
	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}

	secret.Data[v1alpha1.CertificateDataKey] = certData
}

func GetCertificateFromSecret(ctx context.Context, c client.Client, key client.ObjectKey) (*x509.Certificate, error) {
	secret := &corev1.Secret{}
	if err := c.Get(ctx, key, secret); err != nil {
		return nil, err
	}

	return ReadSecret(secret)
}

func ComputeChecksum(certData []byte) string {
	sum := sha256.Sum256(certData)
	return hex.EncodeToString(sum[:])
}

func UpdateChecksum(cert *v1alpha1.Certificate, checksum string) {
	util.SetMetaDataAnnotation(cert, v1alpha1.CertificateChecksumKey, checksum)
}
