package v1alpha1

import (
	"fmt"
	"strings"

	"github.com/appscode/kutil/tools/monitoring/api"
	core "k8s.io/api/core/v1"
)

func (r Memcached) OffshootName() string {
	return r.Name
}

func (r Memcached) OffshootLabels() map[string]string {
	return map[string]string{
		LabelDatabaseName: r.Name,
		LabelDatabaseKind: ResourceKindMemcached,
	}
}

func (r Memcached) DeploymentLabels() map[string]string {
	labels := r.OffshootLabels()
	for key, val := range r.Labels {
		if !strings.HasPrefix(key, GenericKey+"/") && !strings.HasPrefix(key, MemcachedKey+"/") {
			labels[key] = val
		}
	}
	return labels
}

func (r Memcached) DeploymentAnnotations() map[string]string {
	annotations := make(map[string]string)
	for key, val := range r.Annotations {
		if !strings.HasPrefix(key, GenericKey+"/") && !strings.HasPrefix(key, MemcachedKey+"/") {
			annotations[key] = val
		}
	}
	annotations[MemcachedDatabaseVersion] = string(r.Spec.Version)
	return annotations
}

func (r Memcached) ResourceCode() string {
	return ResourceCodeMemcached
}

func (r Memcached) ResourceKind() string {
	return ResourceKindMemcached
}

func (r Memcached) ResourceName() string {
	return ResourceNameMemcached
}

func (r Memcached) ResourceType() string {
	return ResourceTypeMemcached
}

func (s Memcached) ObjectReference() *core.ObjectReference {
	return &core.ObjectReference{
		APIVersion:      SchemeGroupVersion.String(),
		Kind:            ResourceKindMemcached,
		Namespace:       s.Namespace,
		Name:            s.Name,
		UID:             s.UID,
		ResourceVersion: s.ResourceVersion,
	}
}

func (m Memcached) ServiceName() string {
	return m.OffshootName()
}

func (m Memcached) ServiceMonitorName() string {
	return fmt.Sprintf("kubedb-%s-%s", m.Namespace, m.Name)
}

func (m Memcached) Path() string {
	return fmt.Sprintf("/kubedb.com/v1alpha1/namespaces/%s/%s/%s/metrics", m.Namespace, m.ResourceType(), m.Name)
}

func (m Memcached) Scheme() string {
	return ""
}

func (m *Memcached) StatsAccessor() api.StatsAccessor {
	return m
}
