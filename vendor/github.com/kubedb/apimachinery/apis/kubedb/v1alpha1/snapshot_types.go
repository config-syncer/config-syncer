package v1alpha1

import (
	"github.com/appscode/go/encoding/json/types"
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

	// StorageType can be durable or ephemeral.
	// If not given, database storage type will be used.
	// +optional
	StorageType *StorageType `json:"storageType,omitempty"`

	// PodTemplate is an optional configuration for pods used to take database snapshots
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

	// PodVolumeClaimSpec is used to specify temporary storage for backup/restore Job.
	// If not given, database's PvcSpec will be used.
	// If storageType is durable, then a PVC will be created using this PVCSpec.
	// If storageType is ephemeral, then an empty directory will be created of size PvcSpec.Resources.Requests[core.ResourceStorage].
	// +optional
	PodVolumeClaimSpec *core.PersistentVolumeClaimSpec `json:"podVolumeClaimSpec,omitempty"`
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
	ObservedGeneration *types.IntHash `json:"observedGeneration,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SnapshotList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of Snapshot CRD objects
	Items []Snapshot `json:"items,omitempty"`
}
