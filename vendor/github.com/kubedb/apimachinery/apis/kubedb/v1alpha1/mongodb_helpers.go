package v1alpha1

import (
	"fmt"
	"strconv"

	"github.com/appscode/go/types"
	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	meta_util "github.com/appscode/kutil/meta"
	"github.com/kubedb/apimachinery/apis"
	"github.com/kubedb/apimachinery/apis/kubedb"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
)

var _ apis.ResourceInfo = &MongoDB{}

func (m MongoDB) OffshootName() string {
	return m.Name
}

func (m MongoDB) OffshootSelectors() map[string]string {
	return map[string]string{
		LabelDatabaseName: m.Name,
		LabelDatabaseKind: ResourceKindMongoDB,
	}
}

func (m MongoDB) OffshootLabels() map[string]string {
	return meta_util.FilterKeys(GenericKey, m.OffshootSelectors(), m.Labels)
}

func (m MongoDB) ResourceShortCode() string {
	return ResourceCodeMongoDB
}

func (m MongoDB) ResourceKind() string {
	return ResourceKindMongoDB
}

func (m MongoDB) ResourceSingular() string {
	return ResourceSingularMongoDB
}

func (m MongoDB) ResourcePlural() string {
	return ResourcePluralMongoDB
}

func (m MongoDB) ServiceName() string {
	return m.OffshootName()
}

func (m MongoDB) GoverningServiceName() string {
	return m.OffshootName() + "-gvr"
}

// HostAddress returns serviceName for standalone mongodb.
// and, for replica set = <replName>/<host1>,<host2>,<host3>
// Here, host1 = <pod-name>.<governing-serviceName>
// Governing service name is used for replica host because,
// we used governing service name as part of host while adding members
// to replicaset.
func (m MongoDB) HostAddress() string {
	host := m.ServiceName()
	if m.Spec.ReplicaSet != nil {
		host = m.Spec.ReplicaSet.Name + "/" + m.Name + "-0." + m.GoverningServiceName() + "." + m.Namespace + ".svc"
		for i := 1; i < int(types.Int32(m.Spec.Replicas)); i++ {
			host += "," + m.Name + "-" + strconv.Itoa(i) + "." + m.GoverningServiceName() + "." + m.Namespace + ".svc"
		}
	}
	return host
}

type mongoDBApp struct {
	*MongoDB
}

func (r mongoDBApp) Name() string {
	return r.MongoDB.Name
}

func (r mongoDBApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularMongoDB))
}

func (r MongoDB) AppBindingMeta() appcat.AppBindingMeta {
	return &mongoDBApp{&r}
}

type mongoDBStatsService struct {
	*MongoDB
}

func (m mongoDBStatsService) GetNamespace() string {
	return m.MongoDB.GetNamespace()
}

func (m mongoDBStatsService) ServiceName() string {
	return m.OffshootName() + "-stats"
}

func (m mongoDBStatsService) ServiceMonitorName() string {
	return fmt.Sprintf("kubedb-%s-%s", m.Namespace, m.Name)
}

func (m mongoDBStatsService) Path() string {
	return "/metrics"
}

func (m mongoDBStatsService) Scheme() string {
	return ""
}

func (m MongoDB) StatsService() mona.StatsAccessor {
	return &mongoDBStatsService{&m}
}

func (m *MongoDB) GetMonitoringVendor() string {
	if m.Spec.Monitor != nil {
		return m.Spec.Monitor.Agent.Vendor()
	}
	return ""
}

func (m MongoDB) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePluralMongoDB,
		Singular:      ResourceSingularMongoDB,
		Kind:          ResourceKindMongoDB,
		ShortNames:    []string{ResourceCodeMongoDB},
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
		SpecDefinitionName:      "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1.MongoDB",
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

func (m *MongoDB) SetDefaults() {
	if m == nil {
		return
	}
	m.Spec.SetDefaults()
}

func (m *MongoDBSpec) SetDefaults() {
	if m == nil {
		return
	}

	// migrate first to avoid incorrect defaulting
	m.BackupSchedule.SetDefaults()
	if m.DoNotPause {
		m.TerminationPolicy = TerminationPolicyDoNotTerminate
		m.DoNotPause = false
	}
	if len(m.NodeSelector) > 0 {
		m.PodTemplate.Spec.NodeSelector = m.NodeSelector
		m.NodeSelector = nil
	}
	if m.Resources != nil {
		m.PodTemplate.Spec.Resources = *m.Resources
		m.Resources = nil
	}
	if m.Affinity != nil {
		m.PodTemplate.Spec.Affinity = m.Affinity
		m.Affinity = nil
	}
	if len(m.SchedulerName) > 0 {
		m.PodTemplate.Spec.SchedulerName = m.SchedulerName
		m.SchedulerName = ""
	}
	if len(m.Tolerations) > 0 {
		m.PodTemplate.Spec.Tolerations = m.Tolerations
		m.Tolerations = nil
	}
	if len(m.ImagePullSecrets) > 0 {
		m.PodTemplate.Spec.ImagePullSecrets = m.ImagePullSecrets
		m.ImagePullSecrets = nil
	}

	// perform defaulting
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
	m.setDefaultProbes()
}

// setDefaultProbes sets defaults only when probe fields are nil.
// In operator, check if the value of probe fields is "{}".
// For "{}", ignore readinessprobe or livenessprobe in statefulset.
// ref: https://github.com/helm/charts/blob/345ba987722350ffde56ec34d2928c0b383940aa/stable/mongodb/templates/deployment-standalone.yaml#L93
func (m *MongoDBSpec) setDefaultProbes() {
	cmd := []string{
		"mongo",
		"--eval",
		"db.adminCommand('ping')",
	}
	if m.PodTemplate.Spec.LivenessProbe == nil {
		m.PodTemplate.Spec.LivenessProbe = &core.Probe{
			Handler: core.Handler{
				Exec: &core.ExecAction{
					Command: cmd,
				},
			},
			FailureThreshold: 3,
			PeriodSeconds:    10,
			SuccessThreshold: 1,
			TimeoutSeconds:   5,
		}
	}
	if m.PodTemplate.Spec.ReadinessProbe == nil {
		m.PodTemplate.Spec.ReadinessProbe = &core.Probe{
			Handler: core.Handler{
				Exec: &core.ExecAction{
					Command: cmd,
				},
			},
			FailureThreshold: 3,
			PeriodSeconds:    10,
			SuccessThreshold: 1,
			TimeoutSeconds:   1,
		}
	}
}

func (m *MongoDBSpec) GetSecrets() []string {
	if m == nil {
		return nil
	}

	var secrets []string
	if m.DatabaseSecret != nil {
		secrets = append(secrets, m.DatabaseSecret.SecretName)
	}
	if m.ReplicaSet != nil && m.ReplicaSet.KeyFile != nil {
		secrets = append(secrets, m.ReplicaSet.KeyFile.SecretName)
	}
	return secrets
}
