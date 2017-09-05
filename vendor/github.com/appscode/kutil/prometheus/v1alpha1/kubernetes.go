package v1alpha1

import (
	"errors"

	prom "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func GetGroupVersionKind(v interface{}) schema.GroupVersionKind {
	switch v.(type) {
	case *prom.Prometheus:
		return schema.GroupVersionKind{Group: prom.Group, Version: prom.Version, Kind: prom.PrometheusesKind}
	case *prom.ServiceMonitor:
		return schema.GroupVersionKind{Group: prom.Group, Version: prom.Version, Kind: prom.ServiceMonitorsKind}
	case *prom.Alertmanager:
		return schema.GroupVersionKind{Group: prom.Group, Version: prom.Version, Kind: prom.AlertmanagersKind}
	default:
		return schema.GroupVersionKind{}
	}
}

func AssignTypeKind(v interface{}) error {
	switch u := v.(type) {
	case *prom.Prometheus:
		u.APIVersion = schema.GroupVersion{Group: prom.Group, Version: prom.Version}.String()
		u.Kind = prom.PrometheusesKind
		return nil
	case *prom.ServiceMonitor:
		u.APIVersion = schema.GroupVersion{Group: prom.Group, Version: prom.Version}.String()
		u.Kind = prom.ServiceMonitorsKind
		return nil
	case *prom.Alertmanager:
		u.APIVersion = schema.GroupVersion{Group: prom.Group, Version: prom.Version}.String()
		u.Kind = prom.AlertmanagersKind
		return nil
	}
	return errors.New("Unknown api object type")
}
