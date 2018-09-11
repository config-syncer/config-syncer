package v1alpha1

import (
	core "k8s.io/api/core/v1"
	store "kmodules.xyz/objectstore-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

type InitSpec struct {
	ScriptSource   *ScriptSourceSpec      `json:"scriptSource,omitempty"`
	SnapshotSource *SnapshotSourceSpec    `json:"snapshotSource,omitempty"`
	PostgresWAL    *PostgresWALSourceSpec `json:"postgresWAL,omitempty"`
}

type ScriptSourceSpec struct {
	ScriptPath        string `json:"scriptPath,omitempty"`
	core.VolumeSource `json:",inline,omitempty"`
}

type SnapshotSourceSpec struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type BackupScheduleSpec struct {
	CronExpression string `json:"cronExpression,omitempty"`

	store.Backend `json:",inline,omitempty"`

	// PodTemplate is an optional configuration for pods used for backup and recovery
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

	// -------------------------------------------------------------------------

	// Deprecated: Use podTemplate.spec.resources
	Resources *core.ResourceRequirements `json:"resources,omitempty"`
}

type DatabasePhase string

const (
	// used for Databases that are currently running
	DatabasePhaseRunning DatabasePhase = "Running"
	// used for Databases that are currently creating
	DatabasePhaseCreating DatabasePhase = "Creating"
	// used for Databases that are currently initializing
	DatabasePhaseInitializing DatabasePhase = "Initializing"
	// used for Databases that are Failed
	DatabasePhaseFailed DatabasePhase = "Failed"
)

type StorageType string

const (
	// default storage type and requires spec.storage to be configured
	StorageTypeDurable StorageType = "Durable"
	// Uses emptyDir as storage
	StorageTypeEphemeral StorageType = "Ephemeral"
)

type TerminationPolicy string

const (
	// Pauses database into a DormantDatabase
	TerminationPolicyPause TerminationPolicy = "Pause"
	// Deletes database pods, service, pvcs but leave the snapshot data intact. This will not create a DormantDatabase.
	TerminationPolicyDelete TerminationPolicy = "Delete"
	// Deletes database pods, service, pvcs and snapshot data. This will not create a DormantDatabase.
	TerminationPolicyWipeOut TerminationPolicy = "WipeOut"
)
