package v1alpha1

import (
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (a ClusterAlert) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return &apiextensions.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: ResourceTypeClusterAlert + "." + SchemeGroupVersion.Group,
			Labels: map[string]string{
				"app": "searchlight",
			},
		},
		Spec: apiextensions.CustomResourceDefinitionSpec{
			Group:   SchemeGroupVersion.Group,
			Version: SchemeGroupVersion.Version,
			Scope:   apiextensions.NamespaceScoped,
			Names: apiextensions.CustomResourceDefinitionNames{
				Plural:     ResourceTypeClusterAlert,
				Kind:       ResourceKindClusterAlert,
				ShortNames: []string{"ca"},
			},
		},
	}
}

func (a NodeAlert) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return &apiextensions.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: ResourceTypeNodeAlert + "." + SchemeGroupVersion.Group,
			Labels: map[string]string{
				"app": "searchlight",
			},
		},
		Spec: apiextensions.CustomResourceDefinitionSpec{
			Group:   SchemeGroupVersion.Group,
			Version: SchemeGroupVersion.Version,
			Scope:   apiextensions.NamespaceScoped,
			Names: apiextensions.CustomResourceDefinitionNames{
				Plural:     ResourceTypeNodeAlert,
				Kind:       ResourceKindNodeAlert,
				ShortNames: []string{"noa"},
			},
		},
	}
}

func (a PodAlert) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return &apiextensions.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: ResourceTypePodAlert + "." + SchemeGroupVersion.Group,
			Labels: map[string]string{
				"app": "searchlight",
			},
		},
		Spec: apiextensions.CustomResourceDefinitionSpec{
			Group:   SchemeGroupVersion.Group,
			Version: SchemeGroupVersion.Version,
			Scope:   apiextensions.NamespaceScoped,
			Names: apiextensions.CustomResourceDefinitionNames{
				Plural:     ResourceTypePodAlert,
				Kind:       ResourceKindPodAlert,
				ShortNames: []string{"poa"},
			},
		},
	}
}
