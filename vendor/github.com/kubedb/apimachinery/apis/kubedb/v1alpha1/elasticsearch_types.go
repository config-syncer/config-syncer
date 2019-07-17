package v1alpha1

import (
	"github.com/appscode/go/encoding/json/types"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

const (
	ResourceCodeElasticsearch     = "es"
	ResourceKindElasticsearch     = "Elasticsearch"
	ResourceSingularElasticsearch = "elasticsearch"
	ResourcePluralElasticsearch   = "elasticsearches"
)

// +genclient
// +k8s:openapi-gen=true
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
	Version types.StrYo `json:"version"`

	// Number of instances to deploy for a Elasticsearch database.
	Replicas *int32 `json:"replicas,omitempty"`

	// Elasticsearch topology for node specification
	Topology *ElasticsearchClusterTopology `json:"topology,omitempty"`

	// To enable ssl in transport & http layer
	EnableSSL bool `json:"enableSSL,omitempty"`

	// Secret with SSL certificates
	CertificateSecret *core.SecretVolumeSource `json:"certificateSecret,omitempty"`

	// Authentication plugin used by Elasticsearch cluster. If unset, defaults to SearchGuard.
	AuthPlugin ElasticsearchAuthPlugin `json:"authPlugin,omitempty"`

	// Database authentication secret
	DatabaseSecret *core.SecretVolumeSource `json:"databaseSecret,omitempty"`

	// StorageType can be durable (default) or ephemeral
	StorageType StorageType `json:"storageType,omitempty"`

	// Storage to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// Init is used to initialize database
	// +optional
	Init *InitSpec `json:"init,omitempty"`

	// BackupSchedule spec to specify how database backup will be taken
	// +optional
	BackupSchedule *BackupScheduleSpec `json:"backupSchedule,omitempty"`

	// Monitor is used monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`

	// ConfigSource is an optional field to provide custom configuration file for database.
	// If specified, this file will be used as configuration file otherwise default configuration file will be used.
	ConfigSource *core.VolumeSource `json:"configSource,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

	// ServiceTemplate is an optional configuration for service used to expose database
	// +optional
	ServiceTemplate ofst.ServiceTemplateSpec `json:"serviceTemplate,omitempty"`

	// An eviction is allowed if at most "maxUnavailable" pods selected by
	// "selector" are unavailable after the eviction, i.e. even in absence of
	// the evicted pod. For example, one can prevent all voluntary evictions
	// by specifying 0. This is a mutually exclusive setting with "minAvailable".
	// +optional
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty"`

	// updateStrategy indicates the StatefulSetUpdateStrategy that will be
	// employed to update Pods in the StatefulSet when a revision is made to
	// Template.
	UpdateStrategy apps.StatefulSetUpdateStrategy `json:"updateStrategy,omitempty"`

	// TerminationPolicy controls the delete operation for database
	// +optional
	TerminationPolicy TerminationPolicy `json:"terminationPolicy,omitempty"`
}

type ElasticsearchClusterTopology struct {
	Master ElasticsearchNode `json:"master"`
	Data   ElasticsearchNode `json:"data"`
	Client ElasticsearchNode `json:"client"`
}

type ElasticsearchNode struct {
	// Replicas represents number of replica for this specific type of node
	Replicas *int32 `json:"replicas,omitempty"`
	Prefix   string `json:"prefix,omitempty"`
	// Storage to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`
	// Compute Resources required by the sidecar container.
	Resources core.ResourceRequirements `json:"resources,omitempty"`
	// An eviction is allowed if at most "maxUnavailable" pods selected by
	// "selector" are unavailable after the eviction, i.e. even in absence of
	// the evicted pod. For example, one can prevent all voluntary evictions
	// by specifying 0. This is a mutually exclusive setting with "minAvailable".
	// +optional
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty"`
}

type ElasticsearchStatus struct {
	Phase  DatabasePhase `json:"phase,omitempty"`
	Reason string        `json:"reason,omitempty"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration *types.IntHash `json:"observedGeneration,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ElasticsearchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of Elasticsearch CRD objects
	Items []Elasticsearch `json:"items,omitempty"`
}

// +k8s:deepcopy-gen=false
// +k8s:gen-deepcopy=false
// Following structure is used for audit summary report
type ElasticsearchSummary struct {
	IdCount map[string]int64 `json:"idCount"`
	Mapping interface{}      `json:"mapping"`
	Setting interface{}      `json:"setting"`
}

type ElasticsearchAuthPlugin string

const (
	ElasticsearchAuthPluginSearchGuard ElasticsearchAuthPlugin = "SearchGuard" // Default
	ElasticsearchAuthPluginNone        ElasticsearchAuthPlugin = "None"
	ElasticsearchAuthPluginXpack       ElasticsearchAuthPlugin = "X-Pack"
)
