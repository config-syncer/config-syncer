package operator

import (
	acrt "github.com/appscode/go/runtime"
	"github.com/appscode/kubed/pkg/util"
	"github.com/appscode/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	storage "k8s.io/client-go/pkg/apis/storage/v1"
	"k8s.io/client-go/tools/cache"
)

// Blocks caller. Intended to be called as a Go routine.
func (op *Operator) WatchStorageClasss() {
	if !util.IsPreferredAPIResource(op.KubeClient, storage.SchemeGroupVersion.String(), "StorageClass") {
		log.Warningf("Skipping watching non-preferred GroupVersion:%s Kind:%s", storage.SchemeGroupVersion.String(), "StorageClass")
		return
	}

	defer acrt.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return op.KubeClient.StorageV1().StorageClasses().List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.KubeClient.StorageV1().StorageClasses().Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&storage.StorageClass{},
		op.syncPeriod,
		cache.ResourceEventHandlerFuncs{
			DeleteFunc: func(obj interface{}) {
				if sc, ok := obj.(*storage.StorageClass); ok {
					log.Infof("StorageClass %s@%s deleted", sc.Name, sc.Namespace)
					op.TrashCan.Delete(sc.ObjectMeta, obj)
				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}
