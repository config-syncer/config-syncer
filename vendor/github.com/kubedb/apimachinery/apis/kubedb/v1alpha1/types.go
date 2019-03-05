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
	// Arguments to the restore job
	Args []string `json:"args,omitempty"`
}

type BackupScheduleSpec struct {
	CronExpression string `json:"cronExpression,omitempty"`

	store.Backend `json:",inline,omitempty"`

	// PodTemplate is an optional configuration for pods used for backup and recovery
	// +optional
	PodTemplate ofst.PodTemplateSpec `json:"podTemplate,omitempty"`

	// PodVolumeClaimSpec is used to specify temporary storage for backup/restore Job
	// +optional
	PodVolumeClaimSpec *core.PersistentVolumeClaimSpec `json:"podVolumeClaimSpec,omitempty"`

	// -------------------------------------------------------------------------

	// Deprecated: Use podTemplate.spec.resources
	Resources *core.ResourceRequirements `json:"resources,omitempty"`
}

// LeaderElectionConfig contains essential attributes of leader election.
// ref: https://github.com/kubernetes/client-go/blob/6134db91200ea474868bc6775e62cc294a74c6c6/tools/leaderelection/leaderelection.go#L105-L114
type LeaderElectionConfig struct {
	// LeaseDuration is the duration in second that non-leader candidates will
	// wait to force acquire leadership. This is measured against time of
	// last observed ack. Default 15
	LeaseDurationSeconds int32 `json:"leaseDurationSeconds"`
	// RenewDeadline is the duration in second that the acting master will retry
	// refreshing leadership before giving up. Normally, LeaseDuration * 2 / 3.
	// Default 10
	RenewDeadlineSeconds int32 `json:"renewDeadlineSeconds"`
	// RetryPeriod is the duration in second the LeaderElector clients should wait
	// between tries of actions. Normally, LeaseDuration / 3.
	// Default 2
	RetryPeriodSeconds int32 `json:"retryPeriodSeconds"`
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
	// Rejects attempt to delete database using ValidationWebhook. This replaces spec.doNotPause = true
	TerminationPolicyDoNotTerminate TerminationPolicy = "DoNotTerminate"
)
