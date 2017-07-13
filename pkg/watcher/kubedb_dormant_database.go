package watcher

import (
	"errors"
	"fmt"

	acrt "github.com/appscode/go/runtime"
	"github.com/appscode/kubed/pkg/util"
	"github.com/appscode/log"
	tapi "github.com/k8sdb/apimachinery/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
)

// Blocks caller. Intended to be called as a Go routine.
func (c *Controller) WatchDormantDatabases() {
	if !util.IsPreferredAPIResource(c.KubeClient, tapi.V1alpha1SchemeGroupVersion.String(), tapi.ResourceKindDormantDatabase) {
		log.Warningf("Skipping watching non-preferred GroupVersion:%s Kind:%s", tapi.V1alpha1SchemeGroupVersion.String(), tapi.ResourceKindDormantDatabase)
		return
	}

	defer acrt.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return c.KubeDBClient.DormantDatabases(apiv1.NamespaceAll).List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return c.KubeDBClient.DormantDatabases(apiv1.NamespaceAll).Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&tapi.DormantDatabase{},
		c.SyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if drmn, ok := obj.(*tapi.DormantDatabase); ok {
					fmt.Println(drmn)
				}
			},
			UpdateFunc: func(old, new interface{}) {
				oldAlert, ok := old.(*tapi.DormantDatabase)
				if !ok {
					log.Errorln(errors.New("Invalid DormantDatabase object"))
					return
				}
				newAlert, ok := new.(*tapi.DormantDatabase)
				if !ok {
					log.Errorln(errors.New("Invalid DormantDatabase object"))
					return
				}
				fmt.Println(oldAlert, newAlert)
			},
			DeleteFunc: func(obj interface{}) {
				if drmn, ok := obj.(*tapi.DormantDatabase); ok {
					fmt.Println(drmn)
					c.Saver.Save(drmn.ObjectMeta, obj)
				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}
