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

func (p Postgres) OffshootName() string {
	return p.Name
}

func (p Postgres) OffshootSelectors() map[string]string {
	return map[string]string{
		LabelDatabaseName: p.Name,
		LabelDatabaseKind: ResourceKindPostgres,
	}
}

func (p Postgres) OffshootLabels() map[string]string {
	return meta_util.FilterKeys(GenericKey, p.OffshootSelectors(), p.Labels)
}

var _ ResourceInfo = &Postgres{}

func (p Postgres) ResourceShortCode() string {
	return ResourceCodePostgres
}

func (p Postgres) ResourceKind() string {
	return ResourceKindPostgres
}

func (p Postgres) ResourceSingular() string {
	return ResourceSingularPostgres
}

func (p Postgres) ResourcePlural() string {
	return ResourcePluralPostgres
}

func (p Postgres) ServiceName() string {
	return p.OffshootName()
}

type postgresStatsService struct {
	*Postgres
}

func (p postgresStatsService) GetNamespace() string {
	return p.Postgres.GetNamespace()
}

func (p postgresStatsService) ServiceName() string {
	return p.OffshootName() + "-stats"
}

func (p postgresStatsService) ServiceMonitorName() string {
	return fmt.Sprintf("kubedb-%s-%s", p.Namespace, p.Name)
}

func (p postgresStatsService) Path() string {
	return fmt.Sprintf("/kubedb.com/v1alpha1/namespaces/%s/%s/%s/metrics", p.Namespace, p.ResourcePlural(), p.Name)
}

func (p postgresStatsService) Scheme() string {
	return ""
}

func (p Postgres) StatsService() mona.StatsAccessor {
	return &postgresStatsService{&p}
}

func (p *Postgres) GetMonitoringVendor() string {
	if p.Spec.Monitor != nil {
		return p.Spec.Monitor.Agent.Vendor()
	}
	return ""
}

func (p Postgres) ReplicasServiceName() string {
	return fmt.Sprintf("%v-replicas", p.Name)
}

func (p Postgres) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePluralPostgres,
		Singular:      ResourceSingularPostgres,
		Kind:          ResourceKindPostgres,
		ShortNames:    []string{ResourceCodePostgres},
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
		SpecDefinitionName:      "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1.Postgres",
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

func (p *Postgres) Migrate() {
	if p == nil {
		return
	}
	p.Spec.Migrate()
}

func (p *PostgresSpec) Migrate() {
	if p == nil {
		return
	}
	p.BackupSchedule.Migrate()
	if len(p.NodeSelector) > 0 {
		p.PodTemplate.Spec.NodeSelector = p.NodeSelector
		p.NodeSelector = nil
	}
	if p.Resources != nil {
		p.PodTemplate.Spec.Resources = *p.Resources
		p.Resources = nil
	}
	if p.Affinity != nil {
		p.PodTemplate.Spec.Affinity = p.Affinity
		p.Affinity = nil
	}
	if len(p.SchedulerName) > 0 {
		p.PodTemplate.Spec.SchedulerName = p.SchedulerName
		p.SchedulerName = ""
	}
	if len(p.Tolerations) > 0 {
		p.PodTemplate.Spec.Tolerations = p.Tolerations
		p.Tolerations = nil
	}
	if len(p.ImagePullSecrets) > 0 {
		p.PodTemplate.Spec.ImagePullSecrets = p.ImagePullSecrets
		p.ImagePullSecrets = nil
	}
}

func (p *Postgres) AlreadyObserved(other *Postgres) bool {
	if p == nil {
		return other == nil
	}
	if other == nil { // && d != nil
		return false
	}
	if p == other {
		return true
	}

	var match bool

	if EnableStatusSubresource {
		match = p.Status.ObservedGeneration >= p.Generation
	} else {
		match = meta_util.Equal(p.Spec, other.Spec)
	}
	if match {
		match = reflect.DeepEqual(p.Labels, other.Labels)
	}
	if match {
		match = meta_util.EqualAnnotation(p.Annotations, other.Annotations)
	}

	if !match && bool(glog.V(log.LevelDebug)) {
		diff := meta_util.Diff(other, p)
		glog.V(log.LevelDebug).Infof("%s %s/%s has changed. Diff: %s", meta_util.GetKind(p), p.Namespace, p.Name, diff)
	}
	return match
}
