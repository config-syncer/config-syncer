package operator

import (
	"errors"
	"reflect"

	"github.com/appscode/go/log"
	acrt "github.com/appscode/go/runtime"
	kutil "github.com/appscode/kube-mon/prometheus/v1"
	"github.com/appscode/kubed/pkg/util"
	prom "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

// Blocks caller. Intended to be called as a Go routine.
func (op *Operator) WatchPrometheus() {
	if !util.IsPreferredAPIResource(op.KubeClient, prom.Group+"/"+prom.Version, prom.PrometheusesKind) {
		log.Warningf("Skipping watching non-preferred GroupVersion:%s Kind:%s", prom.Group+"/"+prom.Version, prom.PrometheusesKind)
		return
	}

	defer acrt.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return op.PromClient.Prometheuses(core.NamespaceAll).List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.PromClient.Prometheuses(core.NamespaceAll).Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&prom.Prometheus{},
		op.Opt.ResyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if res, ok := obj.(*prom.Prometheus); ok {
					log.Infof("Prometheus %s@%s added", res.Name, res.Namespace)
					kutil.AssignTypeKind(res)

					if op.Config.APIServer.EnableSearchIndex {
						if err := op.SearchIndex.HandleAdd(obj); err != nil {
							log.Errorln(err)
						}
					}

					if op.Config.APIServer.EnableReverseIndex {
						if err := op.ReverseIndex.Prometheus.Add(res); err != nil {
							log.Errorln(err)
						}
					}
				}
			},
			DeleteFunc: func(obj interface{}) {
				if res, ok := obj.(*prom.Prometheus); ok {
					log.Infof("Prometheus %s@%s deleted", res.Name, res.Namespace)
					kutil.AssignTypeKind(res)

					if op.Config.APIServer.EnableSearchIndex {
						if err := op.SearchIndex.HandleDelete(obj); err != nil {
							log.Errorln(err)
						}
					}

					if op.Config.APIServer.EnableReverseIndex {
						if err := op.ReverseIndex.Prometheus.Delete(res); err != nil {
							log.Errorln(err)
						}
					}

					if op.TrashCan != nil {
						op.TrashCan.Delete(res.TypeMeta, res.ObjectMeta, obj)
					}
				}
			},
			UpdateFunc: func(old, new interface{}) {
				oldRes, ok := old.(*prom.Prometheus)
				if !ok {
					log.Errorln(errors.New("invalid Prometheus object"))
					return
				}
				newRes, ok := new.(*prom.Prometheus)
				if !ok {
					log.Errorln(errors.New("invalid Prometheus object"))
					return
				}
				kutil.AssignTypeKind(oldRes)
				kutil.AssignTypeKind(newRes)

				if op.Config.APIServer.EnableSearchIndex {
					op.SearchIndex.HandleUpdate(old, new)
				}

				if op.Config.APIServer.EnableReverseIndex {
					if err := op.ReverseIndex.Prometheus.Update(oldRes, newRes); err != nil {
						log.Errorln(err)
					}
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
