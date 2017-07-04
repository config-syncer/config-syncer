package watcher

import (
	"sync"
	"time"

	"github.com/appscode/kubed/pkg/events"
	"github.com/appscode/kubed/pkg/indexers"
	"github.com/appscode/kubed/pkg/namespacesync"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
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

func eventHandlerFuncs(k *Controller) cache.ResourceEventHandlerFuncs {
	return cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			e := events.New(events.Added, k.KubeClient, obj)
			k.Dispatch(e)
		},
		DeleteFunc: func(obj interface{}) {
			e := events.New(events.Deleted, k.KubeClient, obj)
			k.Dispatch(e)
		},
		UpdateFunc: func(old, new interface{}) {
			e := events.New(events.Updated, k.KubeClient, old, new)
			k.Dispatch(e)
		},
	}
}

func (w *Controller) Dispatch(e *events.Event) {
	if e.Ignorable() {
		return
	}

	switch e.ResourceType {
	case events.Namespace:
		ns := namespacesync.NewHandler(w.KubeClient)
		ns.Handle(e)
	case events.Service:
		if w.RunOptions.EnableReverseIndex {
			w.ReverseIndex.Handle(e)
		}
	}
}
