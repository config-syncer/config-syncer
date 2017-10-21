package indexers

import (
	"encoding/json"
	"net/http"
	"reflect"

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
	clientset "k8s.io/client-go/kubernetes"
)

type ServiceIndexer interface {
	Add(svc *core.Service) error
	Delete(svc *core.Service) error
	AddPodForService(svc *core.Service, pod *core.Pod) error
	DeletePodForService(svc *core.Service, pod *core.Pod) error
	Update(old, new *core.Service) error
	Key(meta metav1.ObjectMeta) []byte
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

var _ ServiceIndexer = &ServiceIndexerImpl{}

type ServiceIndexerImpl struct {
	kubeClient clientset.Interface
	index      bleve.Index
}

func (ri *ServiceIndexerImpl) Add(svc *core.Service) error {
	log.Infof("New service: %v", svc.Name)
	log.V(5).Infof("Service details: %v", svc)

	pods, err := ri.podsForService(svc)
	if err != nil {
		return err
	}

	for _, pod := range pods.Items {
		key := ri.Key(pod.ObjectMeta)
		ri.insert(key, svc)
	}
	return nil
}

func (ri *ServiceIndexerImpl) Delete(svc *core.Service) error {
	pods, err := ri.podsForService(svc)
	if err != nil {
		return err
	}

	for _, pod := range pods.Items {
		key := ri.Key(pod.ObjectMeta)
		ri.remove(key, svc)
	}
	return nil
}

func (ri *ServiceIndexerImpl) Update(old, new *core.Service) error {
	if !reflect.DeepEqual(old.Spec.Selector, new.Spec.Selector) {
		// Only update if selector changes
		err := ri.Delete(old)
		if err != nil {
			return err
		}
		err = ri.Add(new)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ri *ServiceIndexerImpl) AddPodForService(svc *core.Service, pod *core.Pod) error {
	key := ri.Key(svc.ObjectMeta)
	return ri.insert(key, svc)
}

func (ri *ServiceIndexerImpl) DeletePodForService(svc *core.Service, pod *core.Pod) error {
	return ri.remove(ri.Key(pod.ObjectMeta), svc)
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

func (ri *ServiceIndexerImpl) podsForService(svc *core.Service) (*core.PodList, error) {
	// Service have an empty selector. Instead of getting all pod we
	// try to ignore pods for it.
	if len(svc.Spec.Selector) == 0 {
		return &core.PodList{}, nil
	}

	return ri.kubeClient.CoreV1().Pods(metav1.NamespaceAll).List(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(svc.Spec.Selector).String(),
	})
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
