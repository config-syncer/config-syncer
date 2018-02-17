package cmds

import (
	"flag"
	"log"
	"os"
	"strings"

	v "github.com/appscode/go/version"
	prom_util "github.com/appscode/kube-mon/prometheus/v1"
	"github.com/appscode/kutil/tools/analytics"
	searchlightcheme "github.com/appscode/searchlight/client/clientset/versioned/scheme"
	stashscheme "github.com/appscode/stash/client/clientset/versioned/scheme"
	voyagerscheme "github.com/appscode/voyager/client/clientset/versioned/scheme"
	"github.com/jpillora/go-ogle-analytics"
	kubedbscheme "github.com/kubedb/apimachinery/client/clientset/versioned/scheme"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	genericapiserver "k8s.io/apiserver/pkg/server"
	_ "k8s.io/client-go/kubernetes/fake"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
)

const (
	gaTrackingCode = "UA-62096468-20"
)

func NewCmdKubed(version string) *cobra.Command {
	var (
		enableAnalytics = true
	)
	cmd := &cobra.Command{
		Use:               "kubed",
		Short:             `Kubed by AppsCode - A Kubernetes Cluster Operator Daemon`,
		Long:              `Kubed is a Kubernetes daemon to perform cluster management tasks. For more information, visit here: https://github.com/appscode/kubed/tree/master/docs`,
		DisableAutoGenTag: true,
		PersistentPreRun: func(c *cobra.Command, args []string) {
			c.Flags().VisitAll(func(flag *pflag.Flag) {
				log.Printf("FLAG: --%s=%q", flag.Name, flag.Value)
			})
			if enableAnalytics && gaTrackingCode != "" {
				if client, err := ga.NewClient(gaTrackingCode); err == nil {
					client.ClientID(analytics.ClientID())
					parts := strings.Split(c.CommandPath(), " ")
					client.Send(ga.NewEvent(parts[0], strings.Join(parts[1:], "/")).Label(version))
				}
			}
			voyagerscheme.AddToScheme(clientsetscheme.Scheme)
			searchlightcheme.AddToScheme(clientsetscheme.Scheme)
			stashscheme.AddToScheme(clientsetscheme.Scheme)
			kubedbscheme.AddToScheme(clientsetscheme.Scheme)
			prom_util.AddToScheme(clientsetscheme.Scheme)
		},
	}
	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	// ref: https://github.com/kubernetes/kubernetes/issues/17162#issuecomment-225596212
	flag.CommandLine.Parse([]string{})
	cmd.PersistentFlags().BoolVar(&enableAnalytics, "analytics", enableAnalytics, "Send analytical events to Google Analytics")

	stopCh := genericapiserver.SetupSignalHandler()
	cmd.AddCommand(NewCmdRun(os.Stdout, os.Stderr, stopCh))
	cmd.AddCommand(NewCmdSnapshot())
	cmd.AddCommand(NewCmdCheck())
	cmd.AddCommand(v.NewCmdVersion())

	return cmd
}
