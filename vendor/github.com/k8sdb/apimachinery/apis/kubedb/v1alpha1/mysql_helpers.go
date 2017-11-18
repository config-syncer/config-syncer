package v1alpha1

import (
	"fmt"
	"strings"

	"github.com/appscode/kutil/tools/monitoring/api"
	core "k8s.io/api/core/v1"
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
	annotations[MySQLDatabaseVersion] = string(m.Spec.Version)
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
