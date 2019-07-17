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

var _ apis.ResourceInfo = &Elasticsearch{}

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
	out := e.OffshootSelectors()
	out[meta_util.NameLabelKey] = ResourceSingularElasticsearch
	out[meta_util.VersionLabelKey] = string(e.Spec.Version)
	out[meta_util.InstanceLabelKey] = e.Name
	out[meta_util.ComponentLabelKey] = "database"
	out[meta_util.ManagedByLabelKey] = GenericKey
	return meta_util.FilterKeys(GenericKey, out, e.Labels)
}

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

// Snapshot service account name.
func (e Elasticsearch) SnapshotSAName() string {
	return fmt.Sprintf("%v-snapshot", e.OffshootName())
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

type elasticsearchApp struct {
	*Elasticsearch
}

func (r elasticsearchApp) Name() string {
	return r.Elasticsearch.Name
}

func (r elasticsearchApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularElasticsearch))
}

func (r Elasticsearch) AppBindingMeta() appcat.AppBindingMeta {
	return &elasticsearchApp{&r}
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
	return "/metrics"
}

func (e elasticsearchStatsService) Scheme() string {
	return ""
}

func (e Elasticsearch) StatsService() mona.StatsAccessor {
	return &elasticsearchStatsService{&e}
}

func (e Elasticsearch) StatsServiceLabels() map[string]string {
	lbl := meta_util.FilterKeys(GenericKey, e.OffshootSelectors(), e.Labels)
	lbl[LabelRole] = "stats"
	return lbl
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

func (e *Elasticsearch) SetDefaults() {
	if e == nil {
		return
	}
	e.Spec.SetDefaults()

	if e.Spec.PodTemplate.Spec.ServiceAccountName == "" {
		e.Spec.PodTemplate.Spec.ServiceAccountName = e.OffshootName()
	}
}

func (e *ElasticsearchSpec) SetDefaults() {
	if e == nil {
		return
	}

	// perform defaulting
	e.BackupSchedule.SetDefaults()

	if e.AuthPlugin == "" {
		e.AuthPlugin = ElasticsearchAuthPluginSearchGuard
	}
	if e.StorageType == "" {
		e.StorageType = StorageTypeDurable
	}
	if e.UpdateStrategy.Type == "" {
		e.UpdateStrategy.Type = apps.RollingUpdateStatefulSetStrategyType
	}
	if e.TerminationPolicy == "" {
		if e.StorageType == StorageTypeEphemeral {
			e.TerminationPolicy = TerminationPolicyDelete
		} else {
			e.TerminationPolicy = TerminationPolicyPause
		}
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
