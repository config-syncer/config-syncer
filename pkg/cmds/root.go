package cmds

import (
	"flag"
	"log"

	v "github.com/appscode/go/version"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	_ "k8s.io/client-go/kubernetes/fake"
)

func NewCmdKubed(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubed",
		Short: `Kubed by AppsCode - A Kubernetes Cluster Operator Daemon`,
		Long:  `Kubed is a Kubernetes daemon to perform cluster management tasks. For more information, visit here: https://github.com/appscode/kubed/tree/master/docs`,
		PersistentPreRun: func(c *cobra.Command, args []string) {
			c.Flags().VisitAll(func(flag *pflag.Flag) {
				log.Printf("FLAG: --%s=%q", flag.Name, flag.Value)
			})
		},
	}
	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	// ref: https://github.com/kubernetes/kubernetes/issues/17162#issuecomment-225596212
	flag.CommandLine.Parse([]string{})

	cmd.AddCommand(NewCmdRun(version))
	cmd.AddCommand(NewCmdSnapshot())
	cmd.AddCommand(NewCmdCheck())
	cmd.AddCommand(v.NewCmdVersion())

	return cmd
}
