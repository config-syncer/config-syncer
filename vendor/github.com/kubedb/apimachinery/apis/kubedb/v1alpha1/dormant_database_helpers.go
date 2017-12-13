package v1alpha1

import (
	core "k8s.io/api/core/v1"
	crd_api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (d DormantDatabase) OffshootName() string {
	return d.Name
}

func (d DormantDatabase) ResourceCode() string {
	return ResourceCodeDormantDatabase
}

func (d DormantDatabase) ResourceKind() string {
	return ResourceKindDormantDatabase
}

func (d DormantDatabase) ResourceName() string {
	return ResourceNameDormantDatabase
}

func (d DormantDatabase) ResourceType() string {
	return ResourceTypeDormantDatabase
}

func (d DormantDatabase) ObjectReference() *core.ObjectReference {
	return &core.ObjectReference{
		APIVersion:      SchemeGroupVersion.String(),
		Kind:            ResourceKindDormantDatabase,
		Namespace:       d.Namespace,
		Name:            d.Name,
		UID:             d.UID,
		ResourceVersion: d.ResourceVersion,
	}
}

func (d DormantDatabase) CustomResourceDefinition() *crd_api.CustomResourceDefinition {
	resourceName := ResourceTypeDormantDatabase + "." + SchemeGroupVersion.Group
	return &crd_api.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: resourceName,
			Labels: map[string]string{
				"app": "kubedb",
			},
		},
		Spec: crd_api.CustomResourceDefinitionSpec{
			Group:   SchemeGroupVersion.Group,
			Version: SchemeGroupVersion.Version,
			Scope:   crd_api.NamespaceScoped,
			Names: crd_api.CustomResourceDefinitionNames{
				Plural:     ResourceTypeDormantDatabase,
				Kind:       ResourceKindDormantDatabase,
				ShortNames: []string{ResourceCodeDormantDatabase},
			},
		},
	}
}
