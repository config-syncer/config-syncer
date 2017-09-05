package v1

import (
	"errors"

	"k8s.io/apimachinery/pkg/runtime/schema"
	storagev1 "k8s.io/client-go/pkg/apis/storage/v1"
)

func GetGroupVersionKind(v interface{}) schema.GroupVersionKind {
	switch v.(type) {
	case *storagev1.StorageClass:
		return storagev1.SchemeGroupVersion.WithKind("StorageClass")
	default:
		return schema.GroupVersionKind{}
	}
}

func AssignTypeKind(v interface{}) error {
	switch u := v.(type) {
	case *storagev1.StorageClass:
		u.APIVersion = storagev1.SchemeGroupVersion.String()
		u.Kind = "StorageClass"
		return nil
	}
	return errors.New("Unknown api object type")
}
