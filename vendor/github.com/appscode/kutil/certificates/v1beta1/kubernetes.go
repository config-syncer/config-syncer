package v1beta1

import (
	"errors"

	"k8s.io/apimachinery/pkg/runtime/schema"
	certificates "k8s.io/client-go/pkg/apis/certificates/v1beta1"
)

func GetGroupVersionKind(v interface{}) schema.GroupVersionKind {
	switch v.(type) {
	case *certificates.CertificateSigningRequest:
		return certificates.SchemeGroupVersion.WithKind("CertificateSigningRequest")
	default:
		return schema.GroupVersionKind{}
	}
}

func AssignTypeKind(v interface{}) error {
	switch u := v.(type) {
	case *certificates.CertificateSigningRequest:
		u.APIVersion = certificates.SchemeGroupVersion.String()
		u.Kind = "CertificateSigningRequest"
		return nil
	}
	return errors.New("Unknown api object type")
}
