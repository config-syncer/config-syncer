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
			},
			DeleteFunc: func(obj interface{}) {
				w.ReverseIndex.Handle("deleted", obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				w.ReverseIndex.Handle("updated", oldObj, newObj)
			},
		},
	)
	go controller.Run(wait.NeverStop)
}
