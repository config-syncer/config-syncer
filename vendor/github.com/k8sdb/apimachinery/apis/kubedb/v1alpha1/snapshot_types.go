package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

const (
	ResourceCodeSnapshot = "snap"
	ResourceKindSnapshot = "Snapshot"
	ResourceNameSnapshot = "snapshot"
	ResourceTypeSnapshot = "snapshots"
)

// +genclient=true
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
	DatabaseName string `json:"databaseName,omitempty"`
	// Snapshot Spec
	SnapshotStorageSpec `json:",inline,omitempty"`
	// Compute Resources required by the sidecar container.
	Resources apiv1.ResourceRequirements `json:"resources,omitempty"`
}

type SnapshotPhase string

const (
	// used for Snapshots that are currently running
	SnapshotPhaseRunning SnapshotPhase = "Running"
	// used for Snapshots that are Succeeded
	SnapshotPhaseSuccessed SnapshotPhase = "Succeeded"
	// used for Snapshots that are Failed
	SnapshotPhaseFailed SnapshotPhase = "Failed"
)

type SnapshotStatus struct {
	StartTime      *metav1.Time  `json:"startTime,omitempty"`
	CompletionTime *metav1.Time  `json:"completionTime,omitempty"`
	Phase          SnapshotPhase `json:"phase,omitempty"`
	Reason         string        `json:"reason,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SnapshotList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of Snapshot CRD objects
	Items []Snapshot `json:"items,omitempty"`
}
