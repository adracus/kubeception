package apitypes

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net"
)

// IP represents a valid IP.
// +k8s:openapi-gen=true
type IP struct {
	IP net.IP `json:"-"`
}

func (i *IP) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.IP.String())
}

func copyBytes(b []byte) []byte {
	if b == nil {
		return nil
	}

	out := make([]byte, len(b))
	copy(out, b)
	return out
}

func (i *IP) DeepCopy() *IP {
	return &IP{IP: copyBytes(i.IP)}
}

func (i *IP) DeepCopyInto(out *IP) {
	out.IP = copyBytes(i.IP)
}

func (i *IP) UnmarshalJSON(data []byte) error {
	ip := net.ParseIP(string(data))
	if ip == nil {
		return fmt.Errorf("invalid IP %q", string(data))
	}

	i.IP = ip
	return nil
}

// OpenAPISchemaType is used by the kube-openapi generator when constructing
// the OpenAPI spec of this type.
//
// See: https://github.com/kubernetes/kube-openapi/tree/master/pkg/generators
func (_ IP) OpenAPISchemaType() []string { return []string{"string"} }

// OpenAPISchemaFormat is used by the kube-openapi generator when constructing
// the OpenAPI spec of this type.
func (_ IP) OpenAPISchemaFormat() string { return "" }

func NewBigInt(bigInt *big.Int) BigInt {
	return BigInt{*bigInt}
}

// BigInt represents a big integer.
// +k8s:openapi-gen=true
type BigInt struct {
	BigInt big.Int `json:"-"`
}

func (b *BigInt) DeepCopy() *BigInt {
	return &BigInt{b.BigInt}
}

func (b *BigInt) DeepCopyInto(out *BigInt) {
	*out = *b
}

func (b *BigInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.BigInt.String())
}

func (b *BigInt) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	n := new(big.Int)
	n, ok := n.SetString(s, 10)
	if !ok {
		return fmt.Errorf("invalid big int: %q", string(data))
	}

	b.BigInt = *n
	return nil
}

// OpenAPISchemaType is used by the kube-openapi generator when constructing
// the OpenAPI spec of this type.
//
// See: https://github.com/kubernetes/kube-openapi/tree/master/pkg/generators
func (_ BigInt) OpenAPISchemaType() []string { return []string{"string"} }

// OpenAPISchemaFormat is used by the kube-openapi generator when constructing
// the OpenAPI spec of this type.
func (_ BigInt) OpenAPISchemaFormat() string { return "" }
