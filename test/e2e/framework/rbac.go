package framework

import (
	api "github.com/appscode/kubed/apis/kubed/v1alpha1"
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	USER_ROLE_NAME = "appscode:kubed:view"
	USER_ANONYMOUS = "system:anonymous"
)

func (f *Framework) CreateUserRole() *rbac.ClusterRole {
	return &rbac.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: USER_ROLE_NAME,
			Labels: map[string]string{
				"rbac.authorization.k8s.io/aggregate-to-admin": "true",
				"rbac.authorization.k8s.io/aggregate-to-edit":  "true",
				"rbac.authorization.k8s.io/aggregate-to-view":  "true",
			},
		},
		Rules: []rbac.PolicyRule{
			{
				APIGroups: []string{api.GroupName},
				Resources: []string{"searchresults"},
				Verbs:     []string{"get"},
			},
		},
	}
}

func (f *Framework) UserRoleBinding() *rbac.ClusterRoleBinding {
	return &rbac.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: USER_ANONYMOUS,
		},
		Subjects: []rbac.Subject{
			{
				Kind:     rbac.UserKind,
				Name:     USER_ANONYMOUS,
				APIGroup: rbac.GroupName,
			},
		},
		RoleRef: rbac.RoleRef{
			Kind:     "ClusterRole",
			Name:     "admin",
			APIGroup: rbac.GroupName,
		},
	}
}
