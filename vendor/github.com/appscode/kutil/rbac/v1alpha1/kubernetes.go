package v1alpha1

import (
	"errors"

	"k8s.io/apimachinery/pkg/runtime/schema"
	rbacv1alpha1 "k8s.io/client-go/pkg/apis/rbac/v1alpha1"
)

func GetGroupVersionKind(v interface{}) schema.GroupVersionKind {
	switch v.(type) {
	case *rbacv1alpha1.Role:
		return rbacv1alpha1.SchemeGroupVersion.WithKind("Role")
	case *rbacv1alpha1.RoleBinding:
		return rbacv1alpha1.SchemeGroupVersion.WithKind("RoleBinding")
	case *rbacv1alpha1.ClusterRole:
		return rbacv1alpha1.SchemeGroupVersion.WithKind("ClusterRole")
	case *rbacv1alpha1.ClusterRoleBinding:
		return rbacv1alpha1.SchemeGroupVersion.WithKind("ClusterRoleBinding")
	default:
		return schema.GroupVersionKind{}
	}
}

func AssignTypeKind(v interface{}) error {
	switch u := v.(type) {
	case *rbacv1alpha1.Role:
		u.APIVersion = rbacv1alpha1.SchemeGroupVersion.String()
		u.Kind = "Role"
		return nil
	case *rbacv1alpha1.RoleBinding:
		u.APIVersion = rbacv1alpha1.SchemeGroupVersion.String()
		u.Kind = "RoleBinding"
		return nil
	case *rbacv1alpha1.ClusterRole:
		u.APIVersion = rbacv1alpha1.SchemeGroupVersion.String()
		u.Kind = "ClusterRole"
		return nil
	case *rbacv1alpha1.ClusterRoleBinding:
		u.APIVersion = rbacv1alpha1.SchemeGroupVersion.String()
		u.Kind = "ClusterRoleBinding"
		return nil
	}
	return errors.New("Unknown api object type")
}
