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
	out := m.OffshootSelectors()
	out[meta_util.NameLabelKey] = ResourceSingularMemcached
	out[meta_util.VersionLabelKey] = string(m.Spec.Version)
	out[meta_util.InstanceLabelKey] = m.Name
	out[meta_util.ComponentLabelKey] = "database"
	out[meta_util.ManagedByLabelKey] = GenericKey
	return meta_util.FilterKeys(GenericKey, out, m.Labels)
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

func (m Memcached) AppBindingMeta() appcat.AppBindingMeta {
	return &memcachedApp{&m}
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

func (m Memcached) StatsServiceLabels() map[string]string {
	lbl := meta_util.FilterKeys(GenericKey, m.OffshootSelectors(), m.Labels)
	lbl[LabelRole] = "stats"
	return lbl
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

	if m.Spec.PodTemplate.Spec.ServiceAccountName == "" {
		m.Spec.PodTemplate.Spec.ServiceAccountName = m.OffshootName()
	}
}

func (m *MemcachedSpec) SetDefaults() {
	if m == nil {
		return
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
