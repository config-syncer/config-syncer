package v1alpha1

import (
	"errors"

	stash "github.com/appscode/stash/api"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func GetGroupVersionKind(v interface{}) schema.GroupVersionKind {
	switch v.(type) {
	case *stash.Restic:
		return stash.V1alpha1SchemeGroupVersion.WithKind("Restic")
	default:
		return schema.GroupVersionKind{}
	}
}

func AssignTypeKind(v interface{}) error {
	switch u := v.(type) {
	case *stash.Restic:
		u.APIVersion = stash.V1alpha1SchemeGroupVersion.String()
		u.Kind = "Restic"
		return nil
	}
	return errors.New("Unknown api object type")
}
