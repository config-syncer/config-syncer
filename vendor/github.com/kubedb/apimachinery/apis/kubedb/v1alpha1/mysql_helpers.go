package v1alpha1

import (
	"fmt"
	"strings"

	"github.com/appscode/kube-mon/api"
	core "k8s.io/api/core/v1"
	crd_api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (m MySQL) ResourceCode() string {
	return ResourceCodeMySQL
}

func (m MySQL) ResourceKind() string {
	return ResourceKindMySQL
}

func (m MySQL) ResourceName() string {
	return ResourceNameMySQL
}

func (m MySQL) ResourceType() string {
	return ResourceTypeMySQL
}

func (m MySQL) ObjectReference() *core.ObjectReference {
	return &core.ObjectReference{
		APIVersion:      SchemeGroupVersion.String(),
		Kind:            m.ResourceKind(),
		Namespace:       m.Namespace,
		Name:            m.Name,
		UID:             m.UID,
		ResourceVersion: m.ResourceVersion,
	}
}

func (m MySQL) ServiceName() string {
	return m.OffshootName()
}

func (m MySQL) ServiceMonitorName() string {
	return fmt.Sprintf("kubedb-%s-%s", m.Namespace, m.Name)
}

func (m MySQL) Path() string {
	return fmt.Sprintf("/kubedb.com/v1alpha1/namespaces/%s/%s/%s/metrics", m.Namespace, m.ResourceType(), m.Name)
}

func (m MySQL) Scheme() string {
	return ""
}

func (m *MySQL) StatsAccessor() api.StatsAccessor {
	return m
}

func (m MySQL) CustomResourceDefinition() *crd_api.CustomResourceDefinition {
	resourceName := ResourceTypeMySQL + "." + SchemeGroupVersion.Group

	return &crd_api.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: resourceName,
			Labels: map[string]string{
				"app": "kubedb",
			},
		},
		Spec: crd_api.CustomResourceDefinitionSpec{
			Group:   SchemeGroupVersion.Group,
			Version: SchemeGroupVersion.Version,
			Scope:   crd_api.NamespaceScoped,
			Names: crd_api.CustomResourceDefinitionNames{
				Plural:     ResourceTypeMySQL,
				Kind:       ResourceKindMySQL,
				ShortNames: []string{ResourceCodeMySQL},
			},
		},
	}
}
