package indexers

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/appscode/go/arrays"
	"github.com/appscode/go/log"
	"github.com/appscode/kubed/pkg/util"
	"github.com/appscode/pat"
	"github.com/blevesearch/bleve"
	prom "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	clientset "k8s.io/client-go/kubernetes"
)

type PrometheusIndexer interface {
	Add(prom *prom.Prometheus) error
	Delete(prom *prom.Prometheus) error
	Update(old, new *prom.Prometheus) error
	AddServiceMonitor(*prom.ServiceMonitor, []*prom.Prometheus) error
	DeleteServiceMonitor(monitor *prom.ServiceMonitor) error
	Key(meta metav1.ObjectMeta) []byte
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

var _ PrometheusIndexer = &PrometheusIndexerImpl{}

type PrometheusIndexerImpl struct {
	kubeClient clientset.Interface
	promClient prom.MonitoringV1alpha1Interface
	index      bleve.Index
}

func (ri *PrometheusIndexerImpl) Add(prom *prom.Prometheus) error {
	log.Infof("New Prometheus: %v", prom.Name)
	log.V(5).Infof("Prometheus details: %v", prom)

	svcMonitors, err := ri.serviceMonitorsForPrometheus(prom)
	if err != nil {
		return err
	}

	for _, monitors := range svcMonitors.Items {
		key := ri.Key(monitors.ObjectMeta)
		ri.insert(key, prom)
	}
	return nil
}

func (ri *PrometheusIndexerImpl) Delete(prom *prom.Prometheus) error {
	svc, err := ri.serviceMonitorsForPrometheus(prom)
	if err != nil {
		return err
	}

	for _, monitors := range svc.Items {
		key := ri.Key(monitors.ObjectMeta)
		ri.remove(key, prom)
	}
	return nil
}

func (ri *PrometheusIndexerImpl) Update(old, new *prom.Prometheus) error {
	if !reflect.DeepEqual(old.Spec.ServiceMonitorSelector, new.Spec.ServiceMonitorSelector) {
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

func (ri *PrometheusIndexerImpl) AddServiceMonitor(m *prom.ServiceMonitor, prom []*prom.Prometheus) error {
	key := ri.Key(m.ObjectMeta)
	for _, prometheus := range prom {
		selector, err := metav1.LabelSelectorAsSelector(prometheus.Spec.ServiceMonitorSelector)
		if err != nil {
			continue
		}
		if labels.SelectorFromSet(labels.Set(m.Labels)).String() != selector.String() {
			continue
		}

		ri.insert(key, prometheus)
	}
	return nil
}

func (ri *PrometheusIndexerImpl) DeleteServiceMonitor(m *prom.ServiceMonitor) error {
	return ri.index.DeleteInternal(ri.Key(m.ObjectMeta))
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

func (ri *PrometheusIndexerImpl) serviceMonitorsForPrometheus(prometheus *prom.Prometheus) (*prom.ServiceMonitorList, error) {
	selector, err := metav1.LabelSelectorAsSelector(prometheus.Spec.ServiceMonitorSelector)
	if err != nil {
		return &prom.ServiceMonitorList{}, err
	}

	pods, err := ri.promClient.ServiceMonitors(metav1.NamespaceAll).List(metav1.ListOptions{
		LabelSelector: selector.String(),
	})
	if val, ok := pods.(*prom.ServiceMonitorList); ok {
		return val, nil
	}
	return &prom.ServiceMonitorList{}, err
}

func (ri *PrometheusIndexerImpl) equal(a, b *prom.Prometheus) bool {
	if a.Name == b.Name && a.Namespace == b.Namespace {
		return true
	}
	return false
}

func (ri *PrometheusIndexerImpl) Key(meta metav1.ObjectMeta) []byte {
	return []byte(util.GetGroupVersionKind(&prom.ServiceMonitor{}).String() + "/" + meta.Namespace + "/" + meta.Name)
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
