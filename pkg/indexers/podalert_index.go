package indexers

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/appscode/go/arrays"
	"github.com/appscode/kubed/pkg/util"
	"github.com/appscode/log"
	"github.com/appscode/pat"
	searchlight "github.com/appscode/searchlight/api"
	searchlightclient "github.com/appscode/searchlight/client/clientset"
	"github.com/blevesearch/bleve"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	clientset "k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

type PodAlertIndexer interface {
	Add(podAlert *searchlight.PodAlert) error
	Delete(podAlert *searchlight.PodAlert) error
	Update(old, new *searchlight.PodAlert) error
	Key(meta metav1.ObjectMeta) []byte
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

var _ PodAlertIndexer = &PodAlertIndexerImpl{}

type PodAlertIndexerImpl struct {
	kubeClient  clientset.Interface
	alertClient searchlightclient.ExtensionInterface
	index       bleve.Index
}

func (ri *PodAlertIndexerImpl) Add(podAlert *searchlight.PodAlert) error {
	log.Infof("New PodAlert: %v", podAlert.Name)
	log.V(5).Infof("PodAlert details: %v", podAlert)

	pods, err := ri.podsForPodAlert(podAlert)
	if err != nil {
		return err
	}

	for _, pod := range pods.Items {
		key := ri.Key(pod.ObjectMeta)
		ri.insert(key, *podAlert)
	}
	return nil
}

func (ri *PodAlertIndexerImpl) Delete(podAlert *searchlight.PodAlert) error {
	pod, err := ri.podsForPodAlert(podAlert)
	if err != nil {
		return err
	}

	for _, monitors := range pod.Items {
		key := ri.Key(monitors.ObjectMeta)
		ri.remove(key, *podAlert)
	}
	return nil
}

func (ri *PodAlertIndexerImpl) Update(old, new *searchlight.PodAlert) error {
	if !reflect.DeepEqual(old.Spec.Selector, new.Spec.Selector) || (old.Spec.PodName != old.Spec.PodName) {
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

//func (ri *PodAlertIndexerImpl) AddServiceMonitor(m *searchlight.ServiceMonitor, podAlert []*searchlight.PodAlert) error {
//	key := ri.Key(m.ObjectMeta)
//	for _, podAlert := range podAlert {
//		selector, err := metav1.LabelSelectorAsSelector(podAlert.Spec.ServiceMonitorSelector)
//		if err != nil {
//			continue
//		}
//		if labels.SelectorFromSet(labels.Set(m.Labels)).String() != selector.String() {
//			continue
//		}
//
//		ri.insert(key, podAlert)
//	}
//	return nil
//}
//
//func (ri *PodAlertIndexerImpl) DeleteServiceMonitor(m *searchlight.ServiceMonitor) error {
//	return ri.index.DeleteInternal(ri.Key(m.ObjectMeta))
//}

func (ri *PodAlertIndexerImpl) insert(key []byte, podAlert searchlight.PodAlert) error {
	raw, err := ri.index.GetInternal(key)
	if err != nil || len(raw) == 0 {
		data := searchlight.PodAlertList{Items: []searchlight.PodAlert{podAlert}}
		raw, err := json.Marshal(data)
		if err != nil {
			return err
		}
		err = ri.index.SetInternal(key, raw)
		if err != nil {
			return err
		}
	} else {
		var data searchlight.PodAlertList
		err := json.Unmarshal(raw, &data)
		if err != nil {
			return err
		}

		if found, _ := arrays.Contains(data.Items, podAlert); !found {
			data.Items = append(data.Items, podAlert)
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

func (ri *PodAlertIndexerImpl) remove(key []byte, podAlert searchlight.PodAlert) error {
	raw, err := ri.index.GetInternal(key)
	if err != nil {
		return err
	}
	if len(raw) > 0 {
		var data searchlight.PodAlertList
		err := json.Unmarshal(raw, &data)
		if err != nil {
			return err
		}
		var podAlerts []searchlight.PodAlert
		for i, value := range data.Items {
			if ri.equal(podAlert, value) {
				podAlerts = append(data.Items[:i], data.Items[i+1:]...)
				break
			}
		}

		if len(podAlerts) == 0 {
			// Remove unnecessary index
			err = ri.index.DeleteInternal(key)
			if err != nil {
				return err
			}
		} else {
			raw, err := json.Marshal(searchlight.PodAlertList{Items: podAlerts})
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

func (ri *PodAlertIndexerImpl) podsForPodAlert(podAlert *searchlight.PodAlert) (*apiv1.PodList, error) {
	// Following code is a reflection of searchlight PodAlert to Pod Selector logic
	// https://github.com/appscode/searchlight/blob/master/pkg/operator/pod_alerts.go#L153
	// TODO: Change if something changes in upstream
	newSel, err := metav1.LabelSelectorAsSelector(&podAlert.Spec.Selector)
	if err != nil {
		return &apiv1.PodList{}, err
	}
	if podAlert.Spec.PodName != "" {
		if resource, err := ri.kubeClient.CoreV1().Pods(podAlert.Namespace).Get(podAlert.Spec.PodName, metav1.GetOptions{}); err == nil {
			if newSel.Matches(labels.Set(resource.Labels)) {
				return &apiv1.PodList{Items: []apiv1.Pod{*resource}}, nil
			}
		}
	}

	if resources, err := ri.kubeClient.CoreV1().Pods(podAlert.Namespace).List(metav1.ListOptions{LabelSelector: newSel.String()}); err == nil {
		if err != nil {
			return resources, nil
		}
	}
	return &apiv1.PodList{}, nil
}

func (ri *PodAlertIndexerImpl) equal(a, b searchlight.PodAlert) bool {
	if a.Name == b.Name && a.Namespace == b.Namespace {
		return true
	}
	return false
}

func (ri *PodAlertIndexerImpl) Key(meta metav1.ObjectMeta) []byte {
	return []byte(util.GetGroupVersionKind(&apiv1.Pod{}).String() + "/" + meta.Namespace + "/" + meta.Name + "/" + util.GetGroupVersionKind(&searchlight.PodAlert{}).String())
}

func (ri *PodAlertIndexerImpl) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Infoln("Request received at", req.URL.Path)
	params, found := pat.FromContext(req.Context())
	if !found {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	namespace, name := params.Get(":namespace"), params.Get(":name")
	if len(namespace) > 0 && len(name) > 0 {
		key := ri.Key(metav1.ObjectMeta{Name: name, Namespace: namespace})
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
