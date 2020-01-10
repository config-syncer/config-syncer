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
	"io"

	"github.com/appscode/go/log"
	v "github.com/appscode/go/version"
	"github.com/appscode/kubed/pkg/cmds/server"

	"github.com/spf13/cobra"
	"kmodules.xyz/client-go/tools/cli"
)

// runtime.GOPath() + "/src/github.com/appscode/kubed/hack/config/clusterconfig.yaml"
func NewCmdRun(out, errOut io.Writer, stopCh <-chan struct{}) *cobra.Command {
	o := server.NewKubedOptions(out, errOut)

	cmd := &cobra.Command{
		Use:               "run",
		Short:             "Launch Kubernetes Cluster Daemon",
		Long:              "Launch Kubernetes Cluster Daemon",
		DisableAutoGenTag: true,
		PreRun: func(c *cobra.Command, args []string) {
			cli.SendAnalytics(c, v.Version.Version)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Infoln("Starting kubed...")

			if err := o.Complete(); err != nil {
				return err
			}
			if err := o.Validate(args); err != nil {
				return err
			}
			if err := o.Run(stopCh); err != nil {
				return err
			}
			return nil
		},
	}

	o.AddFlags(cmd.Flags())

	return cmd
}
