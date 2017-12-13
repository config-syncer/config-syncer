package v1alpha1

import (
	"fmt"
	"strings"

	"github.com/appscode/kutil/tools/monitoring/api"
	core "k8s.io/api/core/v1"
	crd_api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (m Memcached) CustomResourceDefinition() *crd_api.CustomResourceDefinition {
	resourceName := ResourceTypeMemcached + "." + SchemeGroupVersion.Group

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
				Plural:     ResourceTypeMemcached,
				Kind:       ResourceKindMemcached,
				ShortNames: []string{ResourceCodeMemcached},
			},
		},
	}
}
