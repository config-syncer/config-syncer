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
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	core_lister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

type ServiceMonitorIndexer interface {
	Configure(enable bool)
	ServiceMonitorHandler() cache.ResourceEventHandler
	ServiceHandler() cache.ResourceEventHandler
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

var _ ServiceMonitorIndexer = &ServiceMonitorIndexerImpl{}

type ServiceMonitorIndexerImpl struct {
	svcLister   core_lister.ServiceLister
	smonIndexer cache.Indexer
	index       bleve.Index

	enable bool
	locker sync.RWMutex
}

func NewServiceMonitorIndexer(dir string, svcIndexer cache.Indexer, smonIndexer cache.Indexer) (ServiceMonitorIndexer, error) {
	index, err := ensureIndex(filepath.Join(dir, "smon.indexer"), "indexer")
	if err != nil {
		return nil, err
	}
	return &ServiceMonitorIndexerImpl{
		svcLister:   core_lister.NewServiceLister(svcIndexer),
		smonIndexer: smonIndexer,
		index:       index,
	}, nil
}

func (ri *ServiceMonitorIndexerImpl) Configure(enable bool) {
	ri.locker.Lock()
	defer ri.locker.Unlock()

	ri.enable = enable
}

func (ri *ServiceMonitorIndexerImpl) insert(key []byte, monitor *prom.ServiceMonitor) error {
	raw, err := ri.index.GetInternal(key)
	if err != nil || len(raw) == 0 {
		data := prom.ServiceMonitorList{Items: []*prom.ServiceMonitor{monitor}}
		raw, err := json.Marshal(data)
		if err != nil {
			return err
		}
		err = ri.index.SetInternal(key, raw)
		if err != nil {
			return err
		}
	} else {
		var data prom.ServiceMonitorList
		err := json.Unmarshal(raw, &data)
		if err != nil {
			return err
		}

		if found, _ := arrays.Contains(data.Items, monitor); !found {
			data.Items = append(data.Items, monitor)
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

func (ri *ServiceMonitorIndexerImpl) remove(key []byte, svcMonitor *prom.ServiceMonitor) error {
	raw, err := ri.index.GetInternal(key)
	if err != nil {
		return err
	}
	if len(raw) > 0 {
		var data prom.ServiceMonitorList
		err := json.Unmarshal(raw, &data)
		if err != nil {
			return err
		}
		var monitors []*prom.ServiceMonitor
		for i, valueSvc := range data.Items {
			if ri.equal(svcMonitor, valueSvc) {
				monitors = append(data.Items[:i], data.Items[i+1:]...)
				break
			}
		}

		if len(monitors) == 0 {
			// Remove unnecessary index
			err = ri.index.DeleteInternal(key)
			if err != nil {
				return err
			}
		} else {
			raw, err := json.Marshal(prom.ServiceMonitorList{Items: monitors})
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

func (ri *ServiceMonitorIndexerImpl) serviceForServiceMonitors(svcMonitor *prom.ServiceMonitor) ([]*core.Service, error) {
	selector, err := metav1.LabelSelectorAsSelector(&svcMonitor.Spec.Selector)
	if err != nil {
		return nil, err
	}
	if svcMonitor.Spec.NamespaceSelector.Any {
		return ri.svcLister.Services(metav1.NamespaceAll).List(selector)
	}

	var services []*core.Service
	for _, ns := range svcMonitor.Spec.NamespaceSelector.MatchNames {
		svc, err := ri.svcLister.Services(ns).List(selector)
		if err == nil {
			services = append(services, svc...)
		}
	}
	return services, nil
}

func (ri *ServiceMonitorIndexerImpl) equal(a, b *prom.ServiceMonitor) bool {
	if a.Name == b.Name && a.Namespace == b.Namespace {
		return true
	}
	return false
}

func (ri *ServiceMonitorIndexerImpl) Key(meta metav1.ObjectMeta) []byte {
	return []byte(kutil.GetGroupVersionKind(&core.Service{}).String() + "/" + meta.Namespace + "/" + meta.Name)
}

func (ri *ServiceMonitorIndexerImpl) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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

func (ri *ServiceMonitorIndexerImpl) ServiceMonitorHandler() cache.ResourceEventHandler {
	return &smonServiceMonitorHandler{ri}
}

func (ri *ServiceMonitorIndexerImpl) ServiceHandler() cache.ResourceEventHandler {
	return &smonServiceHandler{ri}
}

type smonServiceMonitorHandler struct {
	*ServiceMonitorIndexerImpl
}

var _ cache.ResourceEventHandler = &smonServiceMonitorHandler{}

func (ri *smonServiceMonitorHandler) OnAdd(obj interface{}) {
	ri.locker.RLock()
	defer ri.locker.RUnlock()

	if !ri.enable {
		return
	}

	smon, ok := obj.(*prom.ServiceMonitor)
	if !ok {
		return
	}
	log.Debugf("New svcMonitor: %v", smon.Name)
	ri.add(smon)
}

func (ri *smonServiceMonitorHandler) add(smon *prom.ServiceMonitor) {
	services, err := ri.serviceForServiceMonitors(smon)
	if err != nil {
		log.Errorln(err)
		return
	}
	for _, service := range services {
		if err = ri.insert(ri.Key(service.ObjectMeta), smon); err != nil {
			log.Errorln(err)
			return
		}
	}
}

func (ri *smonServiceMonitorHandler) OnDelete(obj interface{}) {
	ri.locker.RLock()
	defer ri.locker.RUnlock()

	if !ri.enable {
		return
	}

	smon, ok := obj.(*prom.ServiceMonitor)
	if !ok {
		return
	}
	ri.delete(smon)
}

func (ri *smonServiceMonitorHandler) delete(smon *prom.ServiceMonitor) {
	services, err := ri.serviceForServiceMonitors(smon)
	if err != nil {
		log.Errorln(err)
		return
	}
	for _, pod := range services {
		ri.remove(ri.Key(pod.ObjectMeta), smon)
	}
}

func (ri *smonServiceMonitorHandler) OnUpdate(oldObj, newObj interface{}) {
	ri.locker.RLock()
	defer ri.locker.RUnlock()

	if !ri.enable {
		return
	}

	oldRes, ok := oldObj.(*prom.ServiceMonitor)
	if !ok {
		return
	}

	newRes, ok := newObj.(*prom.ServiceMonitor)
	if !ok {
		return
	}

	if !reflect.DeepEqual(oldRes.Spec.Selector, newRes.Spec.Selector) {
		// Only update if selector changes
		ri.delete(oldRes)
		ri.add(newRes)
	}
}

type smonServiceHandler struct {
	*ServiceMonitorIndexerImpl
}

var _ cache.ResourceEventHandler = &smonServiceHandler{}

func (ri *smonServiceHandler) OnAdd(obj interface{}) {
	ri.locker.RLock()
	defer ri.locker.RUnlock()

	if !ri.enable {
		return
	}

	svc, ok := obj.(*core.Service)
	if !ok {
		return
	}

	var smons []*prom.ServiceMonitor
	cache.ListAllByNamespace(ri.smonIndexer, core.NamespaceAll, labels.Everything(), func(m interface{}) {
		smons = append(smons, m.(*prom.ServiceMonitor))
	})
	key := ri.Key(svc.ObjectMeta)
	for _, monitor := range smons {
		selector, err := metav1.LabelSelectorAsSelector(&monitor.Spec.Selector)
		if err != nil {
			continue
		}
		if labels.SelectorFromSet(labels.Set(svc.Labels)).String() != selector.String() {
			continue
		}
		if err = ri.insert(key, monitor); err != nil {
			log.Errorln(err)
		}
	}
}

func (ri *smonServiceHandler) OnUpdate(oldObj, newObj interface{}) {}

func (ri *smonServiceHandler) OnDelete(obj interface{}) {
	ri.locker.RLock()
	defer ri.locker.RUnlock()

	if !ri.enable {
		return
	}

	svc, ok := obj.(*core.Service)
	if !ok {
		return
	}

	if err := ri.index.DeleteInternal(ri.Key(svc.ObjectMeta)); err != nil {
		log.Errorln(err)
	}
}
