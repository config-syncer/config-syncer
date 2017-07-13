package cmds

import (
	"os"

	"github.com/appscode/go/flags"
	"github.com/appscode/kubed/pkg/backup"
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
		Use:   "snapshot",
		Short: "Takes a snapshot of Kubernetes api objects",
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.EnsureRequiredFlags(cmd, "context", "backup-dir")

			err := os.MkdirAll(backupDir, 0777)
			if err != nil {
				return err
			}
			restConfig, err := createKubeConfig(context)
			if err != nil {
				return err
			}
			return backup.SnapshotCluster(restConfig, backupDir, sanitize)
		},
	}
	cmd.Flags().BoolVar(&sanitize, "sanitize", false, " Sanitize fields in YAML")
	cmd.Flags().StringVar(&backupDir, "backup-dir", "", "Directory where YAML files will be stored")
	cmd.Flags().StringVar(&context, "context", "", "The name of the kubeconfig context to use")
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
