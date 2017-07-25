package operator

import (
	"errors"
	"reflect"

	acrt "github.com/appscode/go/runtime"
	"github.com/appscode/kubed/pkg/util"
	"github.com/appscode/log"
	pcm "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
)

func (op *Operator) watchService() {
	if !util.IsPreferredAPIResource(op.KubeClient, apiv1.SchemeGroupVersion.String(), "Service") {
		log.Warningf("Skipping watching non-preferred GroupVersion:%s Kind:%s", apiv1.SchemeGroupVersion.String(), "Service")
		return
	}

	defer acrt.HandleCrash()
	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return op.KubeClient.CoreV1().Services(apiv1.NamespaceAll).List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.KubeClient.CoreV1().Services(apiv1.NamespaceAll).Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&apiv1.Service{},
		op.syncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if res, ok := obj.(*apiv1.Service); ok {
					log.Infof("Service %s@%s added", res.Name, res.Namespace)
					util.AssignTypeKind(res)

					if op.Opt.EnableSearchIndex {
						if err := op.SearchIndex.HandleAdd(obj); err != nil {
							log.Errorln(err)
						}
					}
					if op.Opt.EnableReverseIndex {
						op.ReverseIndex.Service.Add(res)
						if op.ReverseIndex.ServiceMonitor != nil {
							serviceMonitors, err := op.PromClient.ServiceMonitors(apiv1.NamespaceAll).List(metav1.ListOptions{})
							if err != nil {
								log.Errorln(err)
								return
							}
							if serviceMonitorList, ok := serviceMonitors.(*pcm.ServiceMonitorList); ok {
								op.ReverseIndex.ServiceMonitor.AddService(res, serviceMonitorList.Items)
							}
						}
					}
				}
			},
			DeleteFunc: func(obj interface{}) {
				if res, ok := obj.(*apiv1.Service); ok {
					log.Infof("Service %s@%s deleted", res.Name, res.Namespace)
					util.AssignTypeKind(res)

					if op.Opt.EnableSearchIndex {
						if err := op.SearchIndex.HandleDelete(obj); err != nil {
							log.Errorln(err)
						}
					}
					if op.Opt.EnableReverseIndex {
						op.ReverseIndex.Service.Delete(res)
						if op.ReverseIndex.ServiceMonitor != nil {
							op.ReverseIndex.ServiceMonitor.DeleteService(res)
						}
					}
					if op.TrashCan != nil {
						op.TrashCan.Delete(res.TypeMeta, res.ObjectMeta, obj)
					}
				}
			},
			UpdateFunc: func(old, new interface{}) {
				oldRes, ok := old.(*apiv1.Service)
				if !ok {
					log.Errorln(errors.New("Invalid Service object"))
					return
				}
				newRes, ok := new.(*apiv1.Service)
				if !ok {
					log.Errorln(errors.New("Invalid Service object"))
					return
				}
				util.AssignTypeKind(oldRes)
				util.AssignTypeKind(newRes)

				if op.Opt.EnableSearchIndex {
					op.SearchIndex.HandleUpdate(old, new)
				}
				if op.Opt.EnableReverseIndex {
					op.ReverseIndex.Service.Update(oldRes, newRes)
				}
				if op.TrashCan != nil && op.Config.RecycleBin.HandleUpdate {
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
