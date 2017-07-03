package watcher

import (
	"sync"
	"time"

	"github.com/appscode/kubed/pkg/events"
	"github.com/appscode/kubed/pkg/namespacesync"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type Controller struct {
	// kubernetes client to apiserver
	KubeClient clientset.Interface

	SyncPeriod time.Duration
	sync.Mutex
}

func (w *Controller) Run() {
	w.watchNamespaces()
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

	if e.ResourceType == events.Namespace && e.EventType.IsAdded() {
		ns := namespacesync.NewHandler(w.KubeClient)
		ns.Handle(e)
	}
}
