package operator

import (
	"errors"
	"reflect"

	"github.com/appscode/go/log"
	acrt "github.com/appscode/go/runtime"
	"github.com/appscode/kubed/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
)

// Blocks caller. Intended to be called as a Go routine.
func (op *Operator) WatchNodes() {
	if !util.IsPreferredAPIResource(op.KubeClient, apiv1.SchemeGroupVersion.String(), "Node") {
		log.Warningf("Skipping watching non-preferred GroupVersion:%s Kind:%s", apiv1.SchemeGroupVersion.String(), "Node")
		return
	}

	defer acrt.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return op.KubeClient.CoreV1().Nodes().List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.KubeClient.CoreV1().Nodes().Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&apiv1.Node{},
		op.syncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if res, ok := obj.(*apiv1.Node); ok {
					log.Infof("Node %s@%s added", res.Name, res.Namespace)
					util.AssignTypeKind(res)

					if op.Config.APIServer.EnableSearchIndex {
						if err := op.SearchIndex.HandleAdd(obj); err != nil {
							log.Errorln(err)
						}
					}

					if op.Eventer != nil &&
						op.Config.EventForwarder.NodeAdded.Handle &&
						op.Config.EventForwarder.NodeAdded.IsAllowed(res.Namespace) &&
						util.IsRecentlyAdded(res.ObjectMeta) {
						err := op.Eventer.Forward(res.TypeMeta, res.ObjectMeta, obj)
						if err != nil {
							log.Errorln(err)
						}
					}
				}
			},
			DeleteFunc: func(obj interface{}) {
				if res, ok := obj.(*apiv1.Node); ok {
					log.Infof("Node %s@%s deleted", res.Name, res.Namespace)
					util.AssignTypeKind(res)

					if op.Config.APIServer.EnableSearchIndex {
						if err := op.SearchIndex.HandleDelete(obj); err != nil {
							log.Errorln(err)
						}
					}
					if op.TrashCan != nil {
						op.TrashCan.Delete(res.TypeMeta, res.ObjectMeta, obj)
					}
				}
			},
			UpdateFunc: func(old, new interface{}) {
				oldRes, ok := old.(*apiv1.Node)
				if !ok {
					log.Errorln(errors.New("Invalid Node object"))
					return
				}
				newRes, ok := new.(*apiv1.Node)
				if !ok {
					log.Errorln(errors.New("Invalid Node object"))
					return
				}
				util.AssignTypeKind(oldRes)
				util.AssignTypeKind(newRes)

				if op.Config.APIServer.EnableSearchIndex {
					op.SearchIndex.HandleUpdate(old, new)
				}
				if op.TrashCan != nil && op.Config.RecycleBin.HandleUpdates {
					if !reflect.DeepEqual(oldRes.Labels, newRes.Labels) ||
						!reflect.DeepEqual(oldRes.Annotations, newRes.Annotations) ||
						!reflect.DeepEqual(oldRes.Spec, newRes.Spec) {
						op.TrashCan.Update(newRes.TypeMeta, newRes.ObjectMeta, old, new)
					}
				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}
