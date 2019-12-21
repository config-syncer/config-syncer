/*
Copyright The Kubed Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"io/ioutil"
	"os"
	"path/filepath"

	yc "github.com/appscode/go/encoding/yaml"

	"github.com/ghodss/yaml"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func LoadConfig(configPath string) (*ClusterConfig, error) {
	if _, err := os.Stat(configPath); err != nil {
		return nil, errors.Errorf("failed to find file %s. Reason: %s", configPath, err)
	}
	if err := os.Chmod(configPath, 0600); err != nil {
		return nil, err
	}

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
	err = os.MkdirAll(filepath.Dir(configPath), 0755)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(configPath, data, 0600); err != nil {
		return err
	}
	return nil
}
