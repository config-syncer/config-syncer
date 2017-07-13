package cmds

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/appscode/go/runtime"
	"github.com/appscode/kubed/pkg/analytics"
	"github.com/appscode/kubed/pkg/config"
	"github.com/appscode/kubed/pkg/indexers"
	"github.com/appscode/kubed/pkg/recover"
	"github.com/appscode/kubed/pkg/watcher"
	"github.com/appscode/log"
	"github.com/appscode/pat"
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
	opt := watcher.Options{
		Indexer:            "indexers.bleve",
		EnableReverseIndex: true,
		ServerAddress:      ":32600",
		EnableAnalytics:    true,
		ConfigPath:         runtime.GOPath() + "/src/github.com/appscode/kubed/hack/config/clusterconfig.yaml",
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

func Run(opt watcher.Options) {
	log.Infoln("configurations provided for kubed", opt)
	defer runtime.HandleCrash()

	c, err := clientcmd.BuildConfigFromFlags(opt.Master, opt.KubeConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cfg, err := config.LoadConfig(opt.ConfigPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	w := &watcher.Watchers{
		KubeClient:        clientset.NewForConfigOrDie(c),
		VoyagerClient:     vcs.NewForConfigOrDie(c),
		SearchlightClient: srch_cs.NewForConfigOrDie(c),
		StashClient:       scs.NewForConfigOrDie(c),
		KubeDBClient:      kcs.NewForConfigOrDie(c),

		Opt:    opt,
		Config: *cfg,
		Saver: &recover.RecoverStuff{
			Opt: cfg.Recover,
		},
		SyncPeriod: time.Minute * 2,
	}
	w.PromClient, err = pcm.NewForConfig(c)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// router is default HTTP request multiplexer for kubed. It matches the URL of each
	// incoming request against a list of registered patterns with their associated
	// methods and calls the handler for the pattern that most closely matches the
	// URL.
	//
	// Pattern matching attempts each pattern in the order in which they were
	// registered.
	router := pat.New()

	// Enable full text indexing to have search feature
	if len(opt.Indexer) > 0 {
		indexer, err := indexers.NewResourceIndexer(opt.Indexer)
		if err != nil {
			log.Errorln(err)
		} else {
			indexer.RegisterRouters(router)
			w.Indexer = indexer
		}
	}

	// Enable pod -> service, service -> serviceMonitor indexing
	if opt.EnableReverseIndex {
		ri, err := indexers.NewReverseIndexer(w.KubeClient, opt.Indexer)
		if err != nil {
			log.Errorln("Failed to create indexer", err)
		} else {
			ri.RegisterRouters(router)
			w.ReverseIndex = ri
		}
	}

	log.Infoln("Running kubed watcher")
	go w.Run()

	// TODO: periodic backup
	// TODO: janitor

	http.Handle("/", router)
	log.Fatalln(http.ListenAndServe(opt.ServerAddress, nil))
}
