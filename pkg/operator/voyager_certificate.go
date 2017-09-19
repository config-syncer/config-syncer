package operator

import (
	"errors"
	"reflect"

	"github.com/appscode/go/log"
	acrt "github.com/appscode/go/runtime"
	"github.com/appscode/kubed/pkg/util"
	kutil "github.com/appscode/kutil/voyager/v1beta1"
	tapi "github.com/appscode/voyager/apis/voyager/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
)

// Blocks caller. Intended to be called as a Go routine.
func (op *Operator) WatchVoyagerCertificates() {
	if !util.IsPreferredAPIResource(op.KubeClient, tapi.SchemeGroupVersion.String(), tapi.ResourceKindCertificate) {
		log.Warningf("Skipping watching non-preferred GroupVersion:%s Kind:%s", tapi.SchemeGroupVersion.String(), tapi.ResourceKindCertificate)
		return
	}
	defer acrt.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return op.VoyagerClient.Certificates(apiv1.NamespaceAll).List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.VoyagerClient.Certificates(apiv1.NamespaceAll).Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&tapi.Certificate{},
		op.Opt.ResyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if res, ok := obj.(*tapi.Certificate); ok {
					log.Infof("Certificate %s@%s added", res.Name, res.Namespace)
					kutil.AssignTypeKind(res)

					if op.Config.APIServer.EnableSearchIndex {
						if err := op.SearchIndex.HandleAdd(obj); err != nil {
							log.Errorln(err)
						}
					}
				}
			},
			DeleteFunc: func(obj interface{}) {
				if res, ok := obj.(*tapi.Certificate); ok {
					log.Infof("Certificate %s@%s deleted", res.Name, res.Namespace)
					kutil.AssignTypeKind(res)

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
				oldRes, ok := old.(*tapi.Certificate)
				if !ok {
					log.Errorln(errors.New("Invalid Certificate object"))
					return
				}
				newRes, ok := new.(*tapi.Certificate)
				if !ok {
					log.Errorln(errors.New("Invalid Certificate object"))
					return
				}
				kutil.AssignTypeKind(oldRes)
				kutil.AssignTypeKind(newRes)

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
