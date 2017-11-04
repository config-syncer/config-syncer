package cmds

import (
	"fmt"
	"github.com/appscode/go/flags"
	"github.com/appscode/kutil"
	"github.com/spf13/cobra"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewCmdSnapshot() *cobra.Command {
	var (
		sanitize  bool
		backupDir string
		context   string
	)
	cmd := &cobra.Command{
		Use:               "snapshot",
		Short:             "Takes a snapshot of Kubernetes api objects",
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.EnsureRequiredFlags(cmd, "cluster", "backup-dir")

			restConfig, err := createKubeConfig(context)
			if err != nil {
				return err
			}
			mgr := kutil.NewBackupManager(context, restConfig, sanitize)
			filename, err := mgr.BackupToTar(backupDir)
			if err != nil {
				return err
			}
			fmt.Printf("Cluster objects are stored in %s", filename)
			fmt.Println()
		},
	}
	cmd.Flags().BoolVar(&sanitize, "sanitize", false, " Sanitize fields in YAML")
	cmd.Flags().StringVar(&backupDir, "backup-dir", "", "Directory where YAML files will be stored")
	cmd.Flags().StringVar(&context, "context", "", "Name of the kubeconfig context to use")
	return cmd
}

func createKubeConfig(ctx string) (*rest.Config, error) {
	apiConfig, err := clientcmd.NewDefaultPathOptions().GetStartingConfig()
	if err != nil {
		return nil, err
	}
	overrides := &clientcmd.ConfigOverrides{CurrentContext: ctx}
	return clientcmd.NewDefaultClientConfig(*apiConfig, overrides).ClientConfig()
}
