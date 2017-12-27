package operator

import (
	"errors"
	"reflect"

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

// Blocks caller. Intended to be called as a Go routine.
func (op *Operator) WatchSecrets() {
	defer rt.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return op.KubeClient.CoreV1().Secrets(core.NamespaceAll).List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.KubeClient.CoreV1().Secrets(core.NamespaceAll).Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&core.Secret{},
		op.Opt.ResyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if res, ok := obj.(*core.Secret); ok {
					log.Infof("Secret %s@%s added", res.Name, res.Namespace)
					kutil.AssignTypeKind(res)

					if op.Config.APIServer.EnableSearchIndex {
						if err := op.SearchIndex.HandleAdd(kutil.ObfuscateSecret(*res)); err != nil {
							log.Errorln(err)
						}
					}
					if op.ConfigSyncer != nil {
						op.ConfigSyncer.SyncSecret(nil, res)
					}
				}
			},
			DeleteFunc: func(obj interface{}) {
				if res, ok := obj.(*core.Secret); ok {
					log.Infof("Secret %s@%s deleted", res.Name, res.Namespace)
					kutil.AssignTypeKind(res)

					if op.Config.APIServer.EnableSearchIndex {
						if err := op.SearchIndex.HandleDelete(kutil.ObfuscateSecret(*res)); err != nil {
							log.Errorln(err)
						}
					}
					if op.TrashCan != nil {
						op.TrashCan.Delete(res.TypeMeta, res.ObjectMeta, kutil.ObfuscateSecret(*res))
					}
					if op.ConfigSyncer != nil {
						op.ConfigSyncer.SyncSecret(res, nil)
					}
				}
			},
			UpdateFunc: func(old, new interface{}) {
				oldRes, ok := old.(*core.Secret)
				if !ok {
					log.Errorln(errors.New("invalid Secret object"))
					return
				}
				newRes, ok := new.(*core.Secret)
				if !ok {
					log.Errorln(errors.New("invalid Secret object"))
					return
				}
				kutil.AssignTypeKind(oldRes)
				kutil.AssignTypeKind(newRes)

				if op.Config.APIServer.EnableSearchIndex {
					op.SearchIndex.HandleUpdate(kutil.ObfuscateSecret(*oldRes), kutil.ObfuscateSecret(*newRes))
				}
				if !reflect.DeepEqual(oldRes.Labels, newRes.Labels) ||
					!reflect.DeepEqual(oldRes.Annotations, newRes.Annotations) ||
					!reflect.DeepEqual(oldRes.Data, newRes.Data) {
					if op.TrashCan != nil && op.Config.RecycleBin.HandleUpdates {
						op.TrashCan.Update(newRes.TypeMeta, newRes.ObjectMeta, kutil.ObfuscateSecret(*oldRes), kutil.ObfuscateSecret(*newRes))
					}

					if op.ConfigSyncer != nil {
						op.ConfigSyncer.SyncSecret(oldRes, newRes)
					}
				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}
