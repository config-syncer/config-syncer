package watcher

import (
	"reflect"
	"sync"
	"time"

	"github.com/appscode/kubed/pkg/events"
	"github.com/appscode/kubed/pkg/handlers"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/cache"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/util/wait"
)

type Watcher struct {
	// kubernetes client to apiserver
	KubeClient clientset.Interface

	SyncPeriod time.Duration
	sync.Mutex
}

func (w *Watcher) Run() {
	w.watchNamespaces()
}

func (w *Watcher) watchNamespaces() {
	lw := cache.NewListWatchFromClient(w.KubeClient.Core().RESTClient(), events.Namespace.String(), kapi.NamespaceAll, fields.Everything())
	_, controller := cache.NewInformer(lw, &kapi.Namespace{}, w.SyncPeriod, eventHandlerFuncs(w))
	go controller.Run(wait.NeverStop)
}

func (w *Watcher) Dispatch(e *events.Event) error {
	if e.Ignorable() {
		return nil
	}
	if e.ResourceType == events.Namespace && e.EventType == events.Added {
		h := &handlers.NamespaceHandler{
			KubeClient: w.KubeClient,
		}
		h.Handle(e)
	}
	return nil
}

func eventHandlerFuncs(k *Watcher) cache.ResourceEventHandlerFuncs {
	return cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			e := events.New(events.Added, obj)
			k.Dispatch(e)
		},
		DeleteFunc: func(obj interface{}) {
			e := events.New(events.Deleted, obj)
			k.Dispatch(e)
		},
		UpdateFunc: func(old, new interface{}) {
			if !reflect.DeepEqual(old, new) {
				e := events.New(events.Updated, old, new)
				k.Dispatch(e)
			}
		},
	}
}
