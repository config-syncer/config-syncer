package indexers

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"reflect"
	"sync"

	"github.com/appscode/go/arrays"
	"github.com/appscode/go/log"
	kutil "github.com/appscode/kube-mon/prometheus/v1"
	"github.com/appscode/pat"
	"github.com/blevesearch/bleve"
	prom "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kubernetes/pkg/apis/core"
)

type PrometheusIndexer interface {
	Configure(enable bool)
	PrometheusHandler() cache.ResourceEventHandler
	ServiceMonitorHandler() cache.ResourceEventHandler
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

var _ PrometheusIndexer = &PrometheusIndexerImpl{}

type PrometheusIndexerImpl struct {
	promIndexer cache.Indexer
	smonIndexer cache.Indexer
	index       bleve.Index

	enable bool
	locker sync.RWMutex
}

func NewPrometheusIndexer(dir string, promIndexer cache.Indexer, smonIndexer cache.Indexer) (PrometheusIndexer, error) {
	index, err := ensureIndex(filepath.Join(dir, "prometheus.indexer"), "indexer")
	if err != nil {
		return nil, err
	}
	return &PrometheusIndexerImpl{
		promIndexer: promIndexer,
		smonIndexer: smonIndexer,
		index:       index,
	}, nil
}

func (ri *PrometheusIndexerImpl) Configure(enable bool) {
	ri.locker.Lock()
	defer ri.locker.Unlock()

	ri.enable = enable
}

func (ri *PrometheusIndexerImpl) insert(key []byte, prometheus *prom.Prometheus) error {
	raw, err := ri.index.GetInternal(key)
	if err != nil || len(raw) == 0 {
		data := prom.PrometheusList{Items: []*prom.Prometheus{prometheus}}
		raw, err := json.Marshal(data)
		if err != nil {
			return err
		}
		err = ri.index.SetInternal(key, raw)
		if err != nil {
			return err
		}
	} else {
		var data prom.PrometheusList
		err := json.Unmarshal(raw, &data)
		if err != nil {
			return err
		}

		if found, _ := arrays.Contains(data.Items, prometheus); !found {
			data.Items = append(data.Items, prometheus)
			raw, err := json.Marshal(data)
			if err != nil {
				return err
			}
			err = ri.index.SetInternal(key, raw)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (ri *PrometheusIndexerImpl) remove(key []byte, prometheus *prom.Prometheus) error {
	raw, err := ri.index.GetInternal(key)
	if err != nil {
		return err
	}
	if len(raw) > 0 {
		var data prom.PrometheusList
		err := json.Unmarshal(raw, &data)
		if err != nil {
			return err
		}
		var prometheuss []*prom.Prometheus
		for i, valueSvc := range data.Items {
			if ri.equal(prometheus, valueSvc) {
				prometheuss = append(data.Items[:i], data.Items[i+1:]...)
				break
			}
		}

		if len(prometheuss) == 0 {
			// Remove unnecessary index
			err = ri.index.DeleteInternal(key)
			if err != nil {
				return err
			}
		} else {
			raw, err := json.Marshal(prom.PrometheusList{Items: prometheuss})
			if err != nil {
				return err
			}
			err = ri.index.SetInternal(key, raw)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (ri *PrometheusIndexerImpl) serviceMonitorsForPrometheus(prometheus *prom.Prometheus) ([]*prom.ServiceMonitor, error) {
	selector, err := metav1.LabelSelectorAsSelector(prometheus.Spec.ServiceMonitorSelector)
	if err != nil {
		return nil, err
	}

	var smons []*prom.ServiceMonitor
	err = cache.ListAllByNamespace(ri.smonIndexer, metav1.NamespaceAll, selector, func(m interface{}) {
		smons = append(smons, m.(*prom.ServiceMonitor))
	})
	return smons, err
}

func (ri *PrometheusIndexerImpl) equal(a, b *prom.Prometheus) bool {
	if a.Name == b.Name && a.Namespace == b.Namespace {
		return true
	}
	return false
}

func (ri *PrometheusIndexerImpl) Key(meta metav1.ObjectMeta) []byte {
	return []byte(kutil.GetGroupVersionKind(&prom.ServiceMonitor{}).String() + "/" + meta.Namespace + "/" + meta.Name)
}

func (ri *PrometheusIndexerImpl) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Infoln("Request received at", req.URL.Path)
	params, found := pat.FromContext(req.Context())
	if !found {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	namespace, name := params.Get(":namespace"), params.Get(":name")
	if len(namespace) > 0 && len(name) > 0 {
		key := ri.Key(v1.ObjectMeta{Name: name, Namespace: namespace})
		if val, err := ri.index.GetInternal(key); err == nil && len(val) > 0 {
			if err := json.NewEncoder(w).Encode(json.RawMessage(val)); err == nil {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("x-content-type-options", "nosniff")
				return
			} else {
				http.Error(w, "Server error"+err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.NotFound(w, req)
		}
		return
	}
	http.Error(w, "Bad Request", http.StatusBadRequest)
}

func (ri *PrometheusIndexerImpl) PrometheusHandler() cache.ResourceEventHandler {
	return &promPrometheusHandler{ri}
}

func (ri *PrometheusIndexerImpl) ServiceMonitorHandler() cache.ResourceEventHandler {
	return &promServiceMonitorHandler{ri}
}

type promPrometheusHandler struct {
	*PrometheusIndexerImpl
}

var _ cache.ResourceEventHandler = &promPrometheusHandler{}

func (ri *promPrometheusHandler) OnAdd(obj interface{}) {
	ri.locker.RLock()
	defer ri.locker.RUnlock()

	if !ri.enable {
		return
	}

	p, ok := obj.(*prom.Prometheus)
	if !ok {
		return
	}
	ri.add(p)
}

func (ri *promPrometheusHandler) add(prom *prom.Prometheus) {
	smons, err := ri.serviceMonitorsForPrometheus(prom)
	if err != nil {
		log.Errorln(err)
		return
	}
	for _, monitors := range smons {
		if err = ri.insert(ri.Key(monitors.ObjectMeta), prom); err != nil {
			log.Errorln(err)
		}
	}
}

func (ri *promPrometheusHandler) OnUpdate(oldObj, newObj interface{}) {
	ri.locker.RLock()
	defer ri.locker.RUnlock()

	if !ri.enable {
		return
	}

	oldRes, ok := oldObj.(*prom.Prometheus)
	if !ok {
		return
	}
	newRes, ok := newObj.(*prom.Prometheus)
	if !ok {
		return
	}

	if !reflect.DeepEqual(oldRes.Spec.ServiceMonitorSelector, newRes.Spec.ServiceMonitorSelector) {
		// Only update if selector changes
		ri.delete(oldRes)
		ri.add(newRes)
	}
}

func (ri *promPrometheusHandler) OnDelete(obj interface{}) {
	ri.locker.RLock()
	defer ri.locker.RUnlock()

	if !ri.enable {
		return
	}

	p, ok := obj.(*prom.Prometheus)
	if !ok {
		return
	}
	ri.delete(p)
}

func (ri *promPrometheusHandler) delete(prom *prom.Prometheus) {
	smons, err := ri.serviceMonitorsForPrometheus(prom)
	if err != nil {
		log.Errorln(err)
		return
	}

	for _, monitors := range smons {
		if err = ri.remove(ri.Key(monitors.ObjectMeta), prom); err != nil {
			log.Errorln(err)
		}
	}
}

type promServiceMonitorHandler struct {
	*PrometheusIndexerImpl
}

var _ cache.ResourceEventHandler = &promServiceMonitorHandler{}

func (ri *promServiceMonitorHandler) OnAdd(obj interface{}) {
	ri.locker.RLock()
	defer ri.locker.RUnlock()

	if !ri.enable {
		return
	}

	smon, ok := obj.(*prom.ServiceMonitor)
	if !ok {
		return
	}

	var proms []*prom.Prometheus
	err := cache.ListAllByNamespace(ri.promIndexer, core.NamespaceAll, labels.Everything(), func(m interface{}) {
		proms = append(proms, m.(*prom.Prometheus))
	})
	if err != nil {
		log.Errorln(err)
		return
	}

	key := ri.Key(smon.ObjectMeta)
	for _, prometheus := range proms {
		selector, err := metav1.LabelSelectorAsSelector(prometheus.Spec.ServiceMonitorSelector)
		if err != nil {
			continue
		}
		if labels.SelectorFromSet(labels.Set(smon.Labels)).String() != selector.String() {
			continue
		}

		if err = ri.insert(key, prometheus); err != nil {
			log.Errorln(err)
		}
	}
}

func (ri *promServiceMonitorHandler) OnDelete(obj interface{}) {
	ri.locker.RLock()
	defer ri.locker.RUnlock()

	if !ri.enable {
		return
	}

	smon, ok := obj.(*prom.ServiceMonitor)
	if !ok {
		return
	}

	if err := ri.index.DeleteInternal(ri.Key(smon.ObjectMeta)); err != nil {
		log.Errorln(err)
	}
}

func (ri *promServiceMonitorHandler) OnUpdate(oldObj, newObj interface{}) {}
