package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/appscode/go/hold"
	"github.com/appscode/go/runtime"
	"github.com/appscode/go/wait"
	"github.com/appscode/kubed/pkg/dns"
	"github.com/appscode/kubed/pkg/indexers"
	"github.com/appscode/kubed/pkg/janitor"
	"github.com/appscode/kubed/pkg/watcher"
	"github.com/appscode/log"
	"github.com/appscode/pat"
	"github.com/appscode/searchlight/pkg/client/influxdb"
	"github.com/spf13/cobra"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type RunOptions struct {
	Master                            string
	KubeConfig                        string
	ESEndpoint                        string
	InfluxSecretName                  string
	InfluxSecretNamespace             string
	ClusterName                       string
	ClusterKubedConfigSecretName      string
	ClusterKubedConfigSecretNamespace string
	ReverseIndex                      bool
	ServerAddress                     string
}

func NewCmdRun() *cobra.Command {
	opt := RunOptions{
		InfluxSecretName:                  "appscode-influx",
		InfluxSecretNamespace:             "kube-system",
		ClusterKubedConfigSecretName:      "cluster-kubed-config",
		ClusterKubedConfigSecretNamespace: "kube-system",
		ReverseIndex:                      true,
		ServerAddress:                     ":32600",
	}
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run daemon",
		Run: func(cmd *cobra.Command, args []string) {
			if opt.ClusterName == "" {
				log.Fatalln("Missing required flag: --cluster-name")
			}
			log.Infoln("Starting kubed...")
			go Run(opt)

			hold.Hold()
		},
	}

	cmd.Flags().StringVar(&opt.ClusterName, "cluster-name", opt.ClusterName, "Name of Kubernetes cluster")
	cmd.Flags().StringVar(&opt.ESEndpoint, "es-endpoint", opt.ESEndpoint, "Endpoint of elasticsearch")
	cmd.Flags().StringVar(&opt.InfluxSecretName, "influx-secret", opt.InfluxSecretName, "Influxdb secret name")
	cmd.Flags().StringVar(&opt.ClusterKubedConfigSecretName, "kubed-config-secret-name", opt.ClusterKubedConfigSecretName, "Kubed configuration secret name")
	cmd.Flags().StringVar(&opt.ClusterKubedConfigSecretNamespace, "kubed-config-secret-namespace", opt.ClusterKubedConfigSecretNamespace, "Kubed configuration secret namespace")
	cmd.Flags().StringVar(&opt.InfluxSecretNamespace, "influx-secret-namespace", opt.InfluxSecretNamespace, "Influxdb secret namespace")
	cmd.Flags().StringVar(&opt.KubeConfig, "kubeconfig", opt.KubeConfig, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	cmd.Flags().StringVar(&opt.Master, "master", opt.Master, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	cmd.Flags().BoolVar(&opt.ReverseIndex, "reverse-index", opt.ReverseIndex, "Reverse indexing of pods to service and others")
	cmd.Flags().StringVar(&opt.ServerAddress, "address", opt.ServerAddress, "The address of the Kubed API Server")
	return cmd
}

func Run(opt RunOptions) {
	log.Infoln("configurations provided for kubed", opt)
	defer runtime.HandleCrash()

	c, err := clientcmd.BuildConfigFromFlags(opt.Master, opt.KubeConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	kubeWatcher := &watcher.Watcher{
		KubeClient: clientset.NewForConfigOrDie(c),
		SyncPeriod: time.Minute * 2,
	}

	log.Infoln("Running kubed watcher")
	go kubeWatcher.Run()

	// router is default HTTP request multiplexer for kubed. It matches the URL of each
	// incoming request against a list of registered patterns with their associated
	// methods and calls the handler for the pattern that most closely matches the
	// URL.
	//
	// Pattern matching attempts each pattern in the order in which they were
	// registered.
	router := pat.New()

	if opt.ReverseIndex {
		ri := indexers.NewReverseIndexer(clientset.NewForConfigOrDie(c), time.Second*2)
		ri.RegisterRouters(router)
		go ri.Start()
	}

	// initializing kube janitor tasks
	kubeJanitor := janitor.Janitor{
		KubeClient:                        clientset.NewForConfigOrDie(c),
		ClusterName:                       opt.ClusterName,
		ElasticConfig:                     make(map[string]string),
		ClusterKubedConfigSecretName:      opt.ClusterKubedConfigSecretName,
		ClusterKubedConfigSecretNamespace: opt.ClusterKubedConfigSecretNamespace,
	}

	if opt.ESEndpoint != "" {
		endpoint := opt.ESEndpoint
		if strings.HasPrefix(opt.ESEndpoint, "http") {
			endpoint = opt.ESEndpoint[7:]
		}
		parts := strings.Split(endpoint, ":")
		if len(parts) == 2 {
			esServiceClusterIP, err := dns.GetServiceClusterIP(kubeWatcher.KubeClient, "ES", parts[0])
			if err != nil {
				log.Errorln(err)
			} else {
				kubeJanitor.ElasticConfig[janitor.ESEndpoint] = fmt.Sprintf("http://%v:%v", esServiceClusterIP, parts[1])
			}
		} else {
			log.Errorln("es-endpoint should contain host_name & host_port")
		}
	}

	if opt.InfluxSecretName != "" {
		// InfluxDB client
		influxConfig, err := influxdb.GetInfluxDBConfig(opt.InfluxSecretName, opt.InfluxSecretNamespace)
		if err != nil {
			log.Errorln(err)
		} else {
			kubeJanitor.InfluxConfig = *influxConfig
		}
	}
	go wait.Forever(kubeJanitor.Run, time.Hour*24)

	if len(opt.ServerAddress) > 0 {
		http.Handle("/", router)
		go http.ListenAndServe(opt.ServerAddress, nil)
	}
}
