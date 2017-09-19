package v1alpha1

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/appscode/kutil"
	searchlight "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func GetGroupVersionKind(v interface{}) schema.GroupVersionKind {
	return searchlight.SchemeGroupVersion.WithKind(kutil.GetKind(v))
}

func AssignTypeKind(v interface{}) error {
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return fmt.Errorf("%v must be a pointer", v)
	}

	switch u := v.(type) {
	case *searchlight.ClusterAlert:
		u.APIVersion = searchlight.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	case *searchlight.NodeAlert:
		u.APIVersion = searchlight.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	case *searchlight.PodAlert:
		u.APIVersion = searchlight.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	}
	return errors.New("unknown api object type")
}
