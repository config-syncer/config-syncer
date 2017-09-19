package kubedb

import (
	"github.com/appscode/go/encoding/json/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

const (
	ResourceCodeElasticsearch = "es"
	ResourceKindElasticsearch = "Elasticsearch"
	ResourceNameElasticsearch = "elasticsearch"
	ResourceTypeElasticsearch = "elasticsearchs"
)

// +genclient=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Elasticsearch defines a Elasticsearch database.
type Elasticsearch struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ElasticsearchSpec   `json:"spec,omitempty"`
	Status            ElasticsearchStatus `json:"status,omitempty"`
}

type ElasticsearchSpec struct {
	// Version of Elasticsearch to be deployed.
	Version types.StrYo `json:"version,omitempty"`
	// Number of instances to deploy for a Elasticsearch database.
	Replicas int32 `json:"replicas,omitempty"`
	// Storage to specify how storage shall be used.
	Storage *apiv1.PersistentVolumeClaimSpec `json:"storage,omitempty"`
	// NodeSelector is a selector which must be true for the pod to fit on a node
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// Init is used to initialize database
	// +optional
	Init *InitSpec `json:"init,omitempty"`
	// BackupSchedule spec to specify how database backup will be taken
	// +optional
	BackupSchedule *BackupScheduleSpec `json:"backupSchedule,omitempty"`
	// If DoNotPause is true, controller will prevent to delete this Elasticsearch object.
	// Controller will create same Elasticsearch object and ignore other process.
	// +optional
	DoNotPause bool `json:"doNotPause,omitempty"`
	// Monitor is used monitor database instance
	// +optional
	Monitor *MonitorSpec `json:"monitor,omitempty"`
	// Compute Resources required by the sidecar container.
	Resources apiv1.ResourceRequirements `json:"resources,omitempty"`
}

type ElasticsearchStatus struct {
	CreationTime *metav1.Time  `json:"creationTime,omitempty"`
	Phase        DatabasePhase `json:"phase,omitempty"`
	Reason       string        `json:"reason,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ElasticsearchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of Elasticsearch CRD objects
	Items []Elasticsearch `json:"items,omitempty"`
}

// Following structure is used for audit summary report
type ElasticsearchSummary struct {
	IdCount map[string]int64 `json:"idCount"`
	Mapping interface{}      `json:"mapping"`
	Setting interface{}      `json:"setting"`
}
