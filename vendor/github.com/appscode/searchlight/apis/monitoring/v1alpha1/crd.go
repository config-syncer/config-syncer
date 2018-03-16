package v1alpha1

import (
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (a ClusterAlert) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return &apiextensions.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: ResourcePluralClusterAlert + "." + SchemeGroupVersion.Group,
			Labels: map[string]string{
				"app": "searchlight",
			},
		},
		Spec: apiextensions.CustomResourceDefinitionSpec{
			Group:   SchemeGroupVersion.Group,
			Version: SchemeGroupVersion.Version,
			Scope:   apiextensions.NamespaceScoped,
			Names: apiextensions.CustomResourceDefinitionNames{
				Plural:     ResourcePluralClusterAlert,
				Singular:   ResourceSingularClusterAlert,
				Kind:       ResourceKindClusterAlert,
				ShortNames: []string{"ca"},
			},
		},
	}
}

func (a NodeAlert) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return &apiextensions.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: ResourcePluralNodeAlert + "." + SchemeGroupVersion.Group,
			Labels: map[string]string{
				"app": "searchlight",
			},
		},
		Spec: apiextensions.CustomResourceDefinitionSpec{
			Group:   SchemeGroupVersion.Group,
			Version: SchemeGroupVersion.Version,
			Scope:   apiextensions.NamespaceScoped,
			Names: apiextensions.CustomResourceDefinitionNames{
				Plural:     ResourcePluralNodeAlert,
				Singular:   ResourceSingularNodeAlert,
				Kind:       ResourceKindNodeAlert,
				ShortNames: []string{"noa"},
			},
		},
	}
}

func (a PodAlert) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return &apiextensions.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: ResourcePluralPodAlert + "." + SchemeGroupVersion.Group,
			Labels: map[string]string{
				"app": "searchlight",
			},
		},
		Spec: apiextensions.CustomResourceDefinitionSpec{
			Group:   SchemeGroupVersion.Group,
			Version: SchemeGroupVersion.Version,
			Scope:   apiextensions.NamespaceScoped,
			Names: apiextensions.CustomResourceDefinitionNames{
				Plural:     ResourcePluralPodAlert,
				Singular:   ResourceSingularPodAlert,
				Kind:       ResourceKindPodAlert,
				ShortNames: []string{"poa"},
			},
		},
	}
}

func (a Incident) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return &apiextensions.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: ResourcePluralIncident + "." + SchemeGroupVersion.Group,
			Labels: map[string]string{
				"app": "searchlight",
			},
		},
		Spec: apiextensions.CustomResourceDefinitionSpec{
			Group:   SchemeGroupVersion.Group,
			Version: SchemeGroupVersion.Version,
			Scope:   apiextensions.NamespaceScoped,
			Names: apiextensions.CustomResourceDefinitionNames{
				Plural:   ResourcePluralIncident,
				Singular: ResourceSingularIncident,
				Kind:     ResourceKindIncident,
			},
		},
	}
}
