package watcher

import (
	acrt "github.com/appscode/go/runtime"
	"github.com/appscode/kubed/pkg/util"
	"github.com/appscode/log"
	prom "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
)

// Blocks caller. Intended to be called as a Go routine.
func (w *Watchers) WatchAlertmanagers() {
	if !util.IsPreferredAPIResource(w.KubeClient, prom.TPRGroup+"/"+prom.TPRVersion, prom.TPRAlertmanagersKind) {
		log.Warningf("Skipping watching non-preferred GroupVersion:%s Kind:%s", prom.TPRGroup+"/"+prom.TPRVersion, prom.TPRAlertmanagersKind)
		return
	}

	defer acrt.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return w.PromClient.Alertmanagers(apiv1.NamespaceAll).List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return w.PromClient.Alertmanagers(apiv1.NamespaceAll).Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&prom.Alertmanager{},
		w.SyncPeriod,
		cache.ResourceEventHandlerFuncs{
			DeleteFunc: func(obj interface{}) {
				if mgr, ok := obj.(*prom.Alertmanager); ok {
					log.Infof("Alertmanager %s@%s deleted", mgr.Name, mgr.Namespace)
					w.Saver.Save(mgr.ObjectMeta, obj)
				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}
