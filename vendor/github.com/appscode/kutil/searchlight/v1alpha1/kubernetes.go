package v1alpha1

import (
	"errors"

	searchlight "github.com/appscode/searchlight/api"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func GetGroupVersionKind(v interface{}) schema.GroupVersionKind {
	switch v.(type) {
	case *searchlight.ClusterAlert:
		return searchlight.V1alpha1SchemeGroupVersion.WithKind("ClusterAlert")
	case *searchlight.NodeAlert:
		return searchlight.V1alpha1SchemeGroupVersion.WithKind("NodeAlert")
	case *searchlight.PodAlert:
		return searchlight.V1alpha1SchemeGroupVersion.WithKind("PodAlert")
	default:
		return schema.GroupVersionKind{}
	}
}

func AssignTypeKind(v interface{}) error {
	switch u := v.(type) {
	case *searchlight.ClusterAlert:
		u.APIVersion = searchlight.V1alpha1SchemeGroupVersion.String()
		u.Kind = "ClusterAlert"
		return nil
	case *searchlight.NodeAlert:
		u.APIVersion = searchlight.V1alpha1SchemeGroupVersion.String()
		u.Kind = "NodeAlert"
		return nil
	case *searchlight.PodAlert:
		u.APIVersion = searchlight.V1alpha1SchemeGroupVersion.String()
		u.Kind = "PodAlert"
		return nil
	}
	return errors.New("Unknown api object type")
}
