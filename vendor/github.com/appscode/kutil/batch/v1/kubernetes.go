package v1

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/appscode/kutil"
	"k8s.io/apimachinery/pkg/runtime/schema"
	batchv1 "k8s.io/client-go/pkg/apis/batch/v1"
)

func GetGroupVersionKind(v interface{}) schema.GroupVersionKind {
	return batchv1.SchemeGroupVersion.WithKind(kutil.GetKind(v))
}

func AssignTypeKind(v interface{}) error {
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return fmt.Errorf("%v must be a pointer", v)
	}

	switch u := v.(type) {
	case *batchv1.Job:
		u.APIVersion = batchv1.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	}
	return errors.New("unknown api object type")
}
