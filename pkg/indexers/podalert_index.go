package indexers

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/appscode/kubed/pkg/util"
	"github.com/appscode/log"
	"github.com/appscode/pat"
	searchlight "github.com/appscode/searchlight/api"
	searchlightclient "github.com/appscode/searchlight/client/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	clientset "k8s.io/client-go/kubernetes"
	apps "k8s.io/client-go/pkg/apis/apps/v1beta1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

type PodAlertIndexer interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

var _ PodAlertIndexer = &PodAlertIndexerImpl{}

type PodAlertIndexerImpl struct {
	kubeClient  clientset.Interface
	alertClient searchlightclient.ExtensionInterface
}

func (ri *PodAlertIndexerImpl) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Infoln("Request received at", req.URL.Path)
	params, found := pat.FromContext(req.Context())
	if !found {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	resource, namespace, name := params.Get(":resource"), params.Get(":namespace"), params.Get(":name")
	if len(resource) > 0 && len(namespace) > 0 && len(name) > 0 {
		var selector *metav1.LabelSelector
		var podName string
		switch resource {
		case "deployments":
			if util.IsPreferredAPIResource(ri.kubeClient, apps.SchemeGroupVersion.String(), "Deployment") {
				res, err := ri.kubeClient.AppsV1beta1().Deployments(namespace).Get(name, metav1.GetOptions{})
				if err != nil {
					http.Error(w, "Server error"+err.Error(), http.StatusInternalServerError)
					return
				}
				selector = res.Spec.Selector
			} else if util.IsPreferredAPIResource(ri.kubeClient, extensions.SchemeGroupVersion.String(), "Deployment") {
				res, err := ri.kubeClient.ExtensionsV1beta1().Deployments(namespace).Get(name, metav1.GetOptions{})
				if err != nil {
					http.Error(w, "Server error"+err.Error(), http.StatusInternalServerError)
					return
				}
				selector = res.Spec.Selector
			}
		case "replicasets":
			res, err := ri.kubeClient.ExtensionsV1beta1().ReplicaSets(namespace).Get(name, metav1.GetOptions{})
			if err != nil {
				http.Error(w, "Server error"+err.Error(), http.StatusInternalServerError)
				return
			}
			selector = res.Spec.Selector
		case "replicationcontrollers":
			res, err := ri.kubeClient.CoreV1().ReplicationControllers(namespace).Get(name, metav1.GetOptions{})
			if err != nil {
				http.Error(w, "Server error"+err.Error(), http.StatusInternalServerError)
				return
			}
			selector = metav1.SetAsLabelSelector(labels.Set(res.Spec.Selector))
		case "daemonsets":
			res, err := ri.kubeClient.ExtensionsV1beta1().DaemonSets(namespace).Get(name, metav1.GetOptions{})
			if err != nil {
				http.Error(w, "Server error"+err.Error(), http.StatusInternalServerError)
				return
			}
			selector = res.Spec.Selector
		case "statefulsets":
			res, err := ri.kubeClient.AppsV1beta1().StatefulSets(namespace).Get(name, metav1.GetOptions{})
			if err != nil {
				http.Error(w, "Server error"+err.Error(), http.StatusInternalServerError)
				return
			}
			selector = res.Spec.Selector
		case "pods":
			res, err := ri.kubeClient.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
			if err != nil {
				http.Error(w, "Server error"+err.Error(), http.StatusInternalServerError)
				return
			}
			podName = res.Name
			selector = metav1.SetAsLabelSelector(labels.Set(res.Labels))
		}

		if selector != nil || len(podName) > 0 {
			podAlerts, err := ri.alertClient.PodAlerts(namespace).List(metav1.ListOptions{})
			if err != nil {
				http.Error(w, "Server error"+err.Error(), http.StatusInternalServerError)
				return
			}

			resultSet := make([]searchlight.PodAlert, 0)
			for _, alert := range podAlerts.Items {
				if reflect.DeepEqual(&selector, alert.Spec.Selector) {
					if resource == "pods" && podName != alert.Spec.PodName {
						continue
					}
					resultSet = append(resultSet, alert)
				}
			}

			if err := json.NewEncoder(w).Encode(resultSet); err == nil {
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
