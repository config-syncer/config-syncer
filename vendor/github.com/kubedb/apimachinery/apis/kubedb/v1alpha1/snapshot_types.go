package v1alpha1

import (
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	store "kmodules.xyz/objectstore-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

const (
	ResourceCodeSnapshot     = "snap"
	ResourceKindSnapshot     = "Snapshot"
	ResourceSingularSnapshot = "snapshot"
	ResourcePluralSnapshot   = "snapshots"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Snapshot struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SnapshotSpec   `json:"spec,omitempty"`
	Status            SnapshotStatus `json:"status,omitempty"`
}

type SnapshotSpec struct {
	// Database name
	DatabaseName string `json:"databaseName"`

	// Snapshot Spec
	store.Backend `json:",inline"`

	// PodTemplate is an optional configuration for pods used to take database snapshots
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

	// -------------------------------------------------------------------------

	// Compute Resources required by the pod used to take database snapshots
	// Deprecated: Use podTemplate.spec.resources
	Resources *core.ResourceRequirements `json:"resources,omitempty"`
}

type SnapshotPhase string

const (
	// used for Snapshots that are currently running
	SnapshotPhaseRunning SnapshotPhase = "Running"
	// used for Snapshots that are Succeeded
	SnapshotPhaseSucceeded SnapshotPhase = "Succeeded"
	// used for Snapshots that are Failed
	SnapshotPhaseFailed SnapshotPhase = "Failed"
)

type SnapshotStatus struct {
	StartTime      *metav1.Time  `json:"startTime,omitempty"`
	CompletionTime *metav1.Time  `json:"completionTime,omitempty"`
	Phase          SnapshotPhase `json:"phase,omitempty"`
	Reason         string        `json:"reason,omitempty"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SnapshotList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of Snapshot CRD objects
	Items []Snapshot `json:"items,omitempty"`
}
