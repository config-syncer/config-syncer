package v1alpha1

import (
	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

func (a ClusterAlert) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Version:       SchemeGroupVersion.Version,
		Plural:        ResourcePluralClusterAlert,
		Singular:      ResourceSingularClusterAlert,
		Kind:          ResourceKindClusterAlert,
		ListKind:      ResourceKindClusterAlert + "List",
		ShortNames:    []string{"ca"},
		ResourceScope: string(apiextensions.NamespaceScoped),
		Labels: crdutils.Labels{
			LabelsMap: map[string]string{
				"app": "searchlight",
			},
		},
		SpecDefinitionName:    "github.com/appscode/searchlight/apis/monitoring/v1alpha1.ClusterAlert",
		EnableValidation:      true,
		GetOpenAPIDefinitions: GetOpenAPIDefinitions,
	})
}

func (a NodeAlert) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Version:       SchemeGroupVersion.Version,
		Plural:        ResourcePluralNodeAlert,
		Singular:      ResourceSingularNodeAlert,
		Kind:          ResourceKindNodeAlert,
		ListKind:      ResourceKindNodeAlert + "List",
		ShortNames:    []string{"noa"},
		ResourceScope: string(apiextensions.NamespaceScoped),
		Labels: crdutils.Labels{
			LabelsMap: map[string]string{
				"app": "searchlight",
			},
		},
		SpecDefinitionName:    "github.com/appscode/searchlight/apis/monitoring/v1alpha1.NodeAlert",
		EnableValidation:      true,
		GetOpenAPIDefinitions: GetOpenAPIDefinitions,
	})
}

func (a PodAlert) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Version:       SchemeGroupVersion.Version,
		Plural:        ResourcePluralPodAlert,
		Singular:      ResourceSingularPodAlert,
		Kind:          ResourceKindPodAlert,
		ListKind:      ResourceKindPodAlert + "List",
		ShortNames:    []string{"poa"},
		ResourceScope: string(apiextensions.NamespaceScoped),
		Labels: crdutils.Labels{
			LabelsMap: map[string]string{
				"app": "searchlight",
			},
		},
		SpecDefinitionName:    "github.com/appscode/searchlight/apis/monitoring/v1alpha1.PodAlert",
		EnableValidation:      true,
		GetOpenAPIDefinitions: GetOpenAPIDefinitions,
	})
}

func (a Incident) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Version:       SchemeGroupVersion.Version,
		Plural:        ResourcePluralIncident,
		Singular:      ResourceSingularIncident,
		Kind:          ResourceKindIncident,
		ListKind:      ResourceKindIncident + "List",
		ResourceScope: string(apiextensions.NamespaceScoped),
		Labels: crdutils.Labels{
			LabelsMap: map[string]string{
				"app": "searchlight",
			},
		},
		SpecDefinitionName:    "github.com/appscode/searchlight/apis/monitoring/v1alpha1.Incident",
		EnableValidation:      true,
		GetOpenAPIDefinitions: GetOpenAPIDefinitions,
	})
}
