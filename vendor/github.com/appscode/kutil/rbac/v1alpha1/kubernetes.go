package v1alpha1

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/appscode/kutil"
	"k8s.io/apimachinery/pkg/runtime/schema"
	rbacv1alpha1 "k8s.io/client-go/pkg/apis/rbac/v1alpha1"
)

func GetGroupVersionKind(v interface{}) schema.GroupVersionKind {
	return rbacv1alpha1.SchemeGroupVersion.WithKind(kutil.GetKind(v))
}

func AssignTypeKind(v interface{}) error {
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return fmt.Errorf("%v must be a pointer", v)
	}

	switch u := v.(type) {
	case *rbacv1alpha1.Role:
		u.APIVersion = rbacv1alpha1.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	case *rbacv1alpha1.RoleBinding:
		u.APIVersion = rbacv1alpha1.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	case *rbacv1alpha1.ClusterRole:
		u.APIVersion = rbacv1alpha1.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	case *rbacv1alpha1.ClusterRoleBinding:
		u.APIVersion = rbacv1alpha1.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	}
	return errors.New("unknown api object type")
}
