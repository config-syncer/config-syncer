package v1beta1

import (
	"errors"

	voyager "github.com/appscode/voyager/api"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func GetGroupVersionKind(v interface{}) schema.GroupVersionKind {
	switch v.(type) {
	case *voyager.Ingress:
		return voyager.V1beta1SchemeGroupVersion.WithKind("Ingress")
	case *voyager.Certificate:
		return voyager.V1beta1SchemeGroupVersion.WithKind("Certificate")
	default:
		return schema.GroupVersionKind{}
	}
}

func AssignTypeKind(v interface{}) error {
	switch u := v.(type) {
	case *voyager.Ingress:
		u.APIVersion = voyager.V1beta1SchemeGroupVersion.String()
		u.Kind = "Ingress"
		return nil
	case *voyager.Certificate:
		u.APIVersion = voyager.V1beta1SchemeGroupVersion.String()
		u.Kind = "Certificate"
		return nil
	}
	return errors.New("Unknown api object type")
}
