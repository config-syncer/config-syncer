package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/appscode/errors"
	yc "github.com/appscode/go/encoding/yaml"
	"github.com/ghodss/yaml"
)

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

func (cfg ClusterConfig) Save(configPath string) error {
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

func (cfg ClusterConfig) Validate() error {
	if cfg.EventForwarder != nil && len(cfg.EventForwarder.NodeAdded.Namespaces) > 0 {
		return fmt.Errorf("Namespeces can't be defined for forwarding `nodeAdded` events.")
	}

	for _, j := range cfg.Janitors {
		switch j.Kind {
		case JanitorElasticsearch:
			if j.Elasticsearch == nil {
				return fmt.Errorf("Missing spec for janitor kind %s", j.Kind)
			}
		case JanitorInfluxDB:
			if j.InfluxDB == nil {
				return fmt.Errorf("Missing spec for janitor kind %s", j.Kind)
			}
		default:
			return fmt.Errorf("Unknown janitor kind %s", j.Kind)
		}
	}
	return nil
}

func (b SnapshotSpec) Location(timestamp time.Time) (string, error) {
	ts := "snapshot.tar.gz"
	if b.Overwrite {
		ts = timestamp.UTC().Format(TimestampFormat) + ".tar.gz"
	}
	if b.S3 != nil {
		return filepath.Join(b.S3.Prefix, ts), nil
	} else if b.GCS != nil {
		return filepath.Join(b.GCS.Prefix, ts), nil
	} else if b.Azure != nil {
		return filepath.Join(b.Azure.Prefix, ts), nil
	} else if b.Local != nil {
		return ts, nil
	} else if b.Swift != nil {
		return filepath.Join(b.Swift.Prefix, ts), nil
	}
	return "", errors.New("No storage provider is configured.")
}

func (b Backend) Container() (string, error) {
	if b.S3 != nil {
		return b.S3.Bucket, nil
	} else if b.GCS != nil {
		return b.GCS.Bucket, nil
	} else if b.Azure != nil {
		return b.Azure.Container, nil
	} else if b.Local != nil {
		return b.Local.Path, nil
	} else if b.Swift != nil {
		return b.Swift.Container, nil
	}
	return "", errors.New("No storage provider is configured.")
}
