package config

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ConfigSyncKey = "kubernetes.appscode.com/sync"
)

type ClusterConfig struct {
	Elasticsearch      *ElasticSearchSpec  `json:"elasticsearch,omitempty,omitempty"`
	InfluxDB           *InfluxDBSpec       `json:"influxdb,omitempty"`
	RecycleBin         *RecycleBinSpec     `json:"recycleBin,omitempty"`
	EventForwarder     *EventForwarderSpec `json:"eventForwarder,omitempty"`
	ClusterSnapshot    *SnapshotSpec       `json:"clusterSnapshot,omitempty"`
	NotifierSecretName string              `json:"notifierSecretName,omitempty"`
}

type ElasticSearchSpec struct {
	Endpoint       string          `json:"endpoint,omitempty"`
	LogIndexPrefix string          `json:"logIndexPrefix,omitempty"`
	TTL            metav1.Duration `json:"ttl,omitempty"`
}

type InfluxDBSpec struct {
	Endpoint string          `json:"endpoint,omitempty"`
	Username string          `json:"username,omitempty"`
	Password string          `json:"password,omitempty"`
	TTL      metav1.Duration `json:"ttl,omitempty"`
}

type RecycleBinSpec struct {
	Path         string          `json:"path,omitempty"`
	TTL          metav1.Duration `json:"ttl,omitempty"`
	HandleUpdate bool            `json:"handle_update,omitempty"`
	NotifyVia    string          `json:"notifyVia,omitempty"`
}

type EventForwarderSpec struct {
	NotifyOnStorageAdd   bool     `json:"notifyOnStorageAdd,omitempty"`
	NotifyOnIngressAdd   bool     `json:"notifyOnIngressAdd,omitempty"`
	ForwardWarningEvents bool     `json:"forwardWarningEvents,omitempty"`
	EventNamespaces      []string `json:"eventNamespaces,omitempty"`
	NotifyVia            string   `json:"notifyVia,omitempty"`
}

// For periodic full cluster backup
// https://github.com/appscode/kubed/issues/16
type SnapshotSpec struct {
	Schedule string  `json:"schedule,omitempty"`
	Sanitize bool    `json:"sanitize,omitempty"`
	Storage  Backend `json:",inline"`
}

const (
	AWS_ACCESS_KEY_ID     = "AWS_ACCESS_KEY_ID"
	AWS_SECRET_ACCESS_KEY = "AWS_SECRET_ACCESS_KEY"

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

type Backend struct {
	StorageSecretName string `json:"storageSecretName,omitempty"`

	Local *LocalSpec `json:"local,omitempty"`
	S3    *S3Spec    `json:"s3,omitempty"`
	GCS   *GCSSpec   `json:"gcs,omitempty"`
	Azure *AzureSpec `json:"azure,omitempty"`
	Swift *SwiftSpec `json:"swift,omitempty"`
}

type LocalSpec struct {
	Path string `json:"path,omitempty"`
}

type S3Spec struct {
	Endpoint string `json:"endpoint,omitempty"`
	Bucket   string `json:"bucket,omiempty"`
	Prefix   string `json:"prefix,omitempty"`
}

type GCSSpec struct {
	Bucket string `json:"bucket,omiempty"`
	Prefix string `json:"prefix,omitempty"`
}

type AzureSpec struct {
	Container string `json:"container,omitempty"`
	Prefix    string `json:"prefix,omitempty"`
}

type SwiftSpec struct {
	Container string `json:"container,omitempty"`
	Prefix    string `json:"prefix,omitempty"`
}
