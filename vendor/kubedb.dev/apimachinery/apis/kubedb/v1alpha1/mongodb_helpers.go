package v1alpha1

import (
	"fmt"
	"strconv"

	"github.com/appscode/go/types"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	crdutils "kmodules.xyz/client-go/apiextensions/v1beta1"
	v1 "kmodules.xyz/client-go/core/v1"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/apis/kubedb"
)

var _ apis.ResourceInfo = &MongoDB{}

const (
	MongoDBShardLabelKey  = "mongodb.kubedb.com/node.shard"
	MongoDBConfigLabelKey = "mongodb.kubedb.com/node.config"
	MongoDBMongosLabelKey = "mongodb.kubedb.com/node.mongos"
)

func (m MongoDB) OffshootName() string {
	return m.Name
}

func (m MongoDB) ShardNodeName(nodeNum int32) string {
	if m.Spec.ShardTopology == nil {
		return ""
	}
	shardName := fmt.Sprintf("%v-shard%v", m.OffshootName(), nodeNum)
	return m.Spec.ShardTopology.Shard.Prefix + shardName
}

func (m MongoDB) ConfigSvrNodeName() string {
	if m.Spec.ShardTopology == nil {
		return ""
	}
	configsvrName := fmt.Sprintf("%v-configsvr", m.OffshootName())
	return m.Spec.ShardTopology.ConfigServer.Prefix + configsvrName
}

func (m MongoDB) MongosNodeName() string {
	if m.Spec.ShardTopology == nil {
		return ""
	}
	mongosName := fmt.Sprintf("%v-mongos", m.OffshootName())
	return m.Spec.ShardTopology.Mongos.Prefix + mongosName
}

// RepSetName returns Replicaset name only for spec.replicaset
func (m MongoDB) RepSetName() string {
	if m.Spec.ReplicaSet == nil {
		return ""
	}
	return m.Spec.ReplicaSet.Name
}

func (m MongoDB) ShardRepSetName(nodeNum int32) string {
	repSetName := fmt.Sprintf("shard%v", nodeNum)
	if m.Spec.ShardTopology != nil && m.Spec.ShardTopology.Shard.Prefix != "" {
		repSetName = fmt.Sprintf("%v%v", m.Spec.ShardTopology.Shard.Prefix, nodeNum)
	}
	return repSetName
}

func (m MongoDB) ConfigSvrRepSetName() string {
	repSetName := fmt.Sprintf("cnfRepSet")
	if m.Spec.ShardTopology != nil && m.Spec.ShardTopology.ConfigServer.Prefix != "" {
		repSetName = m.Spec.ShardTopology.ConfigServer.Prefix
	}
	return repSetName
}

func (m MongoDB) OffshootSelectors() map[string]string {
	return map[string]string{
		LabelDatabaseName: m.Name,
		LabelDatabaseKind: ResourceKindMongoDB,
	}
}

func (m MongoDB) ShardSelectors(nodeNum int32) map[string]string {
	return v1.UpsertMap(m.OffshootSelectors(), map[string]string{
		MongoDBShardLabelKey: m.ShardNodeName(nodeNum),
	})
}

func (m MongoDB) ConfigSvrSelectors() map[string]string {
	return v1.UpsertMap(m.OffshootSelectors(), map[string]string{
		MongoDBConfigLabelKey: m.ConfigSvrNodeName(),
	})
}

func (m MongoDB) MongosSelectors() map[string]string {
	return v1.UpsertMap(m.OffshootSelectors(), map[string]string{
		MongoDBMongosLabelKey: m.MongosNodeName(),
	})
}

func (m MongoDB) OffshootLabels() map[string]string {
	out := m.OffshootSelectors()
	out[meta_util.NameLabelKey] = ResourceSingularMongoDB
	out[meta_util.VersionLabelKey] = string(m.Spec.Version)
	out[meta_util.InstanceLabelKey] = m.Name
	out[meta_util.ComponentLabelKey] = "database"
	out[meta_util.ManagedByLabelKey] = GenericKey
	return meta_util.FilterKeys(GenericKey, out, m.Labels)
}

func (m MongoDB) ShardLabels(nodeNum int32) map[string]string {
	return meta_util.FilterKeys(GenericKey, m.OffshootLabels(), m.ShardSelectors(nodeNum))
}

func (m MongoDB) ConfigSvrLabels() map[string]string {
	return meta_util.FilterKeys(GenericKey, m.OffshootLabels(), m.ConfigSvrSelectors())
}

func (m MongoDB) MongosLabels() map[string]string {
	return meta_util.FilterKeys(GenericKey, m.OffshootLabels(), m.MongosSelectors())
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

// Governing Service Name. Here, name parameter is either
// OffshootName, ShardNodeName or ConfigSvrNodeName
func (m MongoDB) GvrSvcName(name string) string {
	return name + "-gvr"
}

// Snapshot service account name.
func (m MongoDB) SnapshotSAName() string {
	return fmt.Sprintf("%v-snapshot", m.OffshootName())
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
		//host = m.Spec.ReplicaSet.Name + "/" + m.Name + "-0." + m.GvrSvcName(m.OffshootName()) + "." + m.Namespace + ".svc"
		host = fmt.Sprintf("%v/", m.RepSetName())
		for i := 0; i < int(types.Int32(m.Spec.Replicas)); i++ {
			if i != 0 {
				host += ","
			}
			host += fmt.Sprintf("%v-%v.%v.%v.svc", m.Name, strconv.Itoa(i), m.GvrSvcName(m.OffshootName()), m.Namespace)
		}
	}
	return host
}

// ShardDSN = <shardReplName>/<host1:port>,<host2:port>,<host3:port>
//// Here, host1 = <pod-name>.<governing-serviceName>.svc
func (m MongoDB) ShardDSN(nodeNum int32) string {
	if m.Spec.ShardTopology == nil {
		return ""
	}
	host := fmt.Sprintf("%v/", m.ShardRepSetName(nodeNum))
	for i := 0; i < int(m.Spec.ShardTopology.Shard.Replicas); i++ {
		//host += "," + m.ShardNodeName(nodeNum) + "-" + strconv.Itoa(i) + "." + m.GvrSvcName(m.ShardNodeName(nodeNum)) + "." + m.Namespace + ".svc"

		if i != 0 {
			host += ","
		}
		host += fmt.Sprintf("%v-%v.%v.%v.svc:%v", m.ShardNodeName(nodeNum), strconv.Itoa(i), m.GvrSvcName(m.ShardNodeName(nodeNum)), m.Namespace, MongoDBShardPort)
	}
	return host
}

// ConfigSvrDSN = <configSvrReplName>/<host1:port>,<host2:port>,<host3:port>
//// Here, host1 = <pod-name>.<governing-serviceName>.svc
func (m MongoDB) ConfigSvrDSN() string {
	if m.Spec.ShardTopology == nil {
		return ""
	}
	//	host := m.ConfigSvrRepSetName() + "/" + m.ConfigSvrNodeName() + "-0." + m.GvrSvcName(m.ConfigSvrNodeName()) + "." + m.Namespace + ".svc"
	host := fmt.Sprintf("%v/", m.ConfigSvrRepSetName())
	for i := 0; i < int(m.Spec.ShardTopology.ConfigServer.Replicas); i++ {
		if i != 0 {
			host += ","
		}
		host += fmt.Sprintf("%v-%v.%v.%v.svc:%v", m.ConfigSvrNodeName(), strconv.Itoa(i), m.GvrSvcName(m.ConfigSvrNodeName()), m.Namespace, MongoDBShardPort)
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

func (m MongoDB) AppBindingMeta() appcat.AppBindingMeta {
	return &mongoDBApp{&m}
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

func (m MongoDB) StatsServiceLabels() map[string]string {
	lbl := meta_util.FilterKeys(GenericKey, m.OffshootSelectors(), m.Labels)
	lbl[LabelRole] = "stats"
	return lbl
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
		SpecDefinitionName:      "kubedb.dev/apimachinery/apis/kubedb/v1alpha1.MongoDB",
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
	if m.Spec.ShardTopology != nil {
		if m.Spec.ShardTopology.ConfigServer.PodTemplate.Spec.ServiceAccountName == "" {
			m.Spec.ShardTopology.ConfigServer.PodTemplate.Spec.ServiceAccountName = m.OffshootName()
		}
		if m.Spec.ShardTopology.Mongos.PodTemplate.Spec.ServiceAccountName == "" {
			m.Spec.ShardTopology.Mongos.PodTemplate.Spec.ServiceAccountName = m.OffshootName()
		}
		if m.Spec.ShardTopology.Shard.PodTemplate.Spec.ServiceAccountName == "" {
			m.Spec.ShardTopology.Shard.PodTemplate.Spec.ServiceAccountName = m.OffshootName()
		}
	} else {
		if m.Spec.PodTemplate == nil {
			m.Spec.PodTemplate = new(ofst.PodTemplateSpec)
		}
		if m.Spec.PodTemplate.Spec.ServiceAccountName == "" {
			m.Spec.PodTemplate.Spec.ServiceAccountName = m.OffshootName()
		}
	}
}

func (m *MongoDBSpec) SetDefaults() {
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

	if m.SSLMode == "" {
		m.SSLMode = SSLModeDisabled
	}

	if (m.ReplicaSet != nil || m.ShardTopology != nil) && m.ClusterAuthMode == "" {
		if m.SSLMode == SSLModeDisabled {
			m.ClusterAuthMode = ClusterAuthModeKeyFile
		} else {
			m.ClusterAuthMode = ClusterAuthModeX509
		}
	}

	// required to upgrade operator from 0.11.0 to 0.12.0
	if m.ReplicaSet != nil && m.ReplicaSet.KeyFile != nil {
		if m.CertificateSecret == nil {
			m.CertificateSecret = m.ReplicaSet.KeyFile
		}
		m.ReplicaSet.KeyFile = nil
	}

	if m.ShardTopology != nil {
		if m.ShardTopology.Mongos.Strategy.Type == "" {
			m.ShardTopology.Mongos.Strategy.Type = apps.RollingUpdateDeploymentStrategyType
		}

		// set default probes
		m.setDefaultProbes(&m.ShardTopology.Shard.PodTemplate)
		m.setDefaultProbes(&m.ShardTopology.ConfigServer.PodTemplate)
		m.setDefaultProbes(&m.ShardTopology.Mongos.PodTemplate)
	} else {
		if m.Replicas == nil {
			m.Replicas = types.Int32P(1)
		}

		if m.PodTemplate == nil {
			m.PodTemplate = new(ofst.PodTemplateSpec)
		}
		// set default probes
		m.setDefaultProbes(m.PodTemplate)
	}

}

// setDefaultProbes sets defaults only when probe fields are nil.
// In operator, check if the value of probe fields is "{}".
// For "{}", ignore readinessprobe or livenessprobe in statefulset.
// ref: https://github.com/helm/charts/blob/345ba987722350ffde56ec34d2928c0b383940aa/stable/mongodb/templates/deployment-standalone.yaml#L93
func (m *MongoDBSpec) setDefaultProbes(podTemplate *ofst.PodTemplateSpec) {
	if podTemplate == nil {
		return
	}

	cmd := []string{
		"mongo",
		"--host=localhost",
		"--eval",
		"db.adminCommand('ping')",
	}

	if m.SSLMode == SSLModeRequireSSL {
		cmd = append(cmd, []string{
			"--ssl",
			"--sslCAFile=/data/configdb/tls.crt",
			"--sslPEMKeyFile=/data/configdb/mongo.pem",
		}...)
	}

	if podTemplate.Spec.LivenessProbe == nil {
		podTemplate.Spec.LivenessProbe = &core.Probe{
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
	if podTemplate.Spec.ReadinessProbe == nil {
		podTemplate.Spec.ReadinessProbe = &core.Probe{
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

// setSecurityContext will set default PodSecurityContext.
// These values will be applied only to newly created objects.
// These defaultings should not be applied to DBs or dormantDBs,
// that is managed by previous operators,
func (m *MongoDBSpec) SetSecurityContext(podTemplate *ofst.PodTemplateSpec) {
	if podTemplate == nil {
		return
	}
	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = new(core.PodSecurityContext)
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = types.Int64P(999)
	}
	if podTemplate.Spec.SecurityContext.RunAsNonRoot == nil {
		podTemplate.Spec.SecurityContext.RunAsNonRoot = types.BoolP(true)
	}
	if podTemplate.Spec.SecurityContext.RunAsUser == nil {
		podTemplate.Spec.SecurityContext.RunAsUser = types.Int64P(999)
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
	if m.CertificateSecret != nil {
		secrets = append(secrets, m.CertificateSecret.SecretName)
	}
	if m.ReplicaSet != nil && m.ReplicaSet.KeyFile != nil {
		secrets = append(secrets, m.ReplicaSet.KeyFile.SecretName)
	}
	return secrets
}
