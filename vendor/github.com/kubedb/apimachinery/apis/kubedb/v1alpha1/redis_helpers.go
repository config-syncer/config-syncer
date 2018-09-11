package v1alpha1

import (
	"fmt"

	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	meta_util "github.com/appscode/kutil/meta"
	apps "k8s.io/api/apps/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
)

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
	return fmt.Sprintf("/kubedb.com/v1alpha1/namespaces/%s/%s/%s/metrics", r.Namespace, r.ResourcePlural(), r.Name)
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
	if r.StorageType == "" {
		r.StorageType = StorageTypeDurable
	}
	if r.UpdateStrategy.Type == "" {
		r.UpdateStrategy.Type = apps.RollingUpdateStatefulSetStrategyType
	}
	if r.TerminationPolicy == "" {
		r.TerminationPolicy = TerminationPolicyPause
	}
}

func (e *RedisSpec) GetSecrets() []string {
	return nil
}
