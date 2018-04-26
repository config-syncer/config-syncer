package cmds

import (
	"fmt"

	"github.com/appscode/go/flags"
	"github.com/appscode/kutil/tools/backup"
	"github.com/appscode/kutil/tools/clientcmd"
	"github.com/spf13/cobra"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
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
