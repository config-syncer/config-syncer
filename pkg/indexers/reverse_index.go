package indexers

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/appscode/log"
	"github.com/blevesearch/bleve"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	clientset "k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

type ReverseIndexer struct {
	// kubeClient to access kube api server
	kubeClient clientset.Interface

	client bleve.Index

	// Channel serializes event to protect cache
	dataChan chan interface{}
}

func NewReverseIndexer(cl clientset.Interface, dst string) (*ReverseIndexer, error) {
	c, err := ensureIndex(strings.TrimRight(dst, "/")+"/reverse.indexer", "indexer")
	if err != nil {
		return nil, err
	}
	return &ReverseIndexer{
		kubeClient: cl,
		client:     c,
		dataChan:   make(chan interface{}, 1),
	}, nil
}

func (ri *ReverseIndexer) Handle(events string, obj ...interface{}) {
	switch obj[0].(type) {
	case *apiv1.Service:
		ri.handleService(events, obj...)
	}
}

func (ri *ReverseIndexer) handleService(events string, obj ...interface{}) {
	switch events {
	case "added":
		ri.dataChan <- obj[0]
		ri.newService()
	case "deleted":
		ri.dataChan <- obj[0]
		ri.removeService()
	case "updated":
		ri.updateService(obj[0], obj[1])
	default:
		log.Errorln("Event type not recognize", events)
	}
}

func (ri *ReverseIndexer) newService() {
	obj := <-ri.dataChan
	if service, ok := assertIsService(obj); ok {
		log.Infof("New service: %v", service.Name)
		log.V(5).Infof("Service details: %v", service)

		pods, err := ri.podsForService(service)
		if err != nil {
			log.Errorln("Failed to list Pods")
			return
		}

		for _, pod := range pods.Items {
			key := namespacerKey(pod.ObjectMeta)
			raw, err := ri.client.GetInternal(key)
			if err != nil || len(raw) == 0 {
				data := []*apiv1.Service{service}
				raw, err := json.Marshal(data)
				if err == nil {
					err := ri.client.SetInternal(key, raw)
					if err != nil {
						log.Errorln("Failed to store internal document", err)
					}
				}
			} else {
				var data []*apiv1.Service
				err := json.Unmarshal(raw, &data)
				if err == nil {
					data = append(data, service)
					raw, err := json.Marshal(data)
					if err == nil {
						err = ri.client.SetInternal(key, raw)
						if err != nil {
							log.Errorln("Failed to store internal document", err)
						}
					}
				}
			}
		}
	}
}

func (ri *ReverseIndexer) removeService() {
	obj := <-ri.dataChan
	if svc, ok := assertIsService(obj); ok {
		pods, err := ri.podsForService(svc)
		if err != nil {
			log.Errorln("Failed to list Pods")
			return
		}

		for _, pod := range pods.Items {
			key := namespacerKey(pod.ObjectMeta)
			raw, _ := ri.client.GetInternal(key)
			if len(raw) > 0 {
				var data []*apiv1.Service
				err := json.Unmarshal(raw, &data)
				if err == nil {
					tempData := data
					for i, valueSvc := range data {
						if equalService(svc, valueSvc) {
							tempData = append(data[:i], data[i+1:]...)
						}
					}

					if len(tempData) == 0 {
						// Remove unnecessary index
						ri.client.DeleteInternal(key)
					} else {
						raw, err := json.Marshal(tempData)
						if err == nil {
							ri.client.SetInternal(key, raw)
						}
					}
				}
			}
		}
	}
}

func (ri *ReverseIndexer) updateService(oldObj, newObj interface{}) {
	if old, ok := assertIsService(newObj); ok {
		if new, ok := assertIsService(oldObj); ok {
			if !reflect.DeepEqual(old.Spec.Selector, new.Spec.Selector) {
				// Only update if selector changes
				ri.dataChan <- oldObj
				ri.removeService()

				ri.dataChan <- newObj
				ri.newService()
			}
		}
	}
}

func (ri *ReverseIndexer) podsForService(svc *apiv1.Service) (*apiv1.PodList, error) {
	// Service have an empty selector. Instead of getting all pod we
	// try to ignore pods for it.
	if len(svc.Spec.Selector) == 0 {
		return &apiv1.PodList{}, nil
	}

	return ri.kubeClient.CoreV1().Pods(metav1.NamespaceAll).List(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(svc.Spec.Selector).String(),
	})
}

func equalService(a, b *apiv1.Service) bool {
	if a.Name == b.Name && a.Namespace == b.Namespace {
		return true
	}
	return false
}

func namespacerKey(meta metav1.ObjectMeta) []byte {
	return []byte(meta.Namespace + "/" + meta.Name)
}

func assertIsService(obj interface{}) (*apiv1.Service, bool) {
	if service, ok := obj.(*apiv1.Service); ok {
		return service, ok
	} else {
		log.Errorf("Type assertion failed! Expected 'Service', got %T", service)
		return nil, ok
	}
}
