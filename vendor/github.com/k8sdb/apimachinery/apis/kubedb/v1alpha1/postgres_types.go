package v1alpha1

import (
	"github.com/appscode/go/encoding/json/types"
	"github.com/appscode/kube-mon/api"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceCodePostgres = "pg"
	ResourceKindPostgres = "Postgres"
	ResourceNamePostgres = "postgres"
	ResourceTypePostgres = "postgreses"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Postgres defines a Postgres database.
type Postgres struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PostgresSpec   `json:"spec,omitempty"`
	Status            PostgresStatus `json:"status,omitempty"`
}

type PostgresSpec struct {
	// Version of Postgres to be deployed.
	Version types.StrYo `json:"version,omitempty"`
	// Number of instances to deploy for a Postgres database.
	Replicas int32 `json:"replicas,omitempty"`
	// Standby mode
	Standby StandbyMode `json:"standby,omitempty"`
	// Streaming mode
	Streaming StreamingMode `json:"streaming,omitempty"`
	// Archive for wal files
	Archiver *PostgresArchiverSpec `json:"archiver,omitempty"`
	// Database authentication secret
	DatabaseSecret *core.SecretVolumeSource `json:"databaseSecret,omitempty"`
	// Storage to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`
	// NodeSelector is a selector which must be true for the pod to fit on a node
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// Init is used to initialize database
	// +optional
	Init *InitSpec `json:"init,omitempty"`
	// BackupSchedule spec to specify how database backup will be taken
	// +optional
	BackupSchedule *BackupScheduleSpec `json:"backupSchedule,omitempty"`
	// If DoNotPause is true, controller will prevent to delete this Postgres object.
	// Controller will create same Postgres object and ignore other process.
	// +optional
	DoNotPause bool `json:"doNotPause,omitempty"`
	// Monitor is used monitor database instance
	// +optional
	Monitor *api.AgentSpec `json:"monitor,omitempty"`
	// Compute Resources required by the sidecar container.
	Resources core.ResourceRequirements `json:"resources,omitempty"`
	// If specified, the pod's scheduling constraints
	// +optional
	Affinity *core.Affinity `json:"affinity,omitempty" protobuf:"bytes,18,opt,name=affinity"`
	// If specified, the pod will be dispatched by specified scheduler.
	// If not specified, the pod will be dispatched by default scheduler.
	// +optional
	SchedulerName string `json:"schedulerName,omitempty" protobuf:"bytes,19,opt,name=schedulerName"`
	// If specified, the pod's tolerations.
	// +optional
	Tolerations []core.Toleration `json:"tolerations,omitempty" protobuf:"bytes,22,opt,name=tolerations"`
}

type PostgresArchiverSpec struct {
	Storage *SnapshotStorageSpec `json:"storage,omitempty"`
	// wal_keep_segments
}

type PostgresStatus struct {
	CreationTime *metav1.Time  `json:"creationTime,omitempty"`
	Phase        DatabasePhase `json:"phase,omitempty"`
	Reason       string        `json:"reason,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PostgresList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of Postgres CRD objects
	Items []Postgres `json:"items,omitempty"`
}

// Following structures are used for audit summary report
type PostgresTableInfo struct {
	TotalRow int64 `json:"totalRow"`
	MaxID    int64 `json:"maxId"`
	NextID   int64 `json:"nextId"`
}

type PostgresSchemaInfo struct {
	Table map[string]*PostgresTableInfo `json:"table"`
}

type PostgresSummary struct {
	Schema map[string]*PostgresSchemaInfo `json:"schema"`
}

type PostgresWALSourceSpec struct {
	BackupName          string `json:"backupName,omitempty"`
	PITR                string `json:"pitr,omitempty"`
	SnapshotStorageSpec `json:",inline,omitempty"`
}

type StandbyMode string

const (
	HotStandby  StandbyMode = "hot"
	WarmStandby StandbyMode = "warm"
)

type StreamingMode string

const (
	SynchronousStreaming  StreamingMode = "synchronous"
	AsynchronousStreaming StreamingMode = "asynchronous"
)
