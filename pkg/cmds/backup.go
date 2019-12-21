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

	"github.com/spf13/cobra"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"kmodules.xyz/client-go/tools/backup"
	"kmodules.xyz/client-go/tools/clientcmd"
)

func NewCmdBackup() *cobra.Command {
	var (
		sanitize       bool
		backupDir      string
		kubeconfigPath string
		context        string
	)
	cmd := &cobra.Command{
		Use:               "backup",
		Short:             "Takes a backup of Kubernetes api objects",
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.EnsureRequiredFlags(cmd, "backup-dir")

			restConfig, err := clientcmd.BuildConfigFromContext(kubeconfigPath, context)
			if err != nil {
				return err
			}
			mgr := backup.NewBackupManager(context, restConfig, sanitize)
			filename, err := mgr.BackupToTar(backupDir)
			if err != nil {
				return err
			}
			fmt.Printf("Cluster objects are stored in %s", filename)
			fmt.Println()
			return nil
		},
	}
	cmd.Flags().BoolVar(&sanitize, "sanitize", false, " Sanitize fields in YAML")
	cmd.Flags().StringVar(&backupDir, "backup-dir", "", "Directory where YAML files will be stored")
	cmd.Flags().StringVar(&kubeconfigPath, "kubeconfig", "", "kubeconfig file pointing at the 'core' kubernetes server")
	cmd.Flags().StringVar(&context, "context", "", "Name of the kubeconfig context to use")
	return cmd
}
