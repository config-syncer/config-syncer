package controller

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

func (c *Controller) watchService() {
	if !util.IsPreferredAPIResource(c.KubeClient, apiv1.SchemeGroupVersion.String(), "Service") {
		log.Warningf("Skipping watching non-preferred GroupVersion:%s Kind:%s", apiv1.SchemeGroupVersion.String(), "Service")
		return
	}

	defer acrt.HandleCrash()
	lw := cache.NewListWatchFromClient(
		c.KubeClient.CoreV1().RESTClient(),
		"services",
		metav1.NamespaceAll,
		fields.Everything())
	_, controller := cache.NewInformer(lw,
		&apiv1.Service{},
		c.SyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if c.Opt.EnableReverseIndex {
					c.ReverseIndex.Handle("added", obj)
				}
				if len(c.Opt.Indexer) > 0 {
					c.Indexer.HandleAdd(obj)
				}
			},
			DeleteFunc: func(obj interface{}) {
				if c.Opt.EnableReverseIndex {
					c.ReverseIndex.Handle("deleted", obj)
				}
				if len(c.Opt.Indexer) > 0 {
					c.Indexer.HandleDelete(obj)
				}
				c.Saver.Save(obj.(*apiv1.Service).ObjectMeta, obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				if c.Opt.EnableReverseIndex {
					c.ReverseIndex.Handle("updated", oldObj, newObj)
				}
				if len(c.Opt.Indexer) > 0 {
					c.Indexer.HandleUpdate(oldObj, newObj)
				}
			},
		},
	)
	go controller.Run(wait.NeverStop)
}
