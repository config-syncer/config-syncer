package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	yc "github.com/appscode/go/encoding/yaml"
	"github.com/ghodss/yaml"
)

type RecoverSpec struct {
	Path              string
	TTL               time.Duration
	HandleSpecUpdates bool
	EmailOneDelete    bool // Notify Via
}

type ClusterConfig struct {
	LogIndexPrefix            string `json:"log_index_prefix"`
	LogStorageLifetime        int64  `json:"log_storage_lifetime"`
	MonitoringStorageLifetime int64  `json:"monitoring_storage_lifetime"`

	// For periodic full cluster backup
	// https://github.com/appscode/kubed/issues/16
	Backup struct {
		Schedule string
		Storage  Backend
	}

	Recover RecoverSpec

	// Email Warning events
	EventLogger struct {
		NotifyVia string
		Namespace []string // only email for a fixed set of namespaces (Optional)
	}

	// Take ConfigMap/Secret with label to other namespaces
	// kubernetes.appscode.com/sync-config: true

	// Search
	// Reverse Index

	ESEndpoint                        string
	InfluxSecretName                  string
	InfluxSecretNamespace             string
	ClusterKubedConfigSecretName      string
	ClusterKubedConfigSecretNamespace string
	NotifyOnCertSoonToBeExpired       bool
	NotifyVia                         string
}

type Backend struct {
	StorageSecretName string `json:"storageSecretName,omitempty"`

	Local *LocalSpec `json:"local"`
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
