package operator

import (
	"github.com/appscode/go/log"
	acrt "github.com/appscode/go/runtime"
	"github.com/appscode/kubed/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/pkg/api"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
)

// Blocks caller. Intended to be called as a Go routine.
func (op *Operator) WatchEvents() {
	if !util.IsPreferredAPIResource(op.KubeClient, apiv1.SchemeGroupVersion.String(), "Event") {
		log.Warningf("Skipping watching non-preferred GroupVersion:%s Kind:%s", apiv1.SchemeGroupVersion.String(), "Event")
		return
	}

	defer acrt.HandleCrash()

	fs := fields.OneTermEqualSelector(api.EventTypeField, apiv1.EventTypeWarning).String()
	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return op.KubeClient.CoreV1().Events(apiv1.NamespaceAll).List(metav1.ListOptions{
				FieldSelector: fs,
			})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.KubeClient.CoreV1().Events(apiv1.NamespaceAll).Watch(metav1.ListOptions{
				FieldSelector: fs,
			})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&apiv1.Event{},
		op.syncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if res, ok := obj.(*apiv1.Event); ok {
					log.Infof("Event %s@%s added", res.Name, res.Namespace)
					if op.Eventer != nil &&
						op.Config.EventForwarder.WarningEvents.Handle &&
						op.Config.EventForwarder.WarningEvents.IsAllowed(res.Namespace) &&
						util.IsRecentlyAdded(res.ObjectMeta) {
						err := op.Eventer.ForwardEvent(res)
						if err != nil {
							log.Errorln(err)
						}
					}
				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}
