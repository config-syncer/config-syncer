package v1alpha1

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	yc "github.com/appscode/go/encoding/yaml"
	"github.com/ghodss/yaml"
	"github.com/json-iterator/go"
	"github.com/pkg/errors"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func LoadConfig(configPath string) (*ClusterConfig, error) {
	if _, err := os.Stat(configPath); err != nil {
		return nil, errors.Errorf("failed to find file %s. Reason: %s", configPath, err)
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
	if cfg.EventForwarder != nil {
		if cfg.EventForwarder.NodeAdded != nil {
			return errors.Errorf("convert `nodeAdded` spec to eventForwarder.rules format")
		}
		if cfg.EventForwarder.StorageAdded != nil {
			return errors.Errorf("convert `storageAdded` spec to eventForwarder.rules format")
		}
		if cfg.EventForwarder.IngressAdded != nil {
			return errors.Errorf("convert `ingressAdded` spec to eventForwarder.rules format")
		}
		if cfg.EventForwarder.WarningEvents != nil {
			return errors.Errorf("convert `warningEvents` spec to eventForwarder.rules format")
		}
		if cfg.EventForwarder.CSREvents != nil {
			return errors.Errorf("convert `csrEvents` spec to eventForwarder.rules format")
		}

		for i, rule := range cfg.EventForwarder.Rules {
			for j, op := range rule.Operations {
				if op != Create && op != Delete {
					return errors.Errorf("unknown operation %s found at eventForwarder.rules[%d].operations[%d]", op, i, j)
				}
			}
		}
	}

	for _, j := range cfg.Janitors {
		switch j.Kind {
		case JanitorElasticsearch:
			if j.Elasticsearch == nil {
				return errors.Errorf("missing spec for janitor kind %s", j.Kind)
			}
		case JanitorInfluxDB:
			if j.InfluxDB == nil {
				return errors.Errorf("missing spec for janitor kind %s", j.Kind)
			}
		default:
			return errors.Errorf("unknown janitor kind %s", j.Kind)
		}
	}
	return nil
}

func (b SnapshotSpec) Location(filename string) (string, error) {
	ts := filename
	if b.Overwrite {
		ts = "snapshot.tar.gz"
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

func (b Backend) GetBucketAndPrefix() (string, string, error) {
	if b.S3 != nil {
		return b.S3.Bucket, b.S3.Prefix, nil
	} else if b.GCS != nil {
		return b.GCS.Bucket, b.GCS.Prefix, nil
	} else if b.Azure != nil {
		return b.Azure.Container, b.Azure.Prefix, nil
	} else if b.Swift != nil {
		return b.Swift.Container, b.Swift.Prefix, nil
	} else if b.Local !=nil{
		return b.Local.Path,"",nil
	}
	return "", "", errors.New("unknown backend type.")
}

func LoadJanitorAuthInfo(data map[string][]byte) *JanitorAuthInfo {
	if data == nil {
		return &JanitorAuthInfo{}
	}
	insecureSkipVerify, _ := strconv.ParseBool(string(data["INSECURE_SKIP_VERIFY"]))

	return &JanitorAuthInfo{
		CACertData:         data["CA_CERT_DATA"],
		ClientCertData:     data["CLIENT_CERT_DATA"],
		ClientKeyData:      data["CLIENT_KEY_DATA"],
		InsecureSkipVerify: insecureSkipVerify,
		Username:           string(data["USERNAME"]),
		Password:           string(data["PASSWORD"]),
		Token:              string(data["TOKEN"]),
	}
}
