package v1alpha1

import (
	"fmt"

	"github.com/kubedb/apimachinery/apis"
	"github.com/kubedb/apimachinery/apis/kubedb"
	apps "k8s.io/api/apps/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	crdutils "kmodules.xyz/client-go/apiextensions/v1beta1"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
)

var _ apis.ResourceInfo = &Memcached{}

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

type memcachedApp struct {
	*Memcached
}

func (r memcachedApp) Name() string {
	return r.Memcached.Name
}

func (r memcachedApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularMemcached))
}

func (r Memcached) AppBindingMeta() appcat.AppBindingMeta {
	return &memcachedApp{&r}
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
	return "/metrics"
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
		EnableStatusSubresource: apis.EnableStatusSubresource,
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
	}, apis.SetNameSchema)
}

func (m *Memcached) SetDefaults() {
	if m == nil {
		return
	}
	m.Spec.SetDefaults()
}

func (m *MemcachedSpec) SetDefaults() {
	if m == nil {
		return
	}

	// migrate first to avoid incorrect defaulting
	if m.DoNotPause {
		m.TerminationPolicy = TerminationPolicyDoNotTerminate
		m.DoNotPause = false
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

	// perform defaulting
	if m.UpdateStrategy.Type == "" {
		m.UpdateStrategy.Type = apps.RollingUpdateDeploymentStrategyType
	}
	if m.TerminationPolicy == "" {
		m.TerminationPolicy = TerminationPolicyPause
	}
}

func (e *MemcachedSpec) GetSecrets() []string {
	return nil
}
