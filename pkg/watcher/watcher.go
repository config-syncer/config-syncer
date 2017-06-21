package watcher

import (
	"sync"
	"time"

	"github.com/appscode/kubed/pkg/events"
	"github.com/appscode/kubed/pkg/handlers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
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
	lw := cache.NewListWatchFromClient(
		w.KubeClient.CoreV1().RESTClient(),
		events.Namespace.String(),
		metav1.NamespaceAll,
		fields.Everything())
	_, controller := cache.NewInformer(lw, &apiv1.Namespace{}, w.SyncPeriod, eventHandlerFuncs(w))
	go controller.Run(wait.NeverStop)
}

func eventHandlerFuncs(k *Watcher) cache.ResourceEventHandlerFuncs {
	return cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			e := events.New(events.Added, obj)
			if e.Ignorable() {
				return
			}
			if e.ResourceType == events.Namespace && e.EventType == events.Added {
				h := &handlers.NamespaceHandler{
					KubeClient: k.KubeClient,
				}
				h.Handle(e)
			}
		},
	}
}
