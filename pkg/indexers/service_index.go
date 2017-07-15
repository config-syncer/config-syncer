package indexers

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/appscode/log"
	"github.com/appscode/pat"
	"github.com/blevesearch/bleve"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	clientset "k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"github.com/appscode/kubed/pkg/util"
)

type ServiceIndexer interface {
	Add(svc *apiv1.Service)
	Delete(svc *apiv1.Service)
	Update(old, new *apiv1.Service)
	Key(meta metav1.ObjectMeta) []byte
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

var _ ServiceIndexer = &ServiceIndexerImpl{}

type ServiceIndexerImpl struct {
	// kubeClient to access kube api server
	kubeClient clientset.Interface

	index bleve.Index
}

func (ri *ServiceIndexerImpl) Add(svc *apiv1.Service) {
	log.Infof("New service: %v", svc.Name)
	log.V(5).Infof("Service details: %v", svc)

	pods, err := ri.podsForService(svc)
	if err != nil {
		log.Errorln("Failed to list Pods")
		return
	}

	for _, pod := range pods.Items {
		key := ri.Key(pod.ObjectMeta)
		raw, err := ri.index.GetInternal(key)
		if err != nil || len(raw) == 0 {
			data := []*apiv1.Service{svc}
			raw, err := json.Marshal(data)
			if err == nil {
				err := ri.index.SetInternal(key, raw)
				if err != nil {
					log.Errorln("Failed to store internal document", err)
				}
			}
		} else {
			var data []*apiv1.Service
			err := json.Unmarshal(raw, &data)
			if err == nil {
				data = append(data, svc)
				raw, err := json.Marshal(data)
				if err == nil {
					err = ri.index.SetInternal(key, raw)
					if err != nil {
						log.Errorln("Failed to store internal document", err)
					}
				}
			}
		}
	}
}

func (ri *ServiceIndexerImpl) Delete(svc *apiv1.Service) {
	pods, err := ri.podsForService(svc)
	if err != nil {
		log.Errorln("Failed to list Pods")
		return
	}

	for _, pod := range pods.Items {
		key := ri.Key(pod.ObjectMeta)
		raw, _ := ri.index.GetInternal(key)
		if len(raw) > 0 {
			var data []*apiv1.Service
			err := json.Unmarshal(raw, &data)
			if err == nil {
				tempData := data
				for i, valueSvc := range data {
					if ri.equal(svc, valueSvc) {
						tempData = append(data[:i], data[i+1:]...)
					}
				}

				if len(tempData) == 0 {
					// Remove unnecessary index
					ri.index.DeleteInternal(key)
				} else {
					raw, err := json.Marshal(tempData)
					if err == nil {
						ri.index.SetInternal(key, raw)
					}
				}
			}
		}
	}
}

func (ri *ServiceIndexerImpl) Update(old, new *apiv1.Service) {
	if !reflect.DeepEqual(old.Spec.Selector, new.Spec.Selector) {
		// Only update if selector changes
		ri.Delete(old)
		ri.Add(new)
	}
}

func (ri *ServiceIndexerImpl) podsForService(svc *apiv1.Service) (*apiv1.PodList, error) {
	// Service have an empty selector. Instead of getting all pod we
	// try to ignore pods for it.
	if len(svc.Spec.Selector) == 0 {
		return &apiv1.PodList{}, nil
	}

	return ri.kubeClient.CoreV1().Pods(metav1.NamespaceAll).List(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(svc.Spec.Selector).String(),
	})
}

func (ri *ServiceIndexerImpl) equal(a, b *apiv1.Service) bool {
	if a.Name == b.Name && a.Namespace == b.Namespace {
		return true
	}
	return false
}

func (ri *ServiceIndexerImpl) Key(meta metav1.ObjectMeta) []byte {
	return []byte(util.GetGroupVersionKind(&apiv1.Pod{}).String() + "/" + meta.Namespace + "/" + meta.Name)
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
