package v1alpha1

import (
	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

var _ ResourceInfo = &MySQLVersion{}

func (p MySQLVersion) ResourceShortCode() string {
	return ResourceCodeMySQLVersion
}

func (p MySQLVersion) ResourceKind() string {
	return ResourceKindMySQLVersion
}

func (p MySQLVersion) ResourceSingular() string {
	return ResourceSingularMySQLVersion
}

func (p MySQLVersion) ResourcePlural() string {
	return ResourcePluralMySQLVersion
}

func (p MySQLVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePluralMySQLVersion,
		Singular:      ResourceSingularMySQLVersion,
		Kind:          ResourceKindMySQLVersion,
		ShortNames:    []string{ResourceCodeMySQLVersion},
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
		SpecDefinitionName:      "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1.MySQLVersion",
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
