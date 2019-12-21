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
	"fmt"

	"github.com/appscode/go/flags"
	api "github.com/appscode/kubed/apis/kubed/v1alpha1"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

func NewCmdCheck() *cobra.Command {
	var (
		configPath string
	)

	cmd := &cobra.Command{
		Use:               "check",
		Short:             "Check cluster config",
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.EnsureRequiredFlags(cmd, "clusterconfig")

			cfg, err := api.LoadConfig(configPath)
			if err != nil {
				return err
			}
			err = cfg.Validate()
			if err != nil {
				return err
			}
			data, err := yaml.Marshal(cfg)
			if err != nil {
				return err
			}
			fmt.Println("Cluster config was parsed successfully.")
			fmt.Println()
			fmt.Println(string(data))
			return nil
		},
	}
	cmd.Flags().StringVar(&configPath, "clusterconfig", configPath, "Path to cluster config file")
	return cmd
}
