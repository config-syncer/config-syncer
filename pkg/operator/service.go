package operator

import (
	acrt "github.com/appscode/go/runtime"
	"github.com/appscode/kubed/pkg/util"
	"github.com/appscode/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
)

func (op *Operator) watchService() {
	if !util.IsPreferredAPIResource(op.KubeClient, apiv1.SchemeGroupVersion.String(), "Service") {
		log.Warningf("Skipping watching non-preferred GroupVersion:%s Kind:%s", apiv1.SchemeGroupVersion.String(), "Service")
		return
	}

	defer acrt.HandleCrash()
	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return op.KubeClient.CoreV1().ConfigMaps(apiv1.NamespaceAll).List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.KubeClient.CoreV1().ConfigMaps(apiv1.NamespaceAll).Watch(metav1.ListOptions{})
		},
	}
	_, controller := cache.NewInformer(lw,
		&apiv1.Service{},
		op.syncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if op.Opt.EnableReverseIndex {
					op.ReverseIndex.Handle("added", obj)
				}
				if op.Opt.EnableSearchIndex {
					op.SearchIndex.HandleAdd(obj)
				}
			},
			DeleteFunc: func(obj interface{}) {
				if op.Opt.EnableReverseIndex {
					op.ReverseIndex.Handle("deleted", obj)
				}
				if op.Opt.EnableSearchIndex {
					op.SearchIndex.HandleDelete(obj)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				if op.Opt.EnableReverseIndex {
					op.ReverseIndex.Handle("updated", oldObj, newObj)
				}
				if op.Opt.EnableSearchIndex {
					op.SearchIndex.HandleUpdate(oldObj, newObj)
				}
			},
		},
	)
	go controller.Run(wait.NeverStop)
}
