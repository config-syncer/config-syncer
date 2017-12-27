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
func (op *Operator) WatchConfigMaps() {
	defer rt.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return op.KubeClient.CoreV1().ConfigMaps(core.NamespaceAll).List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.KubeClient.CoreV1().ConfigMaps(core.NamespaceAll).Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&core.ConfigMap{},
		op.Opt.ResyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if res, ok := obj.(*core.ConfigMap); ok {
					log.Infof("ConfigMap %s@%s added", res.Name, res.Namespace)
					kutil.AssignTypeKind(res)

					if op.Config.APIServer.EnableSearchIndex {
						if err := op.SearchIndex.HandleAdd(obj); err != nil {
							log.Errorln(err)
						}
					}
					if op.ConfigSyncer != nil {
						op.ConfigSyncer.SyncConfigMap(nil, res)
					}
				}
			},
			DeleteFunc: func(obj interface{}) {
				if res, ok := obj.(*core.ConfigMap); ok {
					log.Infof("ConfigMap %s@%s deleted", res.Name, res.Namespace)
					kutil.AssignTypeKind(res)

					if op.Config.APIServer.EnableSearchIndex {
						if err := op.SearchIndex.HandleDelete(obj); err != nil {
							log.Errorln(err)
						}
					}
					if op.TrashCan != nil {
						op.TrashCan.Delete(res.TypeMeta, res.ObjectMeta, obj)
					}
					if op.ConfigSyncer != nil {
						op.ConfigSyncer.SyncConfigMap(res, nil)
					}
				}
			},
			UpdateFunc: func(old, new interface{}) {
				oldRes, ok := old.(*core.ConfigMap)
				if !ok {
					log.Errorln(errors.New("invalid ConfigMap object"))
					return
				}
				newRes, ok := new.(*core.ConfigMap)
				if !ok {
					log.Errorln(errors.New("invalid ConfigMap object"))
					return
				}
				kutil.AssignTypeKind(oldRes)
				kutil.AssignTypeKind(newRes)

				if op.Config.APIServer.EnableSearchIndex {
					op.SearchIndex.HandleUpdate(old, new)
				}
				if !reflect.DeepEqual(oldRes.Labels, newRes.Labels) ||
					!reflect.DeepEqual(oldRes.Annotations, newRes.Annotations) ||
					!reflect.DeepEqual(oldRes.Data, newRes.Data) {
					if op.TrashCan != nil && op.Config.RecycleBin.HandleUpdates {
						op.TrashCan.Update(newRes.TypeMeta, newRes.ObjectMeta, old, new)
					}

					if op.ConfigSyncer != nil {
						op.ConfigSyncer.SyncConfigMap(oldRes, newRes)
					}
				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}
