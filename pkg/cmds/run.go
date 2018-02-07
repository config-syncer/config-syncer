package cmds

import (
	"io"

	"github.com/appscode/go/log"
	"github.com/appscode/kubed/pkg/cmds/server"
	"github.com/spf13/cobra"
)

// runtime.GOPath() + "/src/github.com/appscode/kubed/hack/config/clusterconfig.yaml"
func NewCmdRun(out, errOut io.Writer, stopCh <-chan struct{}) *cobra.Command {
	o := server.NewKubedOptions(out, errOut)

	cmd := &cobra.Command{
		Use:               "run",
		Short:             "Launch Kubernetes Cluster Daemon",
		Long:              "Launch Kubernetes Cluster Daemon",
		DisableAutoGenTag: true,
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
