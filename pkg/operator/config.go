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

package operator

import (
	"path/filepath"
	"time"

	"github.com/appscode/kubed/pkg/eventer"
	"github.com/appscode/kubed/pkg/syncer"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"kmodules.xyz/client-go/discovery"
	"kmodules.xyz/client-go/tools/fsnotify"
)

type Config struct {
	ScratchDir        string
	ConfigPath        string
	OperatorNamespace string
	ResyncPeriod      time.Duration
	Test              bool
}

type OperatorConfig struct {
	Config

	ClientConfig *rest.Config
	KubeClient   kubernetes.Interface
}

func NewOperatorConfig(clientConfig *rest.Config) *OperatorConfig {
	return &OperatorConfig{
		ClientConfig: clientConfig,
	}
}

func (c *OperatorConfig) New() (*Operator, error) {
	if err := discovery.IsDefaultSupportedVersion(c.KubeClient); err != nil {
		return nil, err
	}

	op := &Operator{
		Config:       c.Config,
		ClientConfig: c.ClientConfig,
		KubeClient:   c.KubeClient,
	}

	op.recorder = eventer.NewEventRecorder(op.KubeClient, "kubed")
	op.configSyncer = syncer.New(op.KubeClient, op.recorder)

	op.Configure()

	op.watcher = &fsnotify.Watcher{
		WatchDir: filepath.Dir(c.ConfigPath),
		Reload:   op.Configure,
	}

	// ---------------------------
	op.kubeInformerFactory = informers.NewSharedInformerFactory(op.KubeClient, c.ResyncPeriod)
	// ---------------------------
	op.setupConfigInformers()
	// ---------------------------

	if err := op.Configure(); err != nil {
		return nil, err
	}
	return op, nil
}
