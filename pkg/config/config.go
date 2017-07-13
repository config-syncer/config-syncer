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
	ConfigSyncKey = "kubernetes.appscode.com/sync"
)

type ClusterConfig struct {
	Elasticsearch      *ElasticSearchSpec  `json:"elasticsearch,omitempty,omitempty"`
	InfluxDB           *InfluxDBSpec       `json:"influxdb,omitempty"`
	RecycleBin         *RecycleBinSpec     `json:"recycle_bin,omitempty"`
	EventForwarder     *EventForwarderSpec `json:"event_forwarder,omitempty"`
	Backup             *BackupSpec         `json:"backup,omitempty"`
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
type BackupSpec struct {
	Schedule string  `json:"schedule,omitempty"`
	Sanitize bool    `json:"sanitize,omitempty"`
	Storage  Backend `json:",inline"`
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
