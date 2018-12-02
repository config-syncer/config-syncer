package v1alpha1

import (
	"fmt"

	"github.com/appscode/go/types"
	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	meta_util "github.com/appscode/kutil/meta"
	"github.com/kubedb/apimachinery/apis"
	"github.com/kubedb/apimachinery/apis/kubedb"
	apps "k8s.io/api/apps/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
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
	return meta_util.FilterKeys(GenericKey, r.OffshootSelectors(), r.Labels)
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
		SpecDefinitionName:      "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1.Redis",
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
}

func (r *RedisSpec) SetDefaults() {
	if r == nil {
		return
	}

	// migrate first to avoid incorrect defaulting
	if r.DoNotPause {
		r.TerminationPolicy = TerminationPolicyDoNotTerminate
		r.DoNotPause = false
	}
	if len(r.NodeSelector) > 0 {
		r.PodTemplate.Spec.NodeSelector = r.NodeSelector
		r.NodeSelector = nil
	}
	if r.Resources != nil {
		r.PodTemplate.Spec.Resources = *r.Resources
		r.Resources = nil
	}
	if r.Affinity != nil {
		r.PodTemplate.Spec.Affinity = r.Affinity
		r.Affinity = nil
	}
	if len(r.SchedulerName) > 0 {
		r.PodTemplate.Spec.SchedulerName = r.SchedulerName
		r.SchedulerName = ""
	}
	if len(r.Tolerations) > 0 {
		r.PodTemplate.Spec.Tolerations = r.Tolerations
		r.Tolerations = nil
	}
	if len(r.ImagePullSecrets) > 0 {
		r.PodTemplate.Spec.ImagePullSecrets = r.ImagePullSecrets
		r.ImagePullSecrets = nil
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
