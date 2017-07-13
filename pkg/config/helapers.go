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

func (s Backend) Location(timestamp time.Time) (string, error) {
	spec := s
	if spec.S3 != nil {
		return filepath.Join(spec.S3.Prefix, "kubed", timestamp.UTC().Format(time.RFC3339)), nil
	} else if spec.GCS != nil {
		return filepath.Join(spec.GCS.Prefix, "kubed", timestamp.UTC().Format(time.RFC3339)), nil
	} else if spec.Azure != nil {
		return filepath.Join(spec.Azure.Prefix, "kubed", timestamp.UTC().Format(time.RFC3339)), nil
	} else if spec.Local != nil {
		return filepath.Join("kubed", timestamp.UTC().Format(time.RFC3339)), nil
	} else if spec.Swift != nil {
		return filepath.Join(spec.Swift.Prefix, "kubed", timestamp.UTC().Format(time.RFC3339)), nil
	}
	return "", errors.New("No storage provider is configured.")
}

func (s Backend) Container() (string, error) {
	if s.S3 != nil {
		return s.S3.Bucket, nil
	} else if s.GCS != nil {
		return s.GCS.Bucket, nil
	} else if s.Azure != nil {
		return s.Azure.Container, nil
	} else if s.Local != nil {
		return s.Local.Path, nil
	} else if s.Swift != nil {
		return s.Swift.Container, nil
	}
	return "", errors.New("No storage provider is configured.")
}
