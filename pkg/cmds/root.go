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

package cmds

import (
	"os"

	"github.com/spf13/cobra"
	v "gomodules.xyz/x/version"
	genericapiserver "k8s.io/apiserver/pkg/server"
	_ "k8s.io/client-go/kubernetes/fake"
)

func NewCmdConfigSyncer(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "config-syncer",
		Short:             `Config Syncer by AppsCode - A Kubernetes Configuration Syncer`,
		Long:              `Config Syncer is a Kubernetes controller to sync configmaps and secrets. For more information, visit here: https://github.com/kubeops/config-syncer/tree/master/docs`,
		DisableAutoGenTag: true,
	}

	stopCh := genericapiserver.SetupSignalHandler()
	cmd.AddCommand(NewCmdRun(os.Stdout, os.Stderr, stopCh))
	cmd.AddCommand(v.NewCmdVersion())

	return cmd
}
