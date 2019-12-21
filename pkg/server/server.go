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

package server

import (
	"github.com/appscode/kubed/apis/kubed/install"
	"github.com/appscode/kubed/pkg/operator"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/version"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

var (
	Scheme = runtime.NewScheme()
	Codecs = serializer.NewCodecFactory(Scheme)
)

func init() {
	install.Install(Scheme)

	// we need to add the options to empty v1
	// TODO fix the server code to avoid this
	metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})

	// TODO: keep the generic API server from wanting this
	unversioned := schema.GroupVersion{Group: "", Version: "v1"}
	Scheme.AddUnversionedTypes(unversioned,
		&metav1.Status{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
	)
}

type KubedConfig struct {
	GenericConfig  *genericapiserver.RecommendedConfig
	OperatorConfig *operator.OperatorConfig
}

// KubedServer contains state for a Kubernetes cluster daemon.
type KubedServer struct {
	GenericAPIServer *genericapiserver.GenericAPIServer
	Operator         *operator.Operator
}

func (op *KubedServer) Run(stopCh <-chan struct{}) error {
	go op.Operator.Run(stopCh)
	return op.GenericAPIServer.PrepareRun().Run(stopCh)
}

type completedConfig struct {
	GenericConfig  genericapiserver.CompletedConfig
	OperatorConfig *operator.OperatorConfig
}

type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (c *KubedConfig) Complete() CompletedConfig {
	completedCfg := completedConfig{
		c.GenericConfig.Complete(),
		c.OperatorConfig,
	}

	completedCfg.GenericConfig.Version = &version.Info{
		Major: "1",
		Minor: "1",
	}

	return CompletedConfig{&completedCfg}
}

// New returns a new instance of KubedServer from the given config.
func (c completedConfig) New() (*KubedServer, error) {
	genericServer, err := c.GenericConfig.New("kubed-server", genericapiserver.NewEmptyDelegate()) // completion is done in Complete, no need for a second time
	if err != nil {
		return nil, err
	}
	operator, err := c.OperatorConfig.New()
	if err != nil {
		return nil, err
	}

	s := &KubedServer{
		GenericAPIServer: genericServer,
		Operator:         operator,
	}

	return s, nil
}
