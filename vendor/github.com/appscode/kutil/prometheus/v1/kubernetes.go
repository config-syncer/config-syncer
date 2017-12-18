package v1

import (
	"errors"

	"github.com/appscode/kutil/meta"
	prom "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func GetGroupVersionKind(v interface{}) schema.GroupVersionKind {
	return schema.GroupVersionKind{Group: prom.Group, Version: prom.Version, Kind: meta.GetKind(v)}
}

func AssignTypeKind(v interface{}) error {
	_, err := conversion.EnforcePtr(v)
	if err != nil {
		return err
	}

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
	return errors.New("unknown api object type")
}
