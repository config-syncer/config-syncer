package v1alpha1

import (
	"fmt"

	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	"github.com/kubedb/apimachinery/apis"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

func (r MongoDBRole) RoleName() string {
	cluster := "-"
	if r.ClusterName != "" {
		cluster = r.ClusterName
	}
	return fmt.Sprintf("k8s.%s.%s.%s", cluster, r.Namespace, r.Name)
}

func (r MongoDBRole) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourceMongoDBRoles,
		Singular:      ResourceMongoDBRole,
		Kind:          ResourceKindMongoDBRole,
		Categories:    []string{"datastore", "kubedb", "appscode", "all"},
		ResourceScope: string(apiextensions.NamespaceScoped),
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
		SpecDefinitionName:      "github.com/kubedb/apimachinery/apis/authorization/v1alpha1.MongoDBRole",
		EnableValidation:        true,
		GetOpenAPIDefinitions:   GetOpenAPIDefinitions,
		EnableStatusSubresource: apis.EnableStatusSubresource,
	})
}

func (r MongoDBRole) IsValid() error {
	return nil
}
