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
	ResourceCodeMemcached     = "mc"
	ResourceKindMemcached     = "Memcached"
	ResourceSingularMemcached = "memcached"
	ResourcePluralMemcached   = "memcacheds"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Memcached defines a Memcached database.
type Memcached struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MemcachedSpec   `json:"spec,omitempty"`
	Status            MemcachedStatus `json:"status,omitempty"`
}

type MemcachedSpec struct {
	// Version of Memcached to be deployed.
	Version types.StrYo `json:"version"`

	// Number of instances to deploy for a Memcached database.
	Replicas *int32 `json:"replicas,omitempty"`

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

	// The deployment strategy to use to replace existing pods with new ones.
	// +optional
	UpdateStrategy apps.DeploymentStrategy `json:"strategy,omitempty" protobuf:"bytes,4,opt,name=strategy"`

	// TerminationPolicy controls the delete operation for database
	// +optional
	TerminationPolicy TerminationPolicy `json:"terminationPolicy,omitempty"`

	// -------------------------------------------------------------------------

	// If DoNotPause is true, controller will prevent to delete this Memcached object.
	// Controller will create same Memcached object and ignore other process.
	// +optional
	// Deprecated: Use terminationPolicy = DoNotTerminate
	DoNotPause bool `json:"doNotPause,omitempty"`

	// NodeSelector is a selector which must be true for the pod to fit on a node
	// +optional
	// Deprecated: Use podTemplate.spec.nodeSelector
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Compute Resources required by the sidecar container.
	// Deprecated: Use podTemplate.spec.resources
	Resources *core.ResourceRequirements `json:"resources,omitempty"`

	// If specified, the pod's scheduling constraints
	// +optional
	// Deprecated: Use podTemplate.spec.affinity
	Affinity *core.Affinity `json:"affinity,omitempty"`

	// If specified, the pod will be dispatched by specified scheduler.
	// If not specified, the pod will be dispatched by default scheduler.
	// +optional
	// Deprecated: Use podTemplate.spec.schedulerName
	SchedulerName string `json:"schedulerName,omitempty"`

	// If specified, the pod's tolerations.
	// +optional
	// Deprecated: Use podTemplate.spec.tolerations
	Tolerations []core.Toleration `json:"tolerations,omitempty"`

	// ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec.
	// If specified, these secrets will be passed to individual puller implementations for them to use.
	// +optional
	// Deprecated: Use podTemplate.spec.imagePullSecrets
	ImagePullSecrets []core.LocalObjectReference `json:"imagePullSecrets,omitempty"`
}

type MemcachedStatus struct {
	Phase  DatabasePhase `json:"phase,omitempty"`
	Reason string        `json:"reason,omitempty"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration *types.IntHash `json:"observedGeneration,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type MemcachedList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of Memcached TPR objects
	Items []Memcached `json:"items,omitempty"`
}
