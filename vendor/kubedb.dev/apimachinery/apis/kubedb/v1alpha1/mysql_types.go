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
	ResourceCodeMySQL     = "my"
	ResourceKindMySQL     = "MySQL"
	ResourceSingularMySQL = "mysql"
	ResourcePluralMySQL   = "mysqls"
)

type MySQLClusterMode string

const (
	MySQLClusterModeGroup MySQLClusterMode = "GroupReplication"
)

type MySQLGroupMode string

const (
	MySQLGroupModeSinglePrimary MySQLGroupMode = "Single-Primary"
	MySQLGroupModeMultiPrimary  MySQLGroupMode = "Multi-Primary"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Mysql defines a Mysql database.
type MySQL struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MySQLSpec   `json:"spec,omitempty"`
	Status            MySQLStatus `json:"status,omitempty"`
}

type MySQLSpec struct {
	// Version of MySQL to be deployed.
	Version types.StrYo `json:"version"`

	// Number of instances to deploy for a MySQL database. In case of MySQL group
	// replication, max allowed value is 9 (default 3).
	// (see ref: https://dev.mysql.com/doc/refman/5.7/en/group-replication-frequently-asked-questions.html)
	Replicas *int32 `json:"replicas,omitempty"`

	// MySQL cluster topology
	Topology *MySQLClusterTopology `json:"topology,omitempty"`

	// StorageType can be durable (default) or ephemeral
	StorageType StorageType `json:"storageType,omitempty"`

	// Storage spec to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// Database authentication secret
	DatabaseSecret *core.SecretVolumeSource `json:"databaseSecret,omitempty"`

	// Init is used to initialize database
	// +optional
	Init *InitSpec `json:"init,omitempty"`

	// BackupSchedule spec to specify how database backup will be taken
	// +optional
	BackupSchedule *BackupScheduleSpec `json:"backupSchedule,omitempty"`

	// Monitor is used monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`

	// ConfigSource is an optional field to provide custom configuration file for database (i.e custom-mysql.cnf).
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	ConfigSource *core.VolumeSource `json:"configSource,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

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

type MySQLClusterTopology struct {
	// If set to -
	// "GroupReplication", GroupSpec is required and MySQL servers will start  a replication group
	Mode *MySQLClusterMode `json:"mode,omitempty"`

	// Group replication info for MySQL
	Group *MySQLGroupSpec `json:"group,omitempty"`
}

type MySQLGroupSpec struct {
	// TODO: "Multi-Primary" needs to be implemented
	// Group Replication can be deployed in either "Single-Primary" or "Multi-Primary" mode
	Mode *MySQLGroupMode `json:"mode,omitempty"`

	// Group name is a version 4 UUID
	// ref: https://dev.mysql.com/doc/refman/5.7/en/group-replication-options.html#sysvar_group_replication_group_name
	Name string `json:"name,omitempty"`

	// On a replication master and each replication slave, the --server-id
	// option must be specified to establish a unique replication ID in the
	// range from 1 to 2^32 − 1. “Unique”, means that each ID must be different
	// from every other ID in use by any other replication master or slave.
	// ref: https://dev.mysql.com/doc/refman/5.7/en/server-system-variables.html#sysvar_server_id
	//
	// So, BaseServerID is needed to calculate a unique server_id for each member.
	BaseServerID *uint `json:"baseServerID,omitempty"`
}

type MySQLStatus struct {
	Phase  DatabasePhase `json:"phase,omitempty"`
	Reason string        `json:"reason,omitempty"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration *types.IntHash `json:"observedGeneration,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type MySQLList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of MySQL TPR objects
	Items []MySQL `json:"items,omitempty"`
}
