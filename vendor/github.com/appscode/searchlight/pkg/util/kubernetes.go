package util

import (
	tapi "github.com/appscode/searchlight/api"
	tcs "github.com/appscode/searchlight/client/clientset"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func FindPodAlert(stashClient tcs.ExtensionInterface, obj metav1.ObjectMeta) ([]*tapi.PodAlert, error) {
	alerts, err := stashClient.PodAlerts(obj.Namespace).List(metav1.ListOptions{LabelSelector: labels.Everything().String()})
	if kerr.IsNotFound(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	result := make([]*tapi.PodAlert, 0)
	for i, alert := range alerts.Items {
		if ok, _ := alert.IsValid(); !ok {
			continue
		}
		if alert.Spec.PodName != "" && alert.Spec.PodName != obj.Name {
			continue
		}
		if selector, err := metav1.LabelSelectorAsSelector(&alert.Spec.Selector); err == nil {
			if selector.Matches(labels.Set(obj.Labels)) {
				result = append(result, &alerts.Items[i])
			}
		}
	}
	return result, nil
}

func FindNodeAlert(stashClient tcs.ExtensionInterface, obj metav1.ObjectMeta) ([]*tapi.NodeAlert, error) {
	alerts, err := stashClient.NodeAlerts(obj.Namespace).List(metav1.ListOptions{LabelSelector: labels.Everything().String()})
	if kerr.IsNotFound(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	result := make([]*tapi.NodeAlert, 0)
	for i, alert := range alerts.Items {
		if ok, _ := alert.IsValid(); !ok {
			continue
		}
		if alert.Spec.NodeName != "" && alert.Spec.NodeName != obj.Name {
			continue
		}
		selector := labels.SelectorFromSet(alert.Spec.Selector)
		if selector.Matches(labels.Set(obj.Labels)) {
			result = append(result, &alerts.Items[i])
		}
	}
	return result, nil
}
