package indexers

import (
	"encoding/json"
	"path/filepath"
	"reflect"

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
}

func NewReverseIndexer(cl clientset.Interface, dst string) (*ReverseIndexer, error) {
	c, err := ensureIndex(filepath.Join(dst, "reverse.indexer"), "indexer")
	if err != nil {
		return nil, err
	}
	return &ReverseIndexer{
		kubeClient: cl,
		client:     c,
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
		ri.AddService(obj[0])
	case "deleted":
		ri.RemoveService(obj[0])
	case "updated":
		ri.UpdateService(obj[0], obj[1])
	default:
		log.Errorln("Event type not recognize", events)
	}
}

func (ri *ReverseIndexer) AddService(svc *apiv1.Service) {
	log.Infof("New service: %v", svc.Name)
	log.V(5).Infof("Service details: %v", svc)

	pods, err := ri.podsForService(svc)
	if err != nil {
		log.Errorln("Failed to list Pods")
		return
	}

	for _, pod := range pods.Items {
		key := namespacerKey(pod.ObjectMeta)
		raw, err := ri.client.GetInternal(key)
		if err != nil || len(raw) == 0 {
			data := []*apiv1.Service{svc}
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
				data = append(data, svc)
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

func (ri *ReverseIndexer) RemoveService(svc *apiv1.Service) {
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

func (ri *ReverseIndexer) UpdateService(old, new *apiv1.Service) {
	if !reflect.DeepEqual(old.Spec.Selector, new.Spec.Selector) {
		// Only update if selector changes
		ri.RemoveService(old)
		ri.AddService(new)
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
