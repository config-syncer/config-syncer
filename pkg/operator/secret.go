package operator

import (
	"errors"
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
func (op *Operator) WatchSecrets() {
	if !util.IsPreferredAPIResource(op.KubeClient, apiv1.SchemeGroupVersion.String(), "Secret") {
		log.Warningf("Skipping watching non-preferred GroupVersion:%s Kind:%s", apiv1.SchemeGroupVersion.String(), "Secret")
		return
	}

	defer acrt.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return op.KubeClient.CoreV1().Secrets(apiv1.NamespaceAll).List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.KubeClient.CoreV1().Secrets(apiv1.NamespaceAll).Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&apiv1.Secret{},
		op.syncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if res, ok := obj.(*apiv1.Secret); ok {
					log.Infof("Secret %s@%s added", res.Name, res.Namespace)

					if op.Opt.EnableSearchIndex {
						if err := op.SearchIndex.HandleAdd(obj); err != nil {
							log.Errorln(err)
						}
					}
					if op.ConfigSyncer != nil {
						op.ConfigSyncer.SyncSecret(nil, res)
					}
				}
			},
			DeleteFunc: func(obj interface{}) {
				if res, ok := obj.(*apiv1.Secret); ok {
					log.Infof("Secret %s@%s deleted", res.Name, res.Namespace)
					if op.Opt.EnableSearchIndex {
						if err := op.SearchIndex.HandleDelete(obj); err != nil {
							log.Errorln(err)
						}
					}
					if op.TrashCan != nil {
						op.TrashCan.Delete(res.TypeMeta, res.ObjectMeta, obj)
					}
					if op.ConfigSyncer != nil {
						op.ConfigSyncer.SyncSecret(res, nil)
					}
				}
			},
			UpdateFunc: func(old, new interface{}) {
				oldRes, ok := old.(*apiv1.Secret)
				if !ok {
					log.Errorln(errors.New("Invalid Secret object"))
					return
				}
				newRes, ok := new.(*apiv1.Secret)
				if !ok {
					log.Errorln(errors.New("Invalid Secret object"))
					return
				}
				if op.Opt.EnableSearchIndex {
					op.SearchIndex.HandleUpdate(old, new)
				}
				if op.TrashCan != nil && op.Config.TrashCan.HandleUpdate {
					if !reflect.DeepEqual(oldRes.Labels, newRes.Labels) ||
						!reflect.DeepEqual(oldRes.Annotations, newRes.Annotations) ||
						!reflect.DeepEqual(oldRes.Data, newRes.Data) {
						op.TrashCan.Update(newRes.TypeMeta, newRes.ObjectMeta, old, new)
					}
				}
				if op.ConfigSyncer != nil {
					op.ConfigSyncer.SyncSecret(oldRes, newRes)
				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}
