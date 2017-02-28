package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"appscode.com/kubed/pkg/janitor"
	promwatcher "appscode.com/kubed/pkg/promwatcher"
	"appscode.com/kubed/pkg/watcher"
	"github.com/appscode/client"
	"github.com/appscode/errors"
	"github.com/appscode/go/runtime"
	"github.com/appscode/go/wait"
	_ "github.com/appscode/k8s-addons/api/install"
	acs "github.com/appscode/k8s-addons/client/clientset"
	"github.com/appscode/k8s-addons/pkg/dns"
	acw "github.com/appscode/k8s-addons/pkg/watcher"
	"github.com/appscode/log"
	"github.com/appscode/searchlight/pkg/client/icinga"
	"github.com/appscode/searchlight/pkg/client/influxdb"
	pcm "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1alpha1"
	cgcmd "k8s.io/client-go/tools/clientcmd"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
)

type Config struct {
	APITokenPath          string
	APIEndpoint           string
	ProviderName          string
	ClusterName           string
	LoadbalancerImageName string
	Master                string
	KubeConfig            string
	ESEndpoint            string
	InfluxSecretName      string
	InfluxSecretNamespace string
	IcingaSecretName      string
	IcingaSecretNamespace string
	IngressClass          string
	EnablePromMonitoring  bool
}

func Run(config *Config) {
	log.Infoln("configurations provided for kubed", config)
	defer runtime.HandleCrash()

	c, err := clientcmd.BuildConfigFromFlags(config.Master, config.KubeConfig)
	if err != nil {
		errors.Exit(err)
	}

	apiOptions := client.NewOption(config.APIEndpoint)
	log.Infoln("api options", apiOptions)
	apiOptions.BearerAuth(readAuth(config.APITokenPath))

	kubeWatcher := &watcher.KubedWatcher{
		Watcher: acw.Watcher{
			Client:                  clientset.NewForConfigOrDie(c),
			AppsCodeExtensionClient: acs.NewACExtensionsForConfigOrDie(c),
			SyncPeriod:              time.Minute * 2,
		},
		AppsCodeApiClientOptions: apiOptions,
		ProviderName:             config.ProviderName,
		ClusterName:              config.ClusterName,
		LoadbalancerImage:        config.LoadbalancerImageName,
		IngressClass:             config.IngressClass,
	}

	if config.IcingaSecretName != "" {
		// Icinga client
		icingaClient, err := icinga.NewIcingaClient(kubeWatcher.Client, config.IcingaSecretName, config.IcingaSecretNamespace)
		if err != nil {
			log.Errorln(err)
		}
		kubeWatcher.IcingaClient = icingaClient
	}

	log.Infoln("configuration loadded, running kubed watcher")
	go kubeWatcher.Run()

	if config.EnablePromMonitoring {
		// get rest.Config for "k8s.io/client-go/tools/clientcmd" package
		config, err := cgcmd.BuildConfigFromFlags(config.Master, config.KubeConfig)
		if err == nil {
			// get client for Prometheus TPR monitoring
			client, err := pcm.NewForConfig(config)
			if err != nil {
				log.Fatalln(err)
			}
			watcher := &promwatcher.PromWatcher{
				Watcher:    kubeWatcher.Watcher,
				PromClient: client,
				SyncPeriod: time.Minute * 2,
			}
			log.Infoln("running Prometheus watcher")
			watcher.WatchPrometheus()
		} else {
			log.Fatalln(err)
		}
	}

	// initializing kube janitor tasks
	kubeJanitor := janitor.Janitor{
		KubedWatcher:  kubeWatcher,
		ElasticConfig: make(map[string]string),
	}

	if config.ESEndpoint != "" {
		endpoint := config.ESEndpoint
		if strings.HasPrefix(config.ESEndpoint, "http") {
			endpoint = config.ESEndpoint[7:]
		}
		parts := strings.Split(endpoint, ":")
		if len(parts) == 2 {
			esServiceClusterIP, err := dns.GetServiceClusterIP(kubeWatcher.Client, "ES", parts[0])
			if err != nil {
				log.Errorln(err)
			} else {
				kubeJanitor.ElasticConfig[janitor.ESEndpoint] = fmt.Sprintf("http://%v:%v", esServiceClusterIP, parts[1])
			}
		} else {
			log.Errorln("es-endpoint should contain host_name & host_port")
		}
	}

	if config.InfluxSecretName != "" {
		// InfluxDB client
		influxConfig, err := influxdb.GetInfluxDBConfig(config.InfluxSecretName, config.InfluxSecretNamespace)
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
