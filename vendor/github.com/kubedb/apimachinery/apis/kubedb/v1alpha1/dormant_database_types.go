package v1alpha1

import (
	"github.com/appscode/go/encoding/json/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceCodeDormantDatabase     = "drmn"
	ResourceKindDormantDatabase     = "DormantDatabase"
	ResourceSingularDormantDatabase = "dormantdatabase"
	ResourcePluralDormantDatabase   = "dormantdatabases"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type DormantDatabase struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DormantDatabaseSpec   `json:"spec,omitempty"`
	Status            DormantDatabaseStatus `json:"status,omitempty"`
}

type DormantDatabaseSpec struct {
	// If true, invoke wipe out operation
	// +optional
	WipeOut bool `json:"wipeOut,omitempty"`
	// Origin to store original database information
	Origin Origin `json:"origin"`
}

type Origin struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Origin Spec to store original database Spec
	Spec OriginSpec `json:"spec"`
}

type OriginSpec struct {
	// Elasticsearch Spec
	// +optional
	Elasticsearch *ElasticsearchSpec `json:"elasticsearch,omitempty"`
	// Postgres Spec
	// +optional
	Postgres *PostgresSpec `json:"postgres,omitempty"`
	// MySQL Spec
	// +optional
	MySQL *MySQLSpec `json:"mysql,omitempty"`
	// MongoDB Spec
	// +optional
	MongoDB *MongoDBSpec `json:"mongodb,omitempty"`
	// Redis Spec
	// +optional
	Redis *RedisSpec `json:"redis,omitempty"`
	// Memcached Spec
	// +optional
	Memcached *MemcachedSpec `json:"memcached,omitempty"`
	// Etcd Spec
	// +optional
	Etcd *EtcdSpec `json:"etcd,omitempty"`
}

type DormantDatabasePhase string

const (
	// used for Databases that are paused
	DormantDatabasePhasePaused DormantDatabasePhase = "Paused"
	// used for Databases that are currently pausing
	DormantDatabasePhasePausing DormantDatabasePhase = "Pausing"
	// used for Databases that are wiped out
	DormantDatabasePhaseWipedOut DormantDatabasePhase = "WipedOut"
	// used for Databases that are currently wiping out
	DormantDatabasePhaseWipingOut DormantDatabasePhase = "WipingOut"
	// used for Databases that are currently recovering
	DormantDatabasePhaseResuming DormantDatabasePhase = "Resuming"
)

type DormantDatabaseStatus struct {
	PausingTime *metav1.Time         `json:"pausingTime,omitempty"`
	WipeOutTime *metav1.Time         `json:"wipeOutTime,omitempty"`
	Phase       DormantDatabasePhase `json:"phase,omitempty"`
	Reason      string               `json:"reason,omitempty"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration *types.IntHash `json:"observedGeneration,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type DormantDatabaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of DormantDatabase CRD objects
	Items []DormantDatabase `json:"items,omitempty"`
}
