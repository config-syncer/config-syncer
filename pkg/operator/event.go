package operator

import (
	"github.com/appscode/go/log"
	acrt "github.com/appscode/go/runtime"
	"github.com/appscode/kubed/pkg/util"
	kutil "github.com/appscode/kutil/core/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

// Blocks caller. Intended to be called as a Go routine.
func (op *Operator) WatchEvents() {
	if !util.IsPreferredAPIResource(op.KubeClient, core.SchemeGroupVersion.String(), "Event") {
		log.Warningf("Skipping watching non-preferred GroupVersion:%s Kind:%s", core.SchemeGroupVersion.String(), "Event")
		return
	}

	defer acrt.HandleCrash()

	fs := fields.OneTermEqualSelector("type", core.EventTypeWarning).String()
	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return op.KubeClient.CoreV1().Events(core.NamespaceAll).List(metav1.ListOptions{
				FieldSelector: fs,
			})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.KubeClient.CoreV1().Events(core.NamespaceAll).Watch(metav1.ListOptions{
				FieldSelector: fs,
			})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&core.Event{},
		op.Opt.ResyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if res, ok := obj.(*core.Event); ok {
					log.Infof("Event %s@%s added", res.Name, res.Namespace)
					kutil.AssignTypeKind(res)

					if op.Eventer != nil &&
						op.Config.EventForwarder.WarningEvents.Handle &&
						op.Config.EventForwarder.WarningEvents.IsAllowed(res.Namespace) &&
						util.IsRecent(res.ObjectMeta.CreationTimestamp) {
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
