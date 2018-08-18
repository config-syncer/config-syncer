package v1alpha1

import (
	"fmt"
	"reflect"

	"github.com/appscode/go/log"
	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	meta_util "github.com/appscode/kutil/meta"
	"github.com/golang/glog"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
)

func (m Memcached) OffshootName() string {
	return m.Name
}

func (m Memcached) OffshootSelectors() map[string]string {
	return map[string]string{
		LabelDatabaseKind: ResourceKindMemcached,
		LabelDatabaseName: m.Name,
	}
}

func (m Memcached) OffshootLabels() map[string]string {
	return meta_util.FilterKeys(GenericKey, m.OffshootSelectors(), m.Labels)
}

func (m Memcached) ResourceShortCode() string {
	return ResourceCodeMemcached
}

func (m Memcached) ResourceKind() string {
	return ResourceKindMemcached
}

func (m Memcached) ResourceSingular() string {
	return ResourceSingularMemcached
}

func (m Memcached) ResourcePlural() string {
	return ResourcePluralMemcached
}

func (m Memcached) ServiceName() string {
	return m.OffshootName()
}

type memcachedStatsService struct {
	*Memcached
}

func (m memcachedStatsService) GetNamespace() string {
	return m.Memcached.GetNamespace()
}

func (m memcachedStatsService) ServiceName() string {
	return m.OffshootName() + "-stats"
}

func (m memcachedStatsService) ServiceMonitorName() string {
	return fmt.Sprintf("kubedb-%s-%s", m.Namespace, m.Name)
}

func (m memcachedStatsService) Path() string {
	return fmt.Sprintf("/kubedb.com/v1alpha1/namespaces/%s/%s/%s/metrics", m.Namespace, m.ResourcePlural(), m.Name)
}

func (m memcachedStatsService) Scheme() string {
	return ""
}

func (m Memcached) StatsService() mona.StatsAccessor {
	return &memcachedStatsService{&m}
}

func (m *Memcached) GetMonitoringVendor() string {
	if m.Spec.Monitor != nil {
		return m.Spec.Monitor.Agent.Vendor()
	}
	return ""
}

func (m Memcached) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePluralMemcached,
		Singular:      ResourceSingularMemcached,
		Kind:          ResourceKindMemcached,
		ShortNames:    []string{ResourceCodeMemcached},
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
		SpecDefinitionName:      "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1.Memcached",
		EnableValidation:        true,
		GetOpenAPIDefinitions:   GetOpenAPIDefinitions,
		EnableStatusSubresource: EnableStatusSubresource,
		AdditionalPrinterColumns: []apiextensions.CustomResourceColumnDefinition{
			{
				Name:     "Version",
				Type:     "string",
				JSONPath: ".spec.version",
			},
			{
				Name:     "Status",
				Type:     "string",
				JSONPath: ".status.phase",
			},
			{
				Name:     "Age",
				Type:     "date",
				JSONPath: ".metadata.creationTimestamp",
			},
		},
	}, setNameSchema)
}

func (m *Memcached) Migrate() {
	if m == nil {
		return
	}
	m.Spec.Migrate()
}

func (m *MemcachedSpec) Migrate() {
	if m == nil {
		return
	}
	if len(m.NodeSelector) > 0 {
		m.PodTemplate.Spec.NodeSelector = m.NodeSelector
		m.NodeSelector = nil
	}
	if m.Resources != nil {
		m.PodTemplate.Spec.Resources = *m.Resources
		m.Resources = nil
	}
	if m.Affinity != nil {
		m.PodTemplate.Spec.Affinity = m.Affinity
		m.Affinity = nil
	}
	if len(m.SchedulerName) > 0 {
		m.PodTemplate.Spec.SchedulerName = m.SchedulerName
		m.SchedulerName = ""
	}
	if len(m.Tolerations) > 0 {
		m.PodTemplate.Spec.Tolerations = m.Tolerations
		m.Tolerations = nil
	}
	if len(m.ImagePullSecrets) > 0 {
		m.PodTemplate.Spec.ImagePullSecrets = m.ImagePullSecrets
		m.ImagePullSecrets = nil
	}
}

func (m *Memcached) AlreadyObserved(other *Memcached) bool {
	if m == nil {
		return other == nil
	}
	if other == nil { // && d != nil
		return false
	}
	if m == other {
		return true
	}

	var match bool

	if EnableStatusSubresource {
		match = m.Status.ObservedGeneration >= m.Generation
	} else {
		match = meta_util.Equal(m.Spec, other.Spec)
	}
	if match {
		match = reflect.DeepEqual(m.Labels, other.Labels)
	}
	if match {
		match = meta_util.EqualAnnotation(m.Annotations, other.Annotations)
	}

	if !match && bool(glog.V(log.LevelDebug)) {
		diff := meta_util.Diff(other, m)
		glog.V(log.LevelDebug).Infof("%s %s/%s has changed. Diff: %s", meta_util.GetKind(m), m.Namespace, m.Name, diff)
	}
	return match
}
