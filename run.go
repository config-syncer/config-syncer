package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/appscode/client"
	"github.com/appscode/go/hold"
	"github.com/appscode/go/runtime"
	"github.com/appscode/go/wait"
	"github.com/appscode/kubed/pkg/dns"
	"github.com/appscode/kubed/pkg/janitor"
	"github.com/appscode/kubed/pkg/watcher"
	"github.com/appscode/log"
	"github.com/appscode/searchlight/pkg/client/influxdb"
	"github.com/spf13/cobra"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type RunOptions struct {
	APITokenPath          string
	APIEndpoint           string
	Master                string
	KubeConfig            string
	ESEndpoint            string
	InfluxSecretName      string
	InfluxSecretNamespace string
	ClusterName           string
}

func NewCmdRun() *cobra.Command {
	opt := RunOptions{
		APIEndpoint:           "https://api.appscode.com:3443",
		APITokenPath:          "/var/run/secrets/appscode/api-token",
		InfluxSecretName:      "appscode-influx",
		InfluxSecretNamespace: "kube-system",
	}
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run daemon",
		Run: func(cmd *cobra.Command, args []string) {
			if opt.APITokenPath == "" {
				log.Fatalln("Missing required flag: --api-token")
			}
			if opt.APIEndpoint == "" {
				log.Fatalln("Missing required flag: --api-endpoint")
			}
			if opt.ClusterName == "" {
				log.Fatalln("Missing required flag: --cluster-name")
			}
			log.Infoln("Starting kubed...")
			go Run(opt)

			hold.Hold()
		},
	}

	cmd.Flags().StringVar(&opt.APIEndpoint, "api-endpoint", opt.APIEndpoint, "AppsCode api server address host:port")
	cmd.Flags().StringVar(&opt.APITokenPath, "api-token", opt.APITokenPath, "File path for AppsCode api token.")
	cmd.Flags().StringVar(&opt.ClusterName, "cluster-name", opt.ClusterName, "Name of Kubernetes cluster")
	cmd.Flags().StringVar(&opt.ESEndpoint, "es-endpoint", opt.ESEndpoint, "Endpoint of elasticsearch")
	cmd.Flags().StringVar(&opt.InfluxSecretName, "influx-secret", opt.InfluxSecretName, "Influxdb secret name")
	cmd.Flags().StringVar(&opt.InfluxSecretNamespace, "influx-secret-namespace", opt.InfluxSecretNamespace, "Influxdb secret namespace")
	cmd.Flags().StringVar(&opt.KubeConfig, "kubeconfig", opt.KubeConfig, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	cmd.Flags().StringVar(&opt.Master, "master", opt.Master, "The address of the Kubernetes API server (overrides any value in kubeconfig)")

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

	apiOptions := client.NewOption(opt.APIEndpoint)
	log.Infoln("api options", apiOptions)
	apiOptions.BearerAuth(readAuth(opt.APITokenPath))

	kubeWatcher := &watcher.Watcher{
		KubeClient: clientset.NewForConfigOrDie(c),
		SyncPeriod: time.Minute * 2,
	}

	log.Infoln("configuration loadded, running kubed watcher")
	go kubeWatcher.Run()

	// initializing kube janitor tasks
	kubeJanitor := janitor.Janitor{
		KubeClient:       clientset.NewForConfigOrDie(c),
		ClusterName:      opt.ClusterName,
		APIClientOptions: apiOptions,
		ElasticConfig:    make(map[string]string),
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
}

func readAuth(path string) (string, string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("failed to read api token", err)
	}

	// trying to parse the api token.
	type Token struct {
		Namespace string `json:"namespace,omitempty"`
		Token     string `json:"token,omitempty"`
	}
	a := &Token{}
	err = json.Unmarshal(data, a)
	if err != nil {
		log.Fatalln("failed to masrshel auth data", err)
	}
	log.Debugln("got api credentials for", a.Namespace, "to", a.Token)
	return a.Namespace, a.Token
}
