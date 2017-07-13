package operator

import (
	"errors"
	"fmt"

	acrt "github.com/appscode/go/runtime"
	"github.com/appscode/kubed/pkg/util"
	"github.com/appscode/log"
	tapi "github.com/appscode/stash/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
)

// Blocks caller. Intended to be called as a Go routine.
func (op *Operator) WatchRestics() {
	if !util.IsPreferredAPIResource(op.KubeClient, tapi.V1alpha1SchemeGroupVersion.String(), tapi.ResourceKindRestic) {
		log.Warningf("Skipping watching non-preferred GroupVersion:%s Kind:%s", tapi.V1alpha1SchemeGroupVersion.String(), tapi.ResourceKindRestic)
		return
	}
	defer acrt.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return op.StashClient.Restics(apiv1.NamespaceAll).List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.StashClient.Restics(apiv1.NamespaceAll).Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&tapi.Restic{},
		op.SyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if restic, ok := obj.(*tapi.Restic); ok {
					fmt.Println(restic)
				}
			},
			UpdateFunc: func(old, new interface{}) {
				oldRestic, ok := old.(*tapi.Restic)
				if !ok {
					log.Errorln(errors.New("Invalid Restic object"))
					return
				}
				newRestic, ok := new.(*tapi.Restic)
				if !ok {
					log.Errorln(errors.New("Invalid Restic object"))
					return
				}
				fmt.Println(oldRestic, newRestic)
			},
			DeleteFunc: func(obj interface{}) {
				if restic, ok := obj.(*tapi.Restic); ok {
					fmt.Println(restic)
					op.Saver.Save(restic.ObjectMeta, obj)
				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}
