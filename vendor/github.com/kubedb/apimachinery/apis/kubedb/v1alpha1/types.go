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

const (
	AWS_ACCESS_KEY_ID     = "AWS_ACCESS_KEY_ID"
	AWS_SECRET_ACCESS_KEY = "AWS_SECRET_ACCESS_KEY"
	CA_CERT_DATA          = "CA_CERT_DATA"

	GOOGLE_PROJECT_ID               = "GOOGLE_PROJECT_ID"
	GOOGLE_SERVICE_ACCOUNT_JSON_KEY = "GOOGLE_SERVICE_ACCOUNT_JSON_KEY"

	AZURE_ACCOUNT_NAME = "AZURE_ACCOUNT_NAME"
	AZURE_ACCOUNT_KEY  = "AZURE_ACCOUNT_KEY"

	// swift
	OS_USERNAME    = "OS_USERNAME"
	OS_PASSWORD    = "OS_PASSWORD"
	OS_REGION_NAME = "OS_REGION_NAME"
	OS_AUTH_URL    = "OS_AUTH_URL"

	// v3 specific
	OS_USER_DOMAIN_NAME    = "OS_USER_DOMAIN_NAME"
	OS_PROJECT_NAME        = "OS_PROJECT_NAME"
	OS_PROJECT_DOMAIN_NAME = "OS_PROJECT_DOMAIN_NAME"

	// v2 specific
	OS_TENANT_ID   = "OS_TENANT_ID"
	OS_TENANT_NAME = "OS_TENANT_NAME"

	// v1 specific
	ST_AUTH = "ST_AUTH"
	ST_USER = "ST_USER"
	ST_KEY  = "ST_KEY"

	// Manual authentication
	OS_STORAGE_URL = "OS_STORAGE_URL"
	OS_AUTH_TOKEN  = "OS_AUTH_TOKEN"
)

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
