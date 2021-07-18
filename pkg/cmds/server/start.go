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
	"io"
	"net"

	"kubeops.dev/kubed/pkg/operator"
	"kubeops.dev/kubed/pkg/server"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	"kmodules.xyz/client-go/tools/clientcmd"
)

const defaultEtcdPathPrefix = "/registry/kubed.appscode.com"

type KubedOptions struct {
	RecommendedOptions *genericoptions.RecommendedOptions
	OperatorOptions    *OperatorOptions

	StdOut io.Writer
	StdErr io.Writer
}

func NewKubedOptions(out, errOut io.Writer) *KubedOptions {
	o := &KubedOptions{
		// TODO we will nil out the etcd storage options.  This requires a later level of k8s.io/apiserver
		RecommendedOptions: genericoptions.NewRecommendedOptions(
			defaultEtcdPathPrefix,
			server.Codecs.LegacyCodec(),
		),
		OperatorOptions: NewOperatorOptions(),
		StdOut:          out,
		StdErr:          errOut,
	}
	o.RecommendedOptions.Etcd = nil
	o.RecommendedOptions.Admission = nil

	return o
}

func (o *KubedOptions) AddFlags(fs *pflag.FlagSet) {
	o.RecommendedOptions.AddFlags(fs)
	o.OperatorOptions.AddFlags(fs)
}

func (o KubedOptions) Validate(args []string) error {
	return nil
}

func (o *KubedOptions) Complete() error {
	return nil
}

func (o KubedOptions) Config() (*server.KubedConfig, error) {
	// TODO have a "real" external address
	if err := o.RecommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return nil, errors.Errorf("error creating self-signed certificates: %v", err)
	}

	serverConfig := genericapiserver.NewRecommendedConfig(server.Codecs)
	if err := o.RecommendedOptions.ApplyTo(serverConfig); err != nil {
		return nil, err
	}
	clientcmd.Fix(serverConfig.ClientConfig)

	operatorConfig := operator.NewOperatorConfig(serverConfig.ClientConfig)
	if err := o.OperatorOptions.ApplyTo(operatorConfig); err != nil {
		return nil, err
	}

	config := &server.KubedConfig{
		GenericConfig:  serverConfig,
		OperatorConfig: operatorConfig,
	}
	return config, nil
}

func (o KubedOptions) Run(stopCh <-chan struct{}) error {
	config, err := o.Config()
	if err != nil {
		return err
	}

	s, err := config.Complete().New()
	if err != nil {
		return err
	}

	return s.Run(stopCh)
}
