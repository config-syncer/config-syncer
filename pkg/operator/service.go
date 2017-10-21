package operator

import (
	"errors"
	"reflect"

	"github.com/appscode/go/log"
	acrt "github.com/appscode/go/runtime"
	"github.com/appscode/kubed/pkg/util"
	kutil "github.com/appscode/kutil/core/v1"
	pcm "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

func (op *Operator) watchService() {
	if !util.IsPreferredAPIResource(op.KubeClient, core.SchemeGroupVersion.String(), "Service") {
		log.Warningf("Skipping watching non-preferred GroupVersion:%s Kind:%s", core.SchemeGroupVersion.String(), "Service")
		return
	}

	defer acrt.HandleCrash()
	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return op.KubeClient.CoreV1().Services(core.NamespaceAll).List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.KubeClient.CoreV1().Services(core.NamespaceAll).Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&core.Service{},
		op.Opt.ResyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if res, ok := obj.(*core.Service); ok {
					log.Infof("Service %s@%s added", res.Name, res.Namespace)
					kutil.AssignTypeKind(res)

					if op.Config.APIServer.EnableSearchIndex {
						if err := op.SearchIndex.HandleAdd(obj); err != nil {
							log.Errorln(err)
						}
					}
					if op.Config.APIServer.EnableReverseIndex {
						op.ReverseIndex.Service.Add(res)
						if op.ReverseIndex.ServiceMonitor != nil {
							serviceMonitors, err := op.PromClient.ServiceMonitors(core.NamespaceAll).List(metav1.ListOptions{})
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
				if res, ok := obj.(*core.Service); ok {
					log.Infof("Service %s@%s deleted", res.Name, res.Namespace)
					kutil.AssignTypeKind(res)

					if op.Config.APIServer.EnableSearchIndex {
						if err := op.SearchIndex.HandleDelete(obj); err != nil {
							log.Errorln(err)
						}
					}
					if op.Config.APIServer.EnableReverseIndex {
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
				oldRes, ok := old.(*core.Service)
				if !ok {
					log.Errorln(errors.New("Invalid Service object"))
					return
				}
				newRes, ok := new.(*core.Service)
				if !ok {
					log.Errorln(errors.New("Invalid Service object"))
					return
				}
				kutil.AssignTypeKind(oldRes)
				kutil.AssignTypeKind(newRes)

				if op.Config.APIServer.EnableSearchIndex {
					op.SearchIndex.HandleUpdate(old, new)
				}
				if op.Config.APIServer.EnableReverseIndex {
					op.ReverseIndex.Service.Update(oldRes, newRes)
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
