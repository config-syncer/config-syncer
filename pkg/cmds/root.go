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
	"kmodules.xyz/client-go/tools/cli"
)

func NewCmdConfigSyncer(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "kubed",
		Short:             `Config Syncer by AppsCode - A Kubernetes Cluster Operator Daemon`,
		Long:              `Config Syncer is a Kubernetes daemon to perform cluster management tasks. For more information, visit here: https://github.com/kubeops/config-syncer/tree/master/docs`,
		DisableAutoGenTag: true,
		PersistentPreRun: func(c *cobra.Command, args []string) {
			cli.SendAnalytics(c, v.Version.Version)
		},
	}
	cmd.PersistentFlags().BoolVar(&cli.EnableAnalytics, "enable-analytics", cli.EnableAnalytics, "send usage events to Google Analytics")

	stopCh := genericapiserver.SetupSignalHandler()
	cmd.AddCommand(NewCmdRun(os.Stdout, os.Stderr, stopCh))
	cmd.AddCommand(v.NewCmdVersion())

	return cmd
}
