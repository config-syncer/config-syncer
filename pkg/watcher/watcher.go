package watcher

import (
	"sync"
	"time"

	"github.com/appscode/kubed/pkg/indexers"
	clientset "k8s.io/client-go/kubernetes"
)

type Controller struct {
	// kubernetes client to apiserver
	KubeClient   clientset.Interface
	RunOptions   RunOptions
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
}

func (w *Controller) Run() {
	w.watchNamespaces()
	if w.RunOptions.EnableReverseIndex {
		w.watchService()
	}
}
