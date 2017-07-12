package watcher

import (
	"errors"
	"reflect"

	acrt "github.com/appscode/go/runtime"
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
func (c *Controller) WatchClusterAlerts() {
	defer acrt.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return c.SearchlightClient.ClusterAlerts(apiv1.NamespaceAll).List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return c.SearchlightClient.ClusterAlerts(apiv1.NamespaceAll).Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&tapi.ClusterAlert{},
		c.SyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if alert, ok := obj.(*tapi.ClusterAlert); ok {
				}
			},
			UpdateFunc: func(old, new interface{}) {
				oldAlert, ok := old.(*tapi.ClusterAlert)
				if !ok {
					log.Errorln(errors.New("Invalid ClusterAlert object"))
					return
				}
				newAlert, ok := new.(*tapi.ClusterAlert)
				if !ok {
					log.Errorln(errors.New("Invalid ClusterAlert object"))
					return
				}
				if !reflect.DeepEqual(oldAlert.Spec, newAlert.Spec) {
				}
			},
			DeleteFunc: func(obj interface{}) {
				if alert, ok := obj.(*tapi.ClusterAlert); ok {

				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}
