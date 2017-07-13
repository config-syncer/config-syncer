package watcher

import (
	acrt "github.com/appscode/go/runtime"
	"github.com/appscode/kubed/pkg/util"
	"github.com/appscode/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
)

func (w *Watchers) watchService() {
	if !util.IsPreferredAPIResource(w.KubeClient, apiv1.SchemeGroupVersion.String(), "Service") {
		log.Warningf("Skipping watching non-preferred GroupVersion:%s Kind:%s", apiv1.SchemeGroupVersion.String(), "Service")
		return
	}

	defer acrt.HandleCrash()
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
				if w.Opt.EnableReverseIndex {
					w.ReverseIndex.Handle("added", obj)
				}
				if len(w.Opt.Indexer) > 0 {
					w.Indexer.HandleAdd(obj)
				}
			},
			DeleteFunc: func(obj interface{}) {
				if w.Opt.EnableReverseIndex {
					w.ReverseIndex.Handle("deleted", obj)
				}
				if len(w.Opt.Indexer) > 0 {
					w.Indexer.HandleDelete(obj)
				}
				w.Saver.Save(obj.(*apiv1.Service).ObjectMeta, obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				if w.Opt.EnableReverseIndex {
					w.ReverseIndex.Handle("updated", oldObj, newObj)
				}
				if len(w.Opt.Indexer) > 0 {
					w.Indexer.HandleUpdate(oldObj, newObj)
				}
			},
		},
	)
	go controller.Run(wait.NeverStop)
}
