package v1alpha1

import (
	"fmt"

	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	"github.com/appscode/kutil/meta"
	meta_util "github.com/appscode/kutil/meta"
	apps "k8s.io/api/apps/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
)

func (e Elasticsearch) OffshootName() string {
	return e.Name
}

func (e Elasticsearch) OffshootSelectors() map[string]string {
	return map[string]string{
		LabelDatabaseKind: ResourceKindElasticsearch,
		LabelDatabaseName: e.Name,
	}
}

func (e Elasticsearch) OffshootLabels() map[string]string {
	return meta_util.FilterKeys(GenericKey, e.OffshootSelectors(), e.Labels)
}

var _ ResourceInfo = &Elasticsearch{}

func (e Elasticsearch) ResourceShortCode() string {
	return ResourceCodeElasticsearch
}

func (e Elasticsearch) ResourceKind() string {
	return ResourceKindElasticsearch
}

func (e Elasticsearch) ResourceSingular() string {
	return ResourceSingularElasticsearch
}

func (e Elasticsearch) ResourcePlural() string {
	return ResourcePluralElasticsearch
}

func (e Elasticsearch) ServiceName() string {
	return e.OffshootName()
}

func (e *Elasticsearch) MasterServiceName() string {
	return fmt.Sprintf("%v-master", e.ServiceName())
}

func (e *Elasticsearch) GetConnectionScheme() string {
	scheme := "http"
	if e.Spec.EnableSSL {
		scheme = "https"
	}
	return scheme
}

func (e *Elasticsearch) GetConnectionURL() string {
	return fmt.Sprintf("%v://%s.%s:%d", e.GetConnectionScheme(), e.OffshootName(), e.Namespace, ElasticsearchRestPort)
}

type elasticsearchStatsService struct {
	*Elasticsearch
}

func (e elasticsearchStatsService) GetNamespace() string {
	return e.Elasticsearch.GetNamespace()
}

func (e elasticsearchStatsService) ServiceName() string {
	return e.OffshootName() + "-stats"
}

func (e elasticsearchStatsService) ServiceMonitorName() string {
	return fmt.Sprintf("kubedb-%s-%s", e.Namespace, e.Name)
}

func (e elasticsearchStatsService) Path() string {
	return fmt.Sprintf("/kubedb.com/v1alpha1/namespaces/%s/%s/%s/metrics", e.Namespace, e.ResourcePlural(), e.Name)
}

func (e elasticsearchStatsService) Scheme() string {
	return ""
}

func (e Elasticsearch) StatsService() mona.StatsAccessor {
	return &elasticsearchStatsService{&e}
}

func (e *Elasticsearch) GetMonitoringVendor() string {
	if e.Spec.Monitor != nil {
		return e.Spec.Monitor.Agent.Vendor()
	}
	return ""
}

func (e Elasticsearch) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePluralElasticsearch,
		Singular:      ResourceSingularElasticsearch,
		Kind:          ResourceKindElasticsearch,
		ShortNames:    []string{ResourceCodeElasticsearch},
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
		SpecDefinitionName:      "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1.Elasticsearch",
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

func (e *Elasticsearch) SetDefaults() {
	if e == nil {
		return
	}
	e.Spec.SetDefaults()
}

func (e *ElasticsearchSpec) SetDefaults() {
	if e == nil {
		return
	}

	// migrate first to avoid incorrect defaulting
	e.BackupSchedule.SetDefaults()
	if len(e.NodeSelector) > 0 {
		e.PodTemplate.Spec.NodeSelector = e.NodeSelector
		e.NodeSelector = nil
	}
	if e.Resources != nil {
		e.PodTemplate.Spec.Resources = *e.Resources
		e.Resources = nil
	}
	if e.Affinity != nil {
		e.PodTemplate.Spec.Affinity = e.Affinity
		e.Affinity = nil
	}
	if len(e.SchedulerName) > 0 {
		e.PodTemplate.Spec.SchedulerName = e.SchedulerName
		e.SchedulerName = ""
	}
	if len(e.Tolerations) > 0 {
		e.PodTemplate.Spec.Tolerations = e.Tolerations
		e.Tolerations = nil
	}
	if len(e.ImagePullSecrets) > 0 {
		e.PodTemplate.Spec.ImagePullSecrets = e.ImagePullSecrets
		e.ImagePullSecrets = nil
	}

	// perform defaulting
	if e.StorageType == "" {
		e.StorageType = StorageTypeDurable
	}
	if e.UpdateStrategy.Type == "" {
		e.UpdateStrategy.Type = apps.RollingUpdateStatefulSetStrategyType
	}
	if e.TerminationPolicy == "" {
		e.TerminationPolicy = TerminationPolicyPause
	}
}

func (e *ElasticsearchSpec) GetSecrets() []string {
	if e == nil {
		return nil
	}

	var secrets []string
	if e.DatabaseSecret != nil {
		secrets = append(secrets, e.DatabaseSecret.SecretName)
	}
	if e.CertificateSecret != nil {
		secrets = append(secrets, e.CertificateSecret.SecretName)
	}
	return secrets
}

const (
	ESSearchGuardDisabled = ElasticsearchKey + "/searchguard-disabled"
)

func (e Elasticsearch) SearchGuardDisabled() bool {
	v, _ := meta.GetBoolValue(e.Annotations, ESSearchGuardDisabled)
	return v
}
