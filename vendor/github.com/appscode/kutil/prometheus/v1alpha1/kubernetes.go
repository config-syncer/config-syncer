package v1alpha1

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/appscode/kutil"
	prom "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func GetGroupVersionKind(v interface{}) schema.GroupVersionKind {
	return schema.GroupVersionKind{Group: prom.Group, Version: prom.Version, Kind: kutil.GetKind(v)}
}

func AssignTypeKind(v interface{}) error {
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return fmt.Errorf("%v must be a pointer", v)
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
