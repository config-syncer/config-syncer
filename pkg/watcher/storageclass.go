package watcher

import (
	acrt "github.com/appscode/go/runtime"
	"github.com/appscode/log"
	"github.com/appscode/stash/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	storage "k8s.io/client-go/pkg/apis/storage/v1"
	"k8s.io/client-go/tools/cache"
)

// Blocks caller. Intended to be called as a Go routine.
func (c *Controller) WatchStorageClasss() {
	if !util.IsPreferredAPIResource(c.KubeClient, extensions.SchemeGroupVersion.String(), "StorageClass") {
		log.Warningf("Skipping watching non-preferred GroupVersion:%s Kind:%s", extensions.SchemeGroupVersion.String(), "StorageClass")
		return
	}

	defer acrt.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return c.KubeClient.StorageV1().StorageClasses().List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return c.KubeClient.StorageV1().StorageClasses().Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&storage.StorageClass{},
		c.SyncPeriod,
		cache.ResourceEventHandlerFuncs{
			DeleteFunc: func(obj interface{}) {
				if storage, ok := obj.(*storage.StorageClass); ok {
					log.Infof("StorageClass %s@%s deleted", storage.Name, storage.Namespace)

				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}
