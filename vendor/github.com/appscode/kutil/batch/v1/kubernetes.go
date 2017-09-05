package v1

import (
	"errors"

	"k8s.io/apimachinery/pkg/runtime/schema"
	batchv1 "k8s.io/client-go/pkg/apis/batch/v1"
)

func GetGroupVersionKind(v interface{}) schema.GroupVersionKind {
	switch v.(type) {
	case *batchv1.Job:
		return batchv1.SchemeGroupVersion.WithKind("Job")
	default:
		return schema.GroupVersionKind{}
	}
}

func AssignTypeKind(v interface{}) error {
	switch u := v.(type) {
	case *batchv1.Job:
		u.APIVersion = batchv1.SchemeGroupVersion.String()
		u.Kind = "Job"
		return nil
	}
	return errors.New("Unknown api object type")
}
