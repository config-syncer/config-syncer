/*
Copyright The Config Syncer Authors.

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
	"fmt"
	licenseapi "go.bytebuilders.dev/license-verifier/apis/licenses/v1alpha1"
	"go.bytebuilders.dev/license-verifier/info"
	license "go.bytebuilders.dev/license-verifier/kubernetes"
	"k8s.io/apimachinery/pkg/util/sets"
	"time"

	"kubeops.dev/config-syncer/pkg/operator"

	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
)

type OperatorOptions struct {
	ClusterName           string
	ConfigSourceNamespace string
	KubeConfigFile        string

	QPS          float32
	Burst        int
	ResyncPeriod time.Duration

	LicenseFile       string
	LicenseApiService string
}

func NewOperatorOptions() *OperatorOptions {
	return &OperatorOptions{
		ClusterName:           "",
		ConfigSourceNamespace: "",
		KubeConfigFile:        "",
		// ref: https://github.com/kubernetes/ingress-nginx/blob/e4d53786e771cc6bdd55f180674b79f5b692e552/pkg/ingress/controller/launch.go#L252-L259
		// High enough QPS to fit all expected use cases. QPS=0 is not set here, because client code is overriding it.
		QPS: 1e6,
		// High enough Burst to fit all expected use cases. Burst=0 is not set here, because client code is overriding it.
		Burst:        1e6,
		ResyncPeriod: 10 * time.Minute,
	}
}

func (s *OperatorOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.ClusterName, "cluster-name", s.ClusterName, "Name of cluster")
	fs.StringVar(&s.ConfigSourceNamespace, "config-source-namespace", s.ConfigSourceNamespace, "Config source namespace")
	fs.StringVar(&s.KubeConfigFile, "kubeconfig-file", s.KubeConfigFile, "kubeconfig file")

	fs.Float32Var(&s.QPS, "qps", s.QPS, "The maximum QPS to the master from this client")
	fs.IntVar(&s.Burst, "burst", s.Burst, "The maximum burst for throttle")
	fs.DurationVar(&s.ResyncPeriod, "resync-period", s.ResyncPeriod, "If non-zero, will re-list this often. Otherwise, re-list will be delayed aslong as possible (until the upstream source closes the watch or times out.")
	fs.StringVar(&s.LicenseFile, "license-file", s.LicenseFile, "Path to license file")
	fs.StringVar(&s.LicenseApiService, "license-apiservice", s.LicenseApiService, "Name of the ApiService to use by the addons to identify the respective service and certificate for license verification request")
}

func (s *OperatorOptions) ApplyTo(cfg *operator.OperatorConfig) error {
	var err error

	cfg.LicenseFile = s.LicenseFile
	cfg.LicenseApiService = s.LicenseApiService
	cfg.ClientConfig.QPS = s.QPS
	cfg.ClientConfig.Burst = s.Burst
	cfg.ResyncPeriod = s.ResyncPeriod
	cfg.Test = false

	if cfg.KubeClient, err = kubernetes.NewForConfig(cfg.ClientConfig); err != nil {
		return err
	}

	cfg.ClusterName = s.ClusterName
	cfg.ConfigSourceNamespace = s.ConfigSourceNamespace
	cfg.KubeConfigFile = s.KubeConfigFile

	if cfg.LicenseProvided() {
		l := license.MustLicenseEnforcer(cfg.ClientConfig, cfg.LicenseFile).LoadLicense()
		if l.Status != licenseapi.LicenseActive {
			return fmt.Errorf("license status %s, reason: %s", l.Status, l.Reason)
		}
		if !sets.NewString(l.Features...).HasAny(info.Features()...) {
			return fmt.Errorf("not a valid license for this product")
		}
		cfg.License = l
	}

	return nil
}
