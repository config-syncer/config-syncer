package config

import (
	"encoding/json"
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

func (b Backend) Location(timestamp time.Time) (string, error) {
	ts := timestamp.UTC().Format(time.RFC3339)
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
