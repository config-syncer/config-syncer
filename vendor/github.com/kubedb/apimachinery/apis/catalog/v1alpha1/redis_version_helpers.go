package v1alpha1

import (
	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	"github.com/kubedb/apimachinery/apis"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

var _ apis.ResourceInfo = &RedisVersion{}

func (p RedisVersion) ResourceShortCode() string {
	return ResourceCodeRedisVersion
}

func (p RedisVersion) ResourceKind() string {
	return ResourceKindRedisVersion
}

func (p RedisVersion) ResourceSingular() string {
	return ResourceSingularRedisVersion
}

func (p RedisVersion) ResourcePlural() string {
	return ResourcePluralRedisVersion
}

func (p RedisVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePluralRedisVersion,
		Singular:      ResourceSingularRedisVersion,
		Kind:          ResourceKindRedisVersion,
		ShortNames:    []string{ResourceCodeRedisVersion},
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
		SpecDefinitionName:      "github.com/kubedb/apimachinery/apis/catalog/v1alpha1.RedisVersion",
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
