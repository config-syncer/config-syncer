package v1alpha1

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/appscode/kutil"
	stash "github.com/appscode/stash/apis/stash/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func GetGroupVersionKind(v interface{}) schema.GroupVersionKind {
	return stash.SchemeGroupVersion.WithKind(kutil.GetKind(v))
}

func AssignTypeKind(v interface{}) error {
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return fmt.Errorf("%v must be a pointer", v)
	}

	switch u := v.(type) {
	case *stash.Restic:
		u.APIVersion = stash.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	}
	return errors.New("unknown api object type")
}
