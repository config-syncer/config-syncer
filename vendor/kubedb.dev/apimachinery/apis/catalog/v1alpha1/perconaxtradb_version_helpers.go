package v1alpha1

import (
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	crdutils "kmodules.xyz/client-go/apiextensions/v1beta1"
	"kubedb.dev/apimachinery/apis"
)

var _ apis.ResourceInfo = &PerconaXtraDBVersion{}

func (p PerconaXtraDBVersion) ResourceShortCode() string {
	return ResourceCodePerconaXtraDBVersion
}

func (p PerconaXtraDBVersion) ResourceKind() string {
	return ResourceKindPerconaXtraDBVersion
}

func (p PerconaXtraDBVersion) ResourceSingular() string {
	return ResourceSingularPerconaXtraDBVersion
}

func (p PerconaXtraDBVersion) ResourcePlural() string {
	return ResourcePluralPerconaXtraDBVersion
}

func (p PerconaXtraDBVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePluralPerconaXtraDBVersion,
		Singular:      ResourceSingularPerconaXtraDBVersion,
		Kind:          ResourceKindPerconaXtraDBVersion,
		ShortNames:    []string{ResourceCodePerconaXtraDBVersion},
		Categories:    []string{"datastore", "kubedb", "appscode"},
		ResourceScope: string(apiextensions.ClusterScoped),
		Versions: []apiextensions.CustomResourceDefinitionVersion{
			{
				Name:    SchemeGroupVersion.Version,
				Served:  true,
				Storage: true,
			},
		},
		Labels: crdutils.Labels{
			LabelsMap: map[string]string{"app": "kubedb"},
		},
		SpecDefinitionName:      "kubedb.dev/apimachinery/apis/catalog/v1alpha1.PerconaXtraDBVersion",
		EnableValidation:        true,
		GetOpenAPIDefinitions:   GetOpenAPIDefinitions,
		EnableStatusSubresource: false,
		AdditionalPrinterColumns: []apiextensions.CustomResourceColumnDefinition{
			{
				Name:     "Version",
				Type:     "string",
				JSONPath: ".spec.version",
			},
			{
				Name:     "DB_IMAGE",
				Type:     "string",
				JSONPath: ".spec.db.image",
			},
			{
				Name:     "PROXYSQL_IMAGE",
				Type:     "string",
				JSONPath: ".spec.proxysql.image",
			},
			{
				Name:     "Deprecated",
				Type:     "boolean",
				JSONPath: ".spec.deprecated",
			},
			{
				Name:     "Age",
				Type:     "date",
				JSONPath: ".metadata.creationTimestamp",
			},
		},
	})
}
