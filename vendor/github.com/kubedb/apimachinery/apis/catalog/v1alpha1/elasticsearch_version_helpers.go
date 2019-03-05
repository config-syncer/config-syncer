package v1alpha1

import (
	"github.com/kubedb/apimachinery/apis"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	crdutils "kmodules.xyz/client-go/apiextensions/v1beta1"
)

var _ apis.ResourceInfo = &ElasticsearchVersion{}

func (p ElasticsearchVersion) ResourceShortCode() string {
	return ResourceCodeElasticsearchVersion
}

func (p ElasticsearchVersion) ResourceKind() string {
	return ResourceKindElasticsearchVersion
}

func (p ElasticsearchVersion) ResourceSingular() string {
	return ResourceSingularElasticsearchVersion
}

func (p ElasticsearchVersion) ResourcePlural() string {
	return ResourcePluralElasticsearchVersion
}

func (p ElasticsearchVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePluralElasticsearchVersion,
		Singular:      ResourceSingularElasticsearchVersion,
		Kind:          ResourceKindElasticsearchVersion,
		ShortNames:    []string{ResourceCodeElasticsearchVersion},
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
		SpecDefinitionName:      "github.com/kubedb/apimachinery/apis/catalog/v1alpha1.ElasticsearchVersion",
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
