package indexers

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"reflect"
	"sync"

	"github.com/appscode/go/arrays"
	"github.com/appscode/go/errors"
	"github.com/appscode/go/log"
	kutil "github.com/appscode/kutil/core/v1"
	"github.com/appscode/pat"
	"github.com/blevesearch/bleve"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	core_lister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

type ServiceIndexer interface {
	Configure(enable bool)
	ServiceHandler() cache.ResourceEventHandler
	EndpointHandler() cache.ResourceEventHandler
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

var _ ServiceIndexer = &ServiceIndexerImpl{}

type ServiceIndexerImpl struct {
	podLister core_lister.PodLister
	svcLister core_lister.ServiceLister
	index     bleve.Index

	enable bool
	locker sync.RWMutex
}

func NewServiceIndexer(dir string, podIndexer cache.Indexer, svcIndexer cache.Indexer) (ServiceIndexer, error) {
	index, err := ensureIndex(filepath.Join(dir, "service.indexer"), "indexer")
	if err != nil {
		return nil, err
	}
	return &ServiceIndexerImpl{
		podLister: core_lister.NewPodLister(podIndexer),
		svcLister: core_lister.NewServiceLister(svcIndexer),
		index:     index,
	}, nil

}

func (ri *ServiceIndexerImpl) Configure(enable bool) {
	ri.locker.Lock()
	defer ri.locker.Unlock()

	ri.enable = enable
}

func (ri *ServiceIndexerImpl) insert(key []byte, svc *core.Service) error {
	raw, err := ri.index.GetInternal(key)
	if err != nil || len(raw) == 0 {
		data := core.ServiceList{Items: []core.Service{*svc}}
		raw, err := json.Marshal(data)
		if err != nil {
			return err
		}
		err = ri.index.SetInternal(key, raw)
		if err != nil {
			return err
		}
	} else {
		var data core.ServiceList
		err := json.Unmarshal(raw, &data)
		if err != nil {
			return err
		}

		if found, _ := arrays.Contains(data.Items, *svc); !found {
			data.Items = append(data.Items, *svc)
			raw, err := json.Marshal(data)
			if err != nil {
				return err
			}
			err = ri.index.SetInternal(key, raw)
			if err != nil {
				return errors.FromErr(err).WithMessage("Failed to store internal document").Err()
			}
		}
	}
	return nil
}

func (ri ServiceIndexerImpl) remove(key []byte, svc *core.Service) error {
	raw, err := ri.index.GetInternal(key)
	if err != nil {
		return err
	}
	if len(raw) > 0 {
		var data core.ServiceList
		err := json.Unmarshal(raw, &data)
		if err != nil {
			return err
		}
		var ni []core.Service
		for i, valueSvc := range data.Items {
			if ri.equal(svc, &valueSvc) {
				ni = append(data.Items[:i], data.Items[i+1:]...)
				break
			}
		}

		if len(ni) == 0 {
			// Remove unnecessary index
			err = ri.index.DeleteInternal(key)
			if err != nil {
				return err
			}
		} else {
			raw, err := json.Marshal(core.ServiceList{Items: ni})
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

func (ri *ServiceIndexerImpl) podsForService(svc *core.Service) ([]*core.Pod, error) {
	// Service have an empty selector. Instead of getting all pod we
	// try to ignore pods for it.
	if len(svc.Spec.Selector) == 0 {
		return nil, nil
	}
	return ri.podLister.Pods(metav1.NamespaceAll).List(labels.SelectorFromSet(svc.Spec.Selector))
}

func (ri *ServiceIndexerImpl) equal(a, b *core.Service) bool {
	if a.Name == b.Name && a.Namespace == b.Namespace {
		return true
	}
	return false
}

func (ri *ServiceIndexerImpl) Key(meta metav1.ObjectMeta) []byte {
	return []byte(kutil.GetGroupVersionKind(&core.Pod{}).String() + "/" + meta.Namespace + "/" + meta.Name)
}

func (ri *ServiceIndexerImpl) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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

func (ri *ServiceIndexerImpl) ServiceHandler() cache.ResourceEventHandler {
	return &svcServiceHandler{ri}
}

func (ri *ServiceIndexerImpl) EndpointHandler() cache.ResourceEventHandler {
	return &svcEndpointHandler{ri}
}

type svcServiceHandler struct {
	*ServiceIndexerImpl
}

var _ cache.ResourceEventHandler = &svcServiceHandler{}

func (ri *svcServiceHandler) OnAdd(obj interface{}) {
	ri.locker.RLock()
	defer ri.locker.RUnlock()

	if !ri.enable {
		return
	}

	svc, ok := obj.(*core.Service)
	if !ok {
		return
	}
	log.Debugf("New service: %v", svc.Name)

	ri.add(svc)
	return
}

func (ri *svcServiceHandler) add(svc *core.Service) {
	pods, err := ri.podsForService(svc)
	if err != nil {
		log.Errorln(err)
		return
	}

	for _, pod := range pods {
		key := ri.Key(pod.ObjectMeta)
		ri.insert(key, svc)
	}
}

func (ri *svcServiceHandler) OnUpdate(oldObj, newObj interface{}) {
	ri.locker.RLock()
	defer ri.locker.RUnlock()

	if !ri.enable {
		return
	}

	oldRes, ok := oldObj.(*core.Service)
	if !ok {
		return
	}

	newRes, ok := newObj.(*core.Service)
	if !ok {
		return
	}

	if !reflect.DeepEqual(oldRes.Spec.Selector, newRes.Spec.Selector) {
		// Only update if selector changes
		ri.delete(oldRes)
		ri.add(newRes)
	}
	return
}

func (ri *svcServiceHandler) OnDelete(obj interface{}) {
	ri.locker.RLock()
	defer ri.locker.RUnlock()

	if !ri.enable {
		return
	}

	svc, ok := obj.(*core.Service)
	if !ok {
		return
	}

	ri.delete(svc)
}

func (ri *svcServiceHandler) delete(svc *core.Service) {
	pods, err := ri.podsForService(svc)
	if err != nil {
		log.Errorln(err)
		return
	}

	for _, pod := range pods {
		key := ri.Key(pod.ObjectMeta)
		ri.remove(key, svc)
	}
}

type svcEndpointHandler struct {
	*ServiceIndexerImpl
}

var _ cache.ResourceEventHandler = &svcEndpointHandler{}

func (ri svcEndpointHandler) OnAdd(obj interface{})    {}
func (ri svcEndpointHandler) OnDelete(obj interface{}) {}

func (ri svcEndpointHandler) OnUpdate(oldObj, newObj interface{}) {
	ri.locker.RLock()
	defer ri.locker.RUnlock()

	if !ri.enable {
		return
	}

	oldRes, ok := oldObj.(*core.Endpoints)
	if !ok {
		return
	}
	newRes, ok := newObj.(*core.Endpoints)
	if !ok {
		return
	}
	if reflect.DeepEqual(oldRes.Subsets, newRes.Subsets) {
		return
	}

	svc, err := ri.svcLister.Services(newRes.Namespace).Get(newRes.Name)
	if err != nil {
		log.Errorln(err)
		return
	}

	oldPods := make(map[string]*core.Pod)
	for _, oldEPSubsets := range oldRes.Subsets {
		for _, oldEPPods := range oldEPSubsets.Addresses {
			if podRef := oldEPPods.TargetRef; podRef != nil {
				pod := &core.Pod{ObjectMeta: metav1.ObjectMeta{Name: podRef.Name, Namespace: podRef.Namespace}}
				oldPods[podRef.String()] = pod
			}
		}
	}

	newPods := make(map[string]*core.Pod)
	for _, newEPSubsets := range newRes.Subsets {
		for _, newEPPods := range newEPSubsets.Addresses {
			if podRef := newEPPods.TargetRef; podRef != nil {
				pod := &core.Pod{ObjectMeta: metav1.ObjectMeta{Name: podRef.Name, Namespace: podRef.Namespace}}
				newPods[podRef.String()] = pod
				if _, ok := oldPods[podRef.String()]; !ok {
					// This Pod reference is in update Endpoint, New Pod Added
					if err = ri.insert(ri.Key(pod.ObjectMeta), svc); err != nil {
						log.Errorln(err)
					}
				}
			}
		}
	}

	for ref, pod := range oldPods {
		if _, ok := newPods[ref]; !ok {
			// Pod ref not found in New Endpoint, Removed
			if err = ri.remove(ri.Key(pod.ObjectMeta), svc); err != nil {
				log.Errorln(err)
			}
		}
	}
}
