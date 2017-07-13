package operator

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

func (op *Operator) watchService() {
	if !util.IsPreferredAPIResource(op.KubeClient, apiv1.SchemeGroupVersion.String(), "Service") {
		log.Warningf("Skipping watching non-preferred GroupVersion:%s Kind:%s", apiv1.SchemeGroupVersion.String(), "Service")
		return
	}

	defer acrt.HandleCrash()
	lw := cache.NewListWatchFromClient(
		op.KubeClient.CoreV1().RESTClient(),
		"services",
		metav1.NamespaceAll,
		fields.Everything())
	_, controller := cache.NewInformer(lw,
		&apiv1.Service{},
		op.SyncPeriod,
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
				op.Saver.Save(obj.(*apiv1.Service).ObjectMeta, obj)
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
