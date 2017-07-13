package watcher

import (
	"errors"
	"fmt"

	acrt "github.com/appscode/go/runtime"
	"github.com/appscode/kubed/pkg/util"
	"github.com/appscode/log"
	tapi "github.com/appscode/searchlight/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
)

// Blocks caller. Intended to be called as a Go routine.
func (w *Watchers) WatchPodAlerts() {
	if !util.IsPreferredAPIResource(w.KubeClient, tapi.V1alpha1SchemeGroupVersion.String(), tapi.ResourceKindPodAlert) {
		log.Warningf("Skipping watching non-preferred GroupVersion:%s Kind:%s", tapi.V1alpha1SchemeGroupVersion.String(), tapi.ResourceKindPodAlert)
		return
	}

	defer acrt.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return w.SearchlightClient.PodAlerts(apiv1.NamespaceAll).List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return w.SearchlightClient.PodAlerts(apiv1.NamespaceAll).Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&tapi.PodAlert{},
		w.SyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if alert, ok := obj.(*tapi.PodAlert); ok {
					fmt.Println(alert)
				}
			},
			UpdateFunc: func(old, new interface{}) {
				oldAlert, ok := old.(*tapi.PodAlert)
				if !ok {
					log.Errorln(errors.New("Invalid PodAlert object"))
					return
				}
				newAlert, ok := new.(*tapi.PodAlert)
				if !ok {
					log.Errorln(errors.New("Invalid PodAlert object"))
					return
				}
				fmt.Println(oldAlert, newAlert)
			},
			DeleteFunc: func(obj interface{}) {
				if alert, ok := obj.(*tapi.PodAlert); ok {
					fmt.Println(alert)
					w.Saver.Save(alert.ObjectMeta, obj)
				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}
