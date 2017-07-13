package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	yc "github.com/appscode/go/encoding/yaml"
	"github.com/ghodss/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ConfigSyncKey = "kubernetes.appscode.com/sync-config"
)

type EventForwarderSpec struct {
	SkipForwardingStorageChange bool     `json:"skip_forwarding_storage_change,omitempty"`
	SkipForwardingIngressChange bool     `json:"skip_forwarding_ingress_change,omitempty"`
	SkipForwardingWarningEvents bool     `json:"skip_forwarding_warning_events,omitempty"`
	ForwardingEventNamespaces   []string `json:"forwarding_event_namespace,omitempty"`
	NotifyVia                   string   `json:"notify_via,omitempty"`
}

type RecoverSpec struct {
	Path           string          `json:"path,omitempty"`
	TTL            metav1.Duration `json:"endpoint,omitempty"`
	HandleUpdate   bool            `json:"handle_update,omitempty"`
	NotifyOnChange bool            `json:"notify_on_change,omitempty"`
	NotifyVia      string          `json:"notify_via,omitempty"`
}

type ElasticSearchSpec struct {
	Endpoint           string `json:"endpoint,omitempty"`
	LogIndexPrefix     string `json:"log_index_prefix,omitempty"`
	LogStorageLifetime int64  `json:"log_storage_lifetime,omitempty"`
}

type InfluxDBSpec struct {
	Endpoint string          `json:"endpoint,omitempty"`
	Username string          `json:"username,omitempty"`
	Password string          `json:"password,omitempty"`
	TTL      metav1.Duration `json:"ttl,omitempty"`
}

// For periodic full cluster backup
// https://github.com/appscode/kubed/issues/16
type BackupSpec struct {
	Schedule string  `json:"schedule,omitempty"`
	Sanitize bool    `json:"sanitize,omitempty"`
	Storage  Backend `json:",inline"`
}

type ClusterConfig struct {
	ElasticSearch *ElasticSearchSpec
	InfluxDB      *InfluxDBSpec
	Backup        *BackupSpec

	Recover RecoverSpec

	EventLogger struct {
		NotifyVia string
		Namespace []string // only email for a fixed set of namespaces (Optional)
	}
}

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

func LoadConfig(configPath string) (*ClusterConfig, error) {
	if _, err := os.Stat(configPath); err != nil {
		return nil, err
	}
	os.Chmod(configPath, 0600)

	cfg := &ClusterConfig{}
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	jsonData, err := yc.ToJSON(bytes)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonData, cfg)
	return cfg, err
}

func (cfg *ClusterConfig) Save(configPath string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	os.MkdirAll(filepath.Dir(configPath), 0755)
	if err := ioutil.WriteFile(configPath, data, 0600); err != nil {
		return err
	}
	return nil
}
