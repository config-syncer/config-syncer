package watcher

import (
	"sync"
	"time"

	"github.com/appscode/kubed/pkg/indexers"
	scs "github.com/appscode/stash/client/clientset"
	vcs "github.com/appscode/voyager/client/clientset"
	pcm "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1alpha1"
	kcs "github.com/k8sdb/apimachinery/client/clientset"
	clientset "k8s.io/client-go/kubernetes"
)

type Controller struct {
	KubeClient    clientset.Interface
	VoyagerClient vcs.ExtensionInterface
	StashClient   scs.ExtensionInterface
	PromClient    pcm.MonitoringV1alpha1Interface
	KubeDBClient  kcs.ExtensionInterface

	RunOptions   RunOptions
	Indexer      *indexers.ResourceIndexer
	ReverseIndex *indexers.ReverseIndexer
	SyncPeriod   time.Duration
	sync.Mutex
}

type RunOptions struct {
	Master                            string
	KubeConfig                        string
	ESEndpoint                        string
	InfluxSecretName                  string
	InfluxSecretNamespace             string
	ClusterName                       string
	ClusterKubedConfigSecretName      string
	ClusterKubedConfigSecretNamespace string
	Indexer                           string
	EnableReverseIndex                bool
	ServerAddress                     string
	NotifyOnCertSoonToBeExpeired      bool
	NotifyVia                         string
}

func (w *Controller) Run() {
	w.watchNamespaces()
	if w.RunOptions.EnableReverseIndex || len(w.RunOptions.Indexer) > 0 {
		w.watchService()
	}
}
