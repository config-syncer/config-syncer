package v1alpha1

import (
	"github.com/appscode/go/encoding/json/types"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

const (
	ResourceCodeMongoDB     = "mg"
	ResourceKindMongoDB     = "MongoDB"
	ResourceSingularMongoDB = "mongodb"
	ResourcePluralMongoDB   = "mongodbs"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MongoDB defines a MongoDB database.
type MongoDB struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MongoDBSpec   `json:"spec,omitempty"`
	Status            MongoDBStatus `json:"status,omitempty"`
}

type MongoDBSpec struct {
	// Version of MongoDB to be deployed.
	Version types.StrYo `json:"version"`

	// Number of instances to deploy for a MongoDB database.
	Replicas *int32 `json:"replicas,omitempty"`

	// MongoDB replica set
	ReplicaSet *MongoDBReplicaSet `json:"replicaSet,omitempty"`

	// MongoDB sharding topology.
	ShardTopology *MongoDBShardingTopology `json:"shardTopology,omitempty"`

	// StorageType can be durable (default) or ephemeral
	StorageType StorageType `json:"storageType,omitempty"`

	// Storage spec to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// Database authentication secret
	DatabaseSecret *core.SecretVolumeSource `json:"databaseSecret,omitempty"`

	// Secret for KeyFile or SSL certificates. Contains `tls.pem` or keyfile `key.txt` depending on enableSSL.
	// Currently SSL support is not enabled.
	CertificateSecret *core.SecretVolumeSource `json:"certificateSecret,omitempty"`

	// Init is used to initialize database
	// +optional
	Init *InitSpec `json:"init,omitempty"`

	// BackupSchedule spec to specify how database backup will be taken
	// +optional
	BackupSchedule *BackupScheduleSpec `json:"backupSchedule,omitempty"`

	// Monitor is used monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`

	// ConfigSource is an optional field to provide custom configuration file for database (i.e mongod.cnf).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	ConfigSource *core.VolumeSource `json:"configSource,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate *ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

	// ServiceTemplate is an optional configuration for service used to expose database
	// +optional
	ServiceTemplate ofst.ServiceTemplateSpec `json:"serviceTemplate,omitempty"`

	// updateStrategy indicates the StatefulSetUpdateStrategy that will be
	// employed to update Pods in the StatefulSet when a revision is made to
	// Template.
	UpdateStrategy apps.StatefulSetUpdateStrategy `json:"updateStrategy,omitempty"`

	// TerminationPolicy controls the delete operation for database
	// +optional
	TerminationPolicy TerminationPolicy `json:"terminationPolicy,omitempty"`
}

type MongoDBReplicaSet struct {
	// Name of replicaset
	Name string `json:"name"`

	// Deprecated: Use spec.certificateSecret
	KeyFile *core.SecretVolumeSource `json:"keyFile,omitempty"`
}

type MongoDBShardingTopology struct {
	// Shard component of mongodb.
	// More info: https://docs.mongodb.com/manual/core/sharded-cluster-shards/
	Shard MongoDBShardNode `json:"shard"`

	// Config Server (metadata) component of mongodb.
	// More info: https://docs.mongodb.com/manual/core/sharded-cluster-config-servers/
	ConfigServer MongoDBConfigNode `json:"configServer"`

	// Mongos (router) component of mongodb.
	// More info: https://docs.mongodb.com/manual/core/sharded-cluster-query-router/
	Mongos MongoDBMongosNode `json:"mongos"`
}

type MongoDBShardNode struct {
	// Shards represents number of shards for shard type of node
	// More info: https://docs.mongodb.com/manual/core/sharded-cluster-shards/
	Shards int32 `json:"shards"`

	// MongoDB sharding node configs
	MongoDBNode `json:",inline"`

	// Storage to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`
}

type MongoDBConfigNode struct {
	// MongoDB config server node configs
	MongoDBNode `json:",inline"`

	// Storage to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`
}

type MongoDBMongosNode struct {
	// MongoDB mongos node configs
	MongoDBNode `json:",inline"`

	// The deployment strategy to use to replace existing pods with new ones.
	// +optional
	Strategy apps.DeploymentStrategy `json:"strategy,omitempty" protobuf:"bytes,4,opt,name=strategy"`
}

type MongoDBNode struct {
	// Replicas represents number of replicas of this specific node.
	// If current node has replicaset enabled, then replicas is the amount of replicaset nodes.
	Replicas int32 `json:"replicas"`

	// Prefix is the name prefix of this node.
	Prefix string `json:"prefix,omitempty"`

	// Compute Resources required by the sidecar container.
	Resources core.ResourceRequirements `json:"resources,omitempty"`

	// ConfigSource is an optional field to provide custom configuration file for database (i.e mongod.cnf).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	ConfigSource *core.VolumeSource `json:"configSource,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty"`
}

type MongoDBStatus struct {
	Phase  DatabasePhase `json:"phase,omitempty"`
	Reason string        `json:"reason,omitempty"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration *types.IntHash `json:"observedGeneration,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type MongoDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of MongoDB TPR objects
	Items []MongoDB `json:"items,omitempty"`
}
