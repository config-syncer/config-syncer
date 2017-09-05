package v1beta1

import (
	"errors"

	"k8s.io/apimachinery/pkg/runtime/schema"
	storagev1 "k8s.io/client-go/pkg/apis/storage/v1"
	storagev1beta1 "k8s.io/client-go/pkg/apis/storage/v1beta1"
)

func GetGroupVersionKind(v interface{}) schema.GroupVersionKind {
	switch v.(type) {
	case *storagev1beta1.StorageClass:
		return storagev1beta1.SchemeGroupVersion.WithKind("StorageClass")
	case *storagev1.StorageClass:
		return storagev1.SchemeGroupVersion.WithKind("StorageClass")
	default:
		return schema.GroupVersionKind{}
	}
}

func AssignTypeKind(v interface{}) error {
	switch u := v.(type) {
	case *storagev1beta1.StorageClass:
		u.APIVersion = storagev1beta1.SchemeGroupVersion.String()
		u.Kind = "StorageClass"
		return nil
	case *storagev1.StorageClass:
		u.APIVersion = storagev1.SchemeGroupVersion.String()
		u.Kind = "StorageClass"
		return nil
	}
	return errors.New("Unknown api object type")
}
