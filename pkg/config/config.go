package config

import (
	stringz "github.com/appscode/go/strings"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	TimestampFormat      = "20060102T150405"
	ConfigSyncKey        = "kubed.appscode.com/sync"
	ConfigSyncNsSelector = "kubed.appscode.com/sync-namespace-selector"
	ConfigOriginKey      = "kubed.appscode.com/origin"

	JanitorElasticsearch = "Elasticsearch"
	JanitorInfluxDB      = "InfluxDB"
)

type ClusterConfig struct {
	ClusterName        string              `json:"clusterName,omitempty"`
	Snapshotter        *SnapshotSpec       `json:"snapshotter,omitempty"`
	RecycleBin         *RecycleBinSpec     `json:"recycleBin,omitempty"`
	EventForwarder     *EventForwarderSpec `json:"eventForwarder,omitempty"`
	EnableConfigSyncer bool                `json:"enableConfigSyncer"`
	NotifierSecretName string              `json:"notifierSecretName,omitempty"`
	Janitors           []JanitorSpec       `json:"janitors,omitempty"`
	APIServer          APIServerSpec       `json:"apiServer,omitempty"`
}

type APIServerSpec struct {
	Address            string `json:"address"`
	EnableSearchIndex  bool   `json:"enableSearchIndex"`
	EnableReverseIndex bool   `json:"enableReverseIndex"`
}

type JanitorSpec struct {
	Kind          string             `json:"kind"`
	TTL           metav1.Duration    `json:"ttl"`
	Elasticsearch *ElasticsearchSpec `json:"elasticsearch,omitempty,omitempty"`
	InfluxDB      *InfluxDBSpec      `json:"influxdb,omitempty"`
}

type ElasticsearchSpec struct {
	Endpoint       string `json:"endpoint,omitempty"`
	LogIndexPrefix string `json:"logIndexPrefix,omitempty"`
	SecretName     string `json:"secretName,omitempty"`
}

type InfluxDBSpec struct {
	Endpoint string `json:"endpoint,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type RecycleBinSpec struct {
	Path          string          `json:"path,omitempty"`
	TTL           metav1.Duration `json:"ttl,omitempty"`
	HandleUpdates bool            `json:"handleUpdates"`
	Receivers     []Receiver      `json:"receivers,omitempty"`
}

type EventForwarderSpec struct {
	NodeAdded     ForwarderSpec `json:"nodeAdded,omitempty"`
	StorageAdded  ForwarderSpec `json:"storageAdded,omitempty"`
	IngressAdded  ForwarderSpec `json:"ingressAdded,omitempty"`
	WarningEvents ForwarderSpec `json:"warningEvents,omitempty"`
	CSREvents     ForwarderSpec `json:"csrEvents,omitempty"`
	Receivers     []Receiver    `json:"receivers,omitempty"`
}

type NoNamespacedForwarderSpec struct {
	Namespaces []string `json:"namespaces,omitempty"`
}

type ForwarderSpec struct {
	Handle     bool     `json:"handle"`
	Namespaces []string `json:"namespaces,omitempty"`
}

func (f ForwarderSpec) IsAllowed(ns string) bool {
	return len(f.Namespaces) == 0 || stringz.Contains(f.Namespaces, ns)
}

type Receiver struct {
	// To whom notification will be sent
	To []string `json:"to,omitempty"`

	// How this notification will be sent
	Notifier string `json:"notifier,omitempty"`
}

// For periodic full cluster backup
// https://github.com/appscode/kubed/issues/16
type SnapshotSpec struct {
	Schedule  string `json:"schedule,omitempty"`
	Sanitize  bool   `json:"sanitize,omitempty"`
	Overwrite bool   `json:"overwrite,omitempty"`
	Backend   `json:",inline,omitempty"`
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

type JanitorAuthInfo struct {
	CACertData         []byte `envconfig:"CA_CERT_DATA"`
	ClientCertData     []byte `envconfig:"CLIENT_CERT_DATA"`
	ClientKeyData      []byte `envconfig:"CLIENT_KEY_DATA"`
	InsecureSkipVerify bool   `envconfig:"INSECURE_SKIP_VERIFY"`

	Username string `envconfig:"USERNAME"`
	Password string `envconfig:"PASSWORD"`
	Token    string `envconfig:"TOKEN"`
}
