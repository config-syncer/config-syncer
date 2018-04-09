package v1alpha1

import (
	"fmt"
	"strings"

	"github.com/appscode/kube-mon/api"
	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	crd_api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

func (m MySQL) OffshootName() string {
	return m.Name
}

func (m MySQL) OffshootLabels() map[string]string {
	return map[string]string{
		LabelDatabaseName: m.Name,
		LabelDatabaseKind: ResourceKindMySQL,
	}
}

func (m MySQL) StatefulSetLabels() map[string]string {
	labels := m.OffshootLabels()
	for key, val := range m.Labels {
		if !strings.HasPrefix(key, GenericKey+"/") && !strings.HasPrefix(key, MySQLKey+"/") {
			labels[key] = val
		}
	}
	return labels
}

func (m MySQL) StatefulSetAnnotations() map[string]string {
	annotations := make(map[string]string)
	for key, val := range m.Annotations {
		if !strings.HasPrefix(key, GenericKey+"/") && !strings.HasPrefix(key, MySQLKey+"/") {
			annotations[key] = val
		}
	}
	return annotations
}

var _ ResourceInfo = &MySQL{}

func (m MySQL) ResourceShortCode() string {
	return ResourceCodeMySQL
}

func (m MySQL) ResourceKind() string {
	return ResourceKindMySQL
}

func (m MySQL) ResourceSingular() string {
	return ResourceSingularMySQL
}

func (m MySQL) ResourcePlural() string {
	return ResourcePluralMySQL
}

func (m MySQL) ServiceName() string {
	return m.OffshootName()
}

func (m MySQL) ServiceMonitorName() string {
	return fmt.Sprintf("kubedb-%s-%s", m.Namespace, m.Name)
}

func (m MySQL) Path() string {
	return fmt.Sprintf("/kubedb.com/v1alpha1/namespaces/%s/%s/%s/metrics", m.Namespace, m.ResourcePlural(), m.Name)
}

func (m MySQL) Scheme() string {
	return ""
}

func (m *MySQL) StatsAccessor() api.StatsAccessor {
	return m
}

func (m *MySQL) GetMonitoringVendor() string {
	if m.Spec.Monitor != nil {
		return m.Spec.Monitor.Agent.Vendor()
	}
	return ""
}

func (m MySQL) CustomResourceDefinition() *crd_api.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Version:       SchemeGroupVersion.Version,
		Plural:        ResourcePluralMySQL,
		Singular:      ResourceSingularMySQL,
		Kind:          ResourceKindMySQL,
		ShortNames:    []string{ResourceCodeMySQL},
		ResourceScope: string(apiextensions.NamespaceScoped),
		Labels: crdutils.Labels{
			LabelsMap: map[string]string{"app": "kubedb"},
		},
		SpecDefinitionName:    "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1.MySQL",
		EnableValidation:      true,
		GetOpenAPIDefinitions: GetOpenAPIDefinitions,
	})
}
