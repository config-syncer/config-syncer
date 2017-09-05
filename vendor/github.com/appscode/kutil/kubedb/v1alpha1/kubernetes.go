package v1alpha1

import (
	"errors"

	kubedb "github.com/k8sdb/apimachinery/api"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func GetGroupVersionKind(v interface{}) schema.GroupVersionKind {
	switch v.(type) {
	case *kubedb.Postgres:
		return kubedb.V1alpha1SchemeGroupVersion.WithKind("Postgres")
	case *kubedb.Elasticsearch:
		return kubedb.V1alpha1SchemeGroupVersion.WithKind("Elasticsearch")
	case *kubedb.Snapshot:
		return kubedb.V1alpha1SchemeGroupVersion.WithKind("Snapshot")
	case *kubedb.DormantDatabase:
		return kubedb.V1alpha1SchemeGroupVersion.WithKind("DormantDatabase")
	default:
		return schema.GroupVersionKind{}
	}
}

func AssignTypeKind(v interface{}) error {
	switch u := v.(type) {
	case *kubedb.Postgres:
		u.APIVersion = kubedb.V1alpha1SchemeGroupVersion.String()
		u.Kind = "Postgres"
		return nil
	case *kubedb.Elasticsearch:
		u.APIVersion = kubedb.V1alpha1SchemeGroupVersion.String()
		u.Kind = "Elasticsearch"
		return nil
	case *kubedb.Snapshot:
		u.APIVersion = kubedb.V1alpha1SchemeGroupVersion.String()
		u.Kind = "Snapshot"
		return nil
	case *kubedb.DormantDatabase:
		u.APIVersion = kubedb.V1alpha1SchemeGroupVersion.String()
		u.Kind = "DormantDatabase"
		return nil
	}
	return errors.New("Unknown api object type")
}
