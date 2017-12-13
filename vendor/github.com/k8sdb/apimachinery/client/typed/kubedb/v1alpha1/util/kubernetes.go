package util

import (
	"errors"

	"github.com/appscode/kutil/meta"
	"github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	"github.com/kubernetes/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func GetGroupVersionKind(v interface{}) schema.GroupVersionKind {
	return v1alpha1.SchemeGroupVersion.WithKind(meta.GetKind(v))
}

func AssignTypeKind(v interface{}) error {
	_, err := conversion.EnforcePtr(v)
	if err != nil {
		return err
	}

	switch u := v.(type) {
	case *v1alpha1.Postgres:
		u.APIVersion = v1alpha1.SchemeGroupVersion.String()
		u.Kind = meta.GetKind(v)
		return nil
	case *v1alpha1.MongoDB:
		u.APIVersion = v1alpha1.SchemeGroupVersion.String()
		u.Kind = meta.GetKind(v)
		return nil
	case *v1alpha1.MySQL:
		u.APIVersion = v1alpha1.SchemeGroupVersion.String()
		u.Kind = meta.GetKind(v)
		return nil
	case *v1alpha1.Elasticsearch:
		u.APIVersion = v1alpha1.SchemeGroupVersion.String()
		u.Kind = meta.GetKind(v)
		return nil
	case *v1alpha1.Redis:
		u.APIVersion = v1alpha1.SchemeGroupVersion.String()
		u.Kind = meta.GetKind(v)
		return nil
	case *v1alpha1.Memcached:
		u.APIVersion = v1alpha1.SchemeGroupVersion.String()
		u.Kind = meta.GetKind(v)
		return nil
	case *v1alpha1.Snapshot:
		u.APIVersion = v1alpha1.SchemeGroupVersion.String()
		u.Kind = meta.GetKind(v)
		return nil
	case *v1alpha1.DormantDatabase:
		u.APIVersion = v1alpha1.SchemeGroupVersion.String()
		u.Kind = meta.GetKind(v)
		return nil
	}
	return errors.New("unknown api object type")
}
