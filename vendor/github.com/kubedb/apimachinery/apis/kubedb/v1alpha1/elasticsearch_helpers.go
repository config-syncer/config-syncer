package v1alpha1

import (
	"fmt"
	"reflect"

	"github.com/appscode/go/log"
	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	"github.com/appscode/kutil/meta"
	meta_util "github.com/appscode/kutil/meta"
	"github.com/golang/glog"
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

func (e *Elasticsearch) Migrate() {
	if e == nil {
		return
	}
	e.Spec.Migrate()
}

func (e *ElasticsearchSpec) Migrate() {
	if e == nil {
		return
	}
	e.BackupSchedule.Migrate()
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
}

const (
	ESSearchGuardDisabled = ElasticsearchKey + "/searchguard-disabled"
)

func (e Elasticsearch) SearchGuardDisabled() bool {
	v, _ := meta.GetBoolValue(e.Annotations, ESSearchGuardDisabled)
	return v
}

func (e *Elasticsearch) AlreadyObserved(other *Elasticsearch) bool {
	if e == nil {
		return other == nil
	}
	if other == nil { // && d != nil
		return false
	}
	if e == other {
		return true
	}

	var match bool

	if EnableStatusSubresource {
		match = e.Status.ObservedGeneration >= e.Generation
	} else {
		match = meta_util.Equal(e.Spec, other.Spec)
	}
	if match {
		match = reflect.DeepEqual(e.Labels, other.Labels)
	}
	if match {
		match = meta_util.EqualAnnotation(e.Annotations, other.Annotations)
	}

	if !match && bool(glog.V(log.LevelDebug)) {
		diff := meta_util.Diff(other, e)
		glog.V(log.LevelDebug).Infof("%s %s/%s has changed. Diff: %s", meta_util.GetKind(e), e.Namespace, e.Name, diff)
	}
	return match
}
