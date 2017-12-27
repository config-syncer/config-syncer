package operator

import (
	"github.com/appscode/go/log"
	kutil "github.com/appscode/kutil/core/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	rt "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

func (op *Operator) watchNamespaces() {
	defer rt.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return op.KubeClient.CoreV1().Namespaces().List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.KubeClient.CoreV1().Namespaces().Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&core.Namespace{},
		op.Opt.ResyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if res, ok := obj.(*core.Namespace); ok {
					log.Infof("Namespace %s added", res.Name)
					kutil.AssignTypeKind(res)

					if op.ConfigSyncer != nil {
						op.ConfigSyncer.SyncIntoNamespace(res.Name)
					}
				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}
