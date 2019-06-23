package v1alpha1

import (
	"kubeception.cloud/kubeception/pkg/apitypes"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Type string

const (
	CACert           Type = "CACert"
	ServerCert       Type = "ServerCert"
	ClientCert       Type = "ClientCert"
	ServerClientCert Type = "ServerClientCert"
)

const (
	PrivateKeyDataKey  = "private-key"
	PublicKeyDataKey   = "public-key"
	KeyPairChecksumKey = "keypair.certificate.kubeception.cloud/checksum"

	CertificateDataKey     = "certificate"
	CertificateChecksumKey = "certificate.certificate.kubeception.cloud/checksum"

	EventGeneratingKey            = "GeneratingKey"
	EventGeneratingCertificate    = "GeneratingCertificate"
	EventErrorGenerateCertificate = "GenerateCertificateError"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// KeyPair is an RSA private and public key pair.
type KeyPair struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeyPairSpec   `json:"spec"`
	Status KeyPairStatus `json:"status"`
}

// +kubebuilder:object:root=true

// KeyPairList contains a list of KeyPair.
type KeyPairList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []KeyPair `json:"items,omitempty"`
}

const (
	SecretsSelfProvisioned = "SelfProvisioned"
)

type KeyPairSpec struct {
	Secrets string `json:"secrets,omitempty"`
}

type KeyPairStatus struct {
	Conditions []KeyPairCondition `json:"conditions,omitempty"`
}

type KeyPairConditionType string

const (
	KeyPairPresent KeyPairConditionType = "Present"
	KeyPairValid   KeyPairConditionType = "Valid"
)

type KeyPairCondition struct {
	Type   KeyPairConditionType   `json:"type"`
	Status corev1.ConditionStatus `json:"status"`
	// The last time this condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
}

type CertificateSubject struct {
	CommonName   string   `json:"commonName"`
	Organization []string `json:"organization,omitempty"`
}

type CertificateInfo struct {
	SerialNumber *apitypes.BigInt `json:"serialNumber,omitempty"`

	NotBefore *metav1.Time       `json:"notBefore,omitempty"`
	NotAfter  *metav1.Time       `json:"notAfter,omitempty"`
	Subject   CertificateSubject `json:"subject"`

	DNSNames    []string      `json:"dnsNames,omitempty"`
	IPAddresses []apitypes.IP `json:"ipAddresses,omitempty"`
}

// +kubebuilder:object:root=true

// Certificate is a certificate of a certification authority backed by an authority.
type Certificate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec CertificateSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// CertificateList is a list of Certificates.
type CertificateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Certificate `json:"items,omitempty"`
}

type CertificateSpec struct {
	Type    Type                         `json:"type"`
	Info    CertificateInfo              `json:"info"`
	KeyPair *corev1.LocalObjectReference `json:"keyPair,omitempty"`
	Parent  *corev1.LocalObjectReference `json:"parent,omitempty"`
}

type CertificateSigner struct {
	KeyPairRef     *corev1.LocalObjectReference `json:"keyPairRef,omitempty"`
	CertificateRef *corev1.LocalObjectReference `json:"certificateRef,omitempty"`
}
