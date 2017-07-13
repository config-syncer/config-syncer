package cmds

import (
	"fmt"
	"os"

	"github.com/appscode/go/runtime"
	"github.com/appscode/kubed/pkg/analytics"
	"github.com/appscode/kubed/pkg/operator"
	"github.com/appscode/log"
	srch_cs "github.com/appscode/searchlight/client/clientset"
	scs "github.com/appscode/stash/client/clientset"
	vcs "github.com/appscode/voyager/client/clientset"
	pcm "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1alpha1"
	kcs "github.com/k8sdb/apimachinery/client/clientset"
	"github.com/spf13/cobra"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func NewCmdRun(version string) *cobra.Command {
	opt := operator.Options{
		ConfigPath:         runtime.GOPath() + "/src/github.com/appscode/kubed/hack/config/clusterconfig.yaml",
		Indexer:            "indexers.bleve",
		EnableReverseIndex: true,
		ServerAddress:      ":8081",
		EnableAnalytics:    true,
		EnableConfigSync:   true,
		ScratchDir:         "/tmp",
	}
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run daemon",
		PreRun: func(cmd *cobra.Command, args []string) {
			if opt.EnableAnalytics {
				analytics.Enable()
			}
			analytics.SendEvent("kubed", "started", version)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			analytics.SendEvent("kubed", "stopped", version)
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Infoln("Starting kubed...")

			Run(opt)
		},
	}

	cmd.Flags().StringVar(&opt.KubeConfig, "kubeconfig", opt.KubeConfig, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	cmd.Flags().StringVar(&opt.Master, "master", opt.Master, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	cmd.Flags().StringVar(&opt.ConfigPath, "clusterconfig", opt.ConfigPath, "Path to cluster config file")

	cmd.Flags().StringVar(&opt.Indexer, "indexer", opt.Indexer, "Reverse indexing of pods to service and others")
	cmd.Flags().BoolVar(&opt.EnableReverseIndex, "enable-reverse-index", opt.EnableReverseIndex, "Reverse indexing of pods to service and others")
	cmd.Flags().StringVar(&opt.ServerAddress, "address", opt.ServerAddress, "The address of the Kubed API Server")

	cmd.Flags().BoolVar(&opt.EnableAnalytics, "analytics", opt.EnableAnalytics, "Send analytical events to Google Analytics")

	return cmd
}

func Run(opt operator.Options) {
	log.Infoln("configurations provided for kubed", opt)
	defer runtime.HandleCrash()

	c, err := clientcmd.BuildConfigFromFlags(opt.Master, opt.KubeConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	op := &operator.Operator{
		KubeClient:        clientset.NewForConfigOrDie(c),
		VoyagerClient:     vcs.NewForConfigOrDie(c),
		SearchlightClient: srch_cs.NewForConfigOrDie(c),
		StashClient:       scs.NewForConfigOrDie(c),
		KubeDBClient:      kcs.NewForConfigOrDie(c),
		Opt:               opt,
	}
	op.PromClient, err = pcm.NewForConfig(c)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = op.Setup()
	if err != nil {
		log.Fatalln(err)
	}

	log.Infoln("Running kubed watcher")
	op.RunAndHold()
}
