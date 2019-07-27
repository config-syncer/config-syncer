package v1alpha1

import (
	"fmt"

	"github.com/appscode/go/types"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	crdutils "kmodules.xyz/client-go/apiextensions/v1beta1"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/apis/kubedb"
)

var _ apis.ResourceInfo = &MySQL{}

func (m MySQL) OffshootName() string {
	return m.Name
}

func (m MySQL) OffshootSelectors() map[string]string {
	return map[string]string{
		LabelDatabaseName: m.Name,
		LabelDatabaseKind: ResourceKindMySQL,
	}
}

func (m MySQL) OffshootLabels() map[string]string {
	out := m.OffshootSelectors()
	out[meta_util.NameLabelKey] = ResourceSingularMySQL
	out[meta_util.VersionLabelKey] = string(m.Spec.Version)
	out[meta_util.InstanceLabelKey] = m.Name
	out[meta_util.ComponentLabelKey] = "database"
	out[meta_util.ManagedByLabelKey] = GenericKey
	return meta_util.FilterKeys(GenericKey, out, m.Labels)
}

func (m MySQL) ResourceShortCode() string {
	return ResourceCodeMySQL
}

func (m MySQL) ResourceKind() string {
	return ResourceKindMySQL
}

func (m MySQL) ResourceSingular() string {
	return ResourceSingularMySQL
}

func (m MySQL) ResourcePlural() string {
	return ResourcePluralMySQL
}

func (m MySQL) ServiceName() string {
	return m.OffshootName()
}

func (m MySQL) GoverningServiceName() string {
	return m.OffshootName() + "-gvr"
}

// Snapshot service account name.
func (m MySQL) SnapshotSAName() string {
	return fmt.Sprintf("%v-snapshot", m.OffshootName())
}

type mysqlApp struct {
	*MySQL
}

func (r mysqlApp) Name() string {
	return r.MySQL.Name
}

func (r mysqlApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularMySQL))
}

func (m MySQL) AppBindingMeta() appcat.AppBindingMeta {
	return &mysqlApp{&m}
}

type mysqlStatsService struct {
	*MySQL
}

func (m mysqlStatsService) GetNamespace() string {
	return m.MySQL.GetNamespace()
}

func (m mysqlStatsService) ServiceName() string {
	return m.OffshootName() + "-stats"
}

func (m mysqlStatsService) ServiceMonitorName() string {
	return fmt.Sprintf("kubedb-%s-%s", m.Namespace, m.Name)
}

func (m mysqlStatsService) Path() string {
	return "/metrics"
}

func (m mysqlStatsService) Scheme() string {
	return ""
}

func (m MySQL) StatsService() mona.StatsAccessor {
	return &mysqlStatsService{&m}
}

func (m MySQL) StatsServiceLabels() map[string]string {
	lbl := meta_util.FilterKeys(GenericKey, m.OffshootSelectors(), m.Labels)
	lbl[LabelRole] = "stats"
	return lbl
}

func (m *MySQL) GetMonitoringVendor() string {
	if m.Spec.Monitor != nil {
		return m.Spec.Monitor.Agent.Vendor()
	}
	return ""
}

func (m MySQL) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePluralMySQL,
		Singular:      ResourceSingularMySQL,
		Kind:          ResourceKindMySQL,
		ShortNames:    []string{ResourceCodeMySQL},
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
		SpecDefinitionName:      "kubedb.dev/apimachinery/apis/kubedb/v1alpha1.MySQL",
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

func (m *MySQL) SetDefaults() {
	if m == nil {
		return
	}
	m.Spec.SetDefaults()

	if m.Spec.PodTemplate.Spec.ServiceAccountName == "" {
		m.Spec.PodTemplate.Spec.ServiceAccountName = m.OffshootName()
	}
}

func (m *MySQLSpec) SetDefaults() {
	if m == nil {
		return
	}

	// perform defaulting
	m.BackupSchedule.SetDefaults()

	if m.StorageType == "" {
		m.StorageType = StorageTypeDurable
	}
	if m.UpdateStrategy.Type == "" {
		m.UpdateStrategy.Type = apps.RollingUpdateStatefulSetStrategyType
	}
	if m.TerminationPolicy == "" {
		if m.StorageType == StorageTypeEphemeral {
			m.TerminationPolicy = TerminationPolicyDelete
		} else {
			m.TerminationPolicy = TerminationPolicyPause
		}
	}

	if m.Topology != nil && m.Topology.Mode != nil && *m.Topology.Mode == MySQLClusterModeGroup {
		if m.Replicas == nil {
			m.Replicas = types.Int32P(MySQLDefaultGroupSize)
		}
		m.setDefaultProbes()
	} else {
		if m.Replicas == nil {
			m.Replicas = types.Int32P(1)
		}
	}
}

// setDefaultProbes sets defaults only when probe fields are nil.
// In operator, check if the value of probe fields is "{}".
// For "{}", ignore readinessprobe or livenessprobe in statefulset.
// Ref: https://github.com/mattlord/Docker-InnoDB-Cluster/blob/master/healthcheck.sh#L10
func (m *MySQLSpec) setDefaultProbes() {
	probe := &core.Probe{
		Handler: core.Handler{
			Exec: &core.ExecAction{
				Command: []string{
					"bash",
					"-c",
					`
export MYSQL_PWD=${MYSQL_ROOT_PASSWORD}
mysql -h localhost -nsLNE -e "select member_state from performance_schema.replication_group_members where member_id=@@server_uuid;" 2>/dev/null | grep -v "*" | egrep -v "ERROR|OFFLINE"
`,
				},
			},
		},
		InitialDelaySeconds: 30,
		PeriodSeconds:       5,
	}

	if m.PodTemplate.Spec.LivenessProbe == nil {
		m.PodTemplate.Spec.LivenessProbe = probe
	}
	if m.PodTemplate.Spec.ReadinessProbe == nil {
		m.PodTemplate.Spec.ReadinessProbe = probe
	}
}

func (e *MySQLSpec) GetSecrets() []string {
	if e == nil {
		return nil
	}

	var secrets []string
	if e.DatabaseSecret != nil {
		secrets = append(secrets, e.DatabaseSecret.SecretName)
	}
	return secrets
}
