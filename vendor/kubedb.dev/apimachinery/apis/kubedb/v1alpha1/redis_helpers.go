package v1alpha1

import (
	"fmt"

	"github.com/appscode/go/types"
	apps "k8s.io/api/apps/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	crdutils "kmodules.xyz/client-go/apiextensions/v1beta1"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/apis/kubedb"
)

var _ apis.ResourceInfo = &Redis{}

func (r Redis) OffshootName() string {
	return r.Name
}

func (r Redis) OffshootSelectors() map[string]string {
	return map[string]string{
		LabelDatabaseName: r.Name,
		LabelDatabaseKind: ResourceKindRedis,
	}
}

func (r Redis) OffshootLabels() map[string]string {
	out := r.OffshootSelectors()
	out[meta_util.NameLabelKey] = ResourceSingularRedis
	out[meta_util.VersionLabelKey] = string(r.Spec.Version)
	out[meta_util.InstanceLabelKey] = r.Name
	out[meta_util.ComponentLabelKey] = "database"
	out[meta_util.ManagedByLabelKey] = GenericKey
	return meta_util.FilterKeys(GenericKey, out, r.Labels)
}

func (r Redis) ResourceShortCode() string {
	return ResourceCodeRedis
}

func (r Redis) ResourceKind() string {
	return ResourceKindRedis
}

func (r Redis) ResourceSingular() string {
	return ResourceSingularRedis
}

func (r Redis) ResourcePlural() string {
	return ResourcePluralRedis
}

func (r Redis) ServiceName() string {
	return r.OffshootName()
}

func (r Redis) ConfigMapName() string {
	return r.OffshootName()
}

func (r Redis) BaseNameForShard() string {
	return fmt.Sprintf("%s-shard", r.OffshootName())
}

func (r Redis) StatefulSetNameWithShard(i int) string {
	return fmt.Sprintf("%s%d", r.BaseNameForShard(), i)
}

type redisApp struct {
	*Redis
}

func (r redisApp) Name() string {
	return r.Redis.Name
}

func (r redisApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularRedis))
}

func (r Redis) AppBindingMeta() appcat.AppBindingMeta {
	return &redisApp{&r}
}

type redisStatsService struct {
	*Redis
}

func (r redisStatsService) GetNamespace() string {
	return r.Redis.GetNamespace()
}

func (r redisStatsService) ServiceName() string {
	return r.OffshootName() + "-stats"
}

func (r redisStatsService) ServiceMonitorName() string {
	return fmt.Sprintf("kubedb-%s-%s", r.Namespace, r.Name)
}

func (r redisStatsService) Path() string {
	return "/metrics"
}

func (r redisStatsService) Scheme() string {
	return ""
}

func (r Redis) StatsService() mona.StatsAccessor {
	return &redisStatsService{&r}
}

func (r Redis) StatsServiceLabels() map[string]string {
	lbl := meta_util.FilterKeys(GenericKey, r.OffshootSelectors(), r.Labels)
	lbl[LabelRole] = "stats"
	return lbl
}

func (r *Redis) GetMonitoringVendor() string {
	if r.Spec.Monitor != nil {
		return r.Spec.Monitor.Agent.Vendor()
	}
	return ""
}

func (r Redis) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePluralRedis,
		Singular:      ResourceSingularRedis,
		Kind:          ResourceKindRedis,
		ShortNames:    []string{ResourceCodeRedis},
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
		SpecDefinitionName:      "kubedb.dev/apimachinery/apis/kubedb/v1alpha1.Redis",
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

func (r *Redis) SetDefaults() {
	if r == nil {
		return
	}
	r.Spec.SetDefaults()

	if r.Spec.PodTemplate.Spec.ServiceAccountName == "" {
		r.Spec.PodTemplate.Spec.ServiceAccountName = r.OffshootName()
	}
}

func (r *RedisSpec) SetDefaults() {
	if r == nil {
		return
	}

	// perform defaulting
	if r.Mode == "" {
		r.Mode = RedisModeStandalone
	} else if r.Mode == RedisModeCluster {
		if r.Cluster == nil {
			r.Cluster = &RedisClusterSpec{}
		}
		if r.Cluster.Master == nil {
			r.Cluster.Master = types.Int32P(3)
		}
		if r.Cluster.Replicas == nil {
			r.Cluster.Replicas = types.Int32P(1)
		}
	}
	if r.StorageType == "" {
		r.StorageType = StorageTypeDurable
	}
	if r.UpdateStrategy.Type == "" {
		r.UpdateStrategy.Type = apps.RollingUpdateStatefulSetStrategyType
	}
	if r.TerminationPolicy == "" {
		if r.StorageType == StorageTypeEphemeral {
			r.TerminationPolicy = TerminationPolicyDelete
		} else {
			r.TerminationPolicy = TerminationPolicyPause
		}
	}
}

func (e *RedisSpec) GetSecrets() []string {
	return nil
}
