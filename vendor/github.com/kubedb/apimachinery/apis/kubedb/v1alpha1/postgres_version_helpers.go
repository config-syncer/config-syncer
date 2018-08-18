package v1alpha1

import (
	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

var _ ResourceInfo = &PostgresVersion{}

func (p PostgresVersion) ResourceShortCode() string {
	return ResourceCodePostgresVersion
}

func (p PostgresVersion) ResourceKind() string {
	return ResourceKindPostgresVersion
}

func (p PostgresVersion) ResourceSingular() string {
	return ResourceSingularPostgresVersion
}

func (p PostgresVersion) ResourcePlural() string {
	return ResourcePluralPostgresVersion
}

func (p PostgresVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePluralPostgresVersion,
		Singular:      ResourceSingularPostgresVersion,
		Kind:          ResourceKindPostgresVersion,
		ShortNames:    []string{ResourceCodePostgresVersion},
		Categories:    []string{"datastore", "kubedb", "appscode", "all"},
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
		SpecDefinitionName:      "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1.PostgresVersion",
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
				Name:     "DbImage",
				Type:     "string",
				JSONPath: ".spec.db.image",
			},
			{
				Name:     "ExporterImage",
				Type:     "string",
				JSONPath: ".spec.exporter.image",
			},
			{
				Name:     "ToolsImage",
				Type:     "string",
				JSONPath: ".spec.tools.image",
			},
			{
				Name:     "Age",
				Type:     "date",
				JSONPath: ".metadata.creationTimestamp",
			},
		},
	})
}
