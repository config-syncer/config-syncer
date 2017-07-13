package watcher

import (
	"github.com/appscode/kubed/pkg/configsync"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
)

func (w *Watchers) watchNamespaces() {
	_, controller := cache.NewInformer(
		cache.NewListWatchFromClient(
			w.KubeClient.CoreV1().RESTClient(),
			"namespaces",
			metav1.NamespaceAll,
			fields.Everything()),
		&apiv1.Namespace{},
		w.SyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				ns := configsync.NewHandler(w.KubeClient)
				ns.Handle(obj)
			},
		},
	)
	go controller.Run(wait.NeverStop)
}
