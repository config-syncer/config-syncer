package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/appscode/go/runtime"
	"github.com/appscode/go/wait"

	"appscode.com/kubed/pkg/janitor"
	"appscode.com/kubed/pkg/watcher"
	"github.com/appscode/client"
	"github.com/appscode/errors"
	_ "github.com/appscode/k8s-addons/api/install"
	acs "github.com/appscode/k8s-addons/client/clientset"
	"github.com/appscode/k8s-addons/pkg/dns"
	acw "github.com/appscode/k8s-addons/pkg/watcher"
	"github.com/appscode/log"
	"github.com/appscode/searchlight/pkg/client/icinga"
	"github.com/appscode/searchlight/pkg/client/influxdb"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
	clientcmdapi "k8s.io/kubernetes/pkg/client/unversioned/clientcmd/api"
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
}

func Run(config *Config) {
	log.Infoln("configurations provided for kubed", config)
	defer runtime.HandleCrash()

	// ref; https://github.com/kubernetes/kubernetes/blob/ba1666fb7b946febecfc836465d22903b687118e/cmd/kube-proxy/app/server.go#L168
	// Create a Kube Client
	// define api config source
	if config.KubeConfig == "" && config.Master == "" {
		log.Warningf("Neither --kubeconfig nor --master was specified.  Using default API client.  This might not work.")
	}
	// This creates a client, first loading any specified kubeconfig
	// file, and then overriding the Master flag, if non-empty.
	c, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: config.KubeConfig},
		&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: config.Master}}).ClientConfig()
	if err != nil {
		errors.Exit(err)
	}

	apiOptions := client.NewOption(config.APIEndpoint)
	log.Infoln("api options", apiOptions)
	apiOptions.BearerAuth(readAuth(config.APITokenPath))
	ac := acs.NewACExtensionsForConfigOrDie(c)

	kubeWatcher := &watcher.KubedWatcher{
		Watcher: acw.Watcher{
			Client:                  clientset.NewForConfigOrDie(c),
			AppsCodeExtensionClient: ac,
			SyncPeriod:              time.Minute * 2,
		},
		AppsCodeApiClientOptions: apiOptions,
		ProviderName:             config.ProviderName,
		ClusterName:              config.ClusterName,
		LoadbalancerImage:        config.LoadbalancerImageName,
	}

	// Icinga client
	icingaClient, err := icinga.NewInClusterClient(kubeWatcher.Client)
	if err != nil {
		log.Errorln(err)
	}
	kubeWatcher.IcingaClient = icingaClient

	log.Infoln("configuration loadded, running kubed watcher")
	go kubeWatcher.Run()

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

	// InfluxDB client
	influxConfig, err := influxdb.LoadConfig(kubeWatcher.Client)
	if err != nil {
		log.Errorln(err)
	}
	kubeJanitor.InfluxConfig = *influxConfig

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
