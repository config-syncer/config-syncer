package v1alpha1

import (
	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	core "k8s.io/api/core/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	crd_api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

func (d DormantDatabase) OffshootName() string {
	return d.Name
}

func (d DormantDatabase) ResourceShortCode() string {
	return ResourceCodeDormantDatabase
}

func (d DormantDatabase) ResourceKind() string {
	return ResourceKindDormantDatabase
}

func (d DormantDatabase) ResourceSingular() string {
	return ResourceSingularDormantDatabase
}

func (d DormantDatabase) ResourcePlural() string {
	return ResourcePluralDormantDatabase
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
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Version:       SchemeGroupVersion.Version,
		Plural:        ResourcePluralDormantDatabase,
		Singular:      ResourceSingularDormantDatabase,
		Kind:          ResourceKindDormantDatabase,
		ListKind:      ResourceKindDormantDatabase + "List",
		ShortNames:    []string{ResourceCodeDormantDatabase},
		ResourceScope: string(apiextensions.NamespaceScoped),
		Labels: crdutils.Labels{
			LabelsMap: map[string]string{"app": "kubedb"},
		},
		SpecDefinitionName:    "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1.DormantDatabase",
		EnableValidation:      true,
		GetOpenAPIDefinitions: GetOpenAPIDefinitions,
	})
}
