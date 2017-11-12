package util

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/appscode/kutil"
	"github.com/k8sdb/apimachinery/apis/kubedb/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func GetGroupVersionKind(v interface{}) schema.GroupVersionKind {
	return v1alpha1.SchemeGroupVersion.WithKind(kutil.GetKind(v))
}

func AssignTypeKind(v interface{}) error {
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return fmt.Errorf("%v must be a pointer", v)
	}

	switch u := v.(type) {
	case *v1alpha1.Postgres:
		u.APIVersion = v1alpha1.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	case *v1alpha1.MySQL:
		u.APIVersion = v1alpha1.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	case *v1alpha1.Elasticsearch:
		u.APIVersion = v1alpha1.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	case *v1alpha1.Snapshot:
		u.APIVersion = v1alpha1.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	case *v1alpha1.DormantDatabase:
		u.APIVersion = v1alpha1.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	}
	return errors.New("unknown api object type")
}
