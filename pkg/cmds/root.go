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

package cmds

import (
	"flag"
	"os"

	"github.com/appscode/go/flags"
	v "github.com/appscode/go/version"

	"github.com/spf13/cobra"
	genericapiserver "k8s.io/apiserver/pkg/server"
	_ "k8s.io/client-go/kubernetes/fake"
	"kmodules.xyz/client-go/logs"
	"kmodules.xyz/client-go/tools/cli"
)

func NewCmdKubed(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "kubed",
		Short:             `Kubed by AppsCode - A Kubernetes Cluster Operator Daemon`,
		Long:              `Kubed is a Kubernetes daemon to perform cluster management tasks. For more information, visit here: https://github.com/appscode/kubed/tree/master/docs`,
		DisableAutoGenTag: true,
		PersistentPreRun: func(c *cobra.Command, args []string) {
			flags.DumpAll(c.Flags())
			cli.SendAnalytics(c, v.Version.Version)
		},
	}
	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	logs.ParseFlags()
	cmd.PersistentFlags().BoolVar(&cli.EnableAnalytics, "enable-analytics", cli.EnableAnalytics, "send usage events to Google Analytics")

	stopCh := genericapiserver.SetupSignalHandler()
	cmd.AddCommand(NewCmdRun(os.Stdout, os.Stderr, stopCh))
	cmd.AddCommand(v.NewCmdVersion())

	return cmd
}
