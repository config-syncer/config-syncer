package watcher

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
)

func (w *Controller) watchService() {
	lw := cache.NewListWatchFromClient(
		w.KubeClient.CoreV1().RESTClient(),
		"services",
		metav1.NamespaceAll,
		fields.Everything())
	_, controller := cache.NewInformer(lw,
		&apiv1.Service{},
		w.SyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				w.ReverseIndex.Handle("added", obj)
				w.Indexer.HandleAdd(obj)
			},
			DeleteFunc: func(obj interface{}) {
				w.ReverseIndex.Handle("deleted", obj)
				w.Indexer.HandleDelete(obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				w.ReverseIndex.Handle("updated", oldObj, newObj)
				w.Indexer.HandleUpdate(oldObj, newObj)
			},
		},
	)
	go controller.Run(wait.NeverStop)
}
