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

var _ apis.ResourceInfo = &Postgres{}

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
	out := p.OffshootSelectors()
	out[meta_util.NameLabelKey] = ResourceSingularPostgres
	out[meta_util.VersionLabelKey] = string(p.Spec.Version)
	out[meta_util.InstanceLabelKey] = p.Name
	out[meta_util.ComponentLabelKey] = "database"
	out[meta_util.ManagedByLabelKey] = GenericKey
	return meta_util.FilterKeys(GenericKey, out, p.Labels)
}

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

// Snapshot service account name.
func (p Postgres) SnapshotSAName() string {
	return fmt.Sprintf("%v-snapshot", p.OffshootName())
}

type postgresApp struct {
	*Postgres
}

func (r postgresApp) Name() string {
	return r.Postgres.Name
}

func (r postgresApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularPostgres))
}

func (p Postgres) AppBindingMeta() appcat.AppBindingMeta {
	return &postgresApp{&p}
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
	return "/metrics"
}

func (p postgresStatsService) Scheme() string {
	return ""
}

func (p Postgres) StatsService() mona.StatsAccessor {
	return &postgresStatsService{&p}
}

func (p Postgres) StatsServiceLabels() map[string]string {
	lbl := meta_util.FilterKeys(GenericKey, p.OffshootSelectors(), p.Labels)
	lbl[LabelRole] = "stats"
	return lbl
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
		SpecDefinitionName:      "kubedb.dev/apimachinery/apis/kubedb/v1alpha1.Postgres",
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

func (p *Postgres) SetDefaults() {
	if p == nil {
		return
	}
	p.Spec.SetDefaults()

	if p.Spec.PodTemplate.Spec.ServiceAccountName == "" {
		p.Spec.PodTemplate.Spec.ServiceAccountName = p.OffshootName()
	}
}

func (p *PostgresSpec) SetDefaults() {
	if p == nil {
		return
	}

	// perform defaulting
	p.BackupSchedule.SetDefaults()

	if p.StorageType == "" {
		p.StorageType = StorageTypeDurable
	}
	if p.UpdateStrategy.Type == "" {
		p.UpdateStrategy.Type = apps.RollingUpdateStatefulSetStrategyType
	}
	if p.TerminationPolicy == "" {
		if p.StorageType == StorageTypeEphemeral {
			p.TerminationPolicy = TerminationPolicyDelete
		} else {
			p.TerminationPolicy = TerminationPolicyPause
		}
	}
	if p.Init != nil && p.Init.PostgresWAL != nil && p.Init.PostgresWAL.PITR != nil {
		pitr := p.Init.PostgresWAL.PITR

		if pitr.TargetInclusive == nil {
			pitr.TargetInclusive = types.BoolP(true)
		}

		p.Init.PostgresWAL.PITR = pitr
	}

	if p.LeaderElection == nil {
		// Default values: https://github.com/kubernetes/apiserver/blob/e85ad7b666fef0476185731329f4cff1536efff8/pkg/apis/config/v1alpha1/defaults.go#L26-L52
		p.LeaderElection = &LeaderElectionConfig{
			LeaseDurationSeconds: 15,
			RenewDeadlineSeconds: 10,
			RetryPeriodSeconds:   2,
		}
	}
}

func (e *PostgresSpec) GetSecrets() []string {
	if e == nil {
		return nil
	}

	var secrets []string
	if e.DatabaseSecret != nil {
		secrets = append(secrets, e.DatabaseSecret.SecretName)
	}
	return secrets
}
