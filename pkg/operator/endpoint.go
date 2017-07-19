package operator

import (
	"reflect"

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

// Blocks caller. Intended to be called as a Go routine.
func (op *Operator) WatchEndpoints() {
	if !util.IsPreferredAPIResource(op.KubeClient, apiv1.SchemeGroupVersion.String(), "Endpoints") {
		log.Warningf("Skipping watching non-preferred GroupVersion:%s Kind:%s", apiv1.SchemeGroupVersion.String(), "Endpoints")
		return
	}

	defer acrt.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return op.KubeClient.CoreV1().Endpoints(apiv1.NamespaceAll).List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.KubeClient.CoreV1().Endpoints(apiv1.NamespaceAll).Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&apiv1.Endpoints{},
		op.syncPeriod,
		cache.ResourceEventHandlerFuncs{
			UpdateFunc: func(oldObj, newObj interface{}) {
				if oldRes, ok := oldObj.(*apiv1.Endpoints); ok {
					if newRes, ok := newObj.(*apiv1.Endpoints); ok {
						if !reflect.DeepEqual(oldRes.Subsets, newRes.Subsets) && op.ReverseIndex.Service != nil {
							svc, err := op.KubeClient.CoreV1().Services(newRes.Namespace).Get(newRes.Name, metav1.GetOptions{})
							if err == nil && svc != nil {
								op.ReverseIndex.Service.Add(svc)
							}
						}
					}
				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}
