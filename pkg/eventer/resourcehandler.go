package eventer

import (
	"fmt"
	"time"

	"github.com/appscode/go/log"
	stringz "github.com/appscode/go/strings"
	"github.com/appscode/kubed/pkg/config"
	"github.com/appscode/kutil/discovery"
	meta_util "github.com/appscode/kutil/meta"
	"github.com/ghodss/yaml"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/cache"
)

const (
	MaxSyncInterval = 10 * time.Minute
)

var _ cache.ResourceEventHandler = &EventForwarder{}

func (f *EventForwarder) OnAdd(obj interface{}) {
	f.lock.RLock()
	defer f.lock.RUnlock()

	if err := f.forward(config.Create, obj); err != nil {
		log.Errorln(err)
		return
	}
}

func (f *EventForwarder) OnUpdate(oldObj, newObj interface{}) {}

func (f *EventForwarder) OnDelete(obj interface{}) {
	f.lock.RLock()
	defer f.lock.RUnlock()

	if err := f.forward(config.Create, obj); err != nil {
		log.Errorln(err)
		return
	}
}

func recentEvent(t metav1.Time) bool {
	return time.Now().Sub(t.Time) < MaxSyncInterval
}

// Check whether the rule's resource fields match the request attrs.
func ruleMatchesResource(r config.PolicyRule, attrs attributes) bool {
	if len(r.Namespaces) > 0 {
		if !hasString(r.Namespaces, attrs.accessor.GetNamespace()) { // Non-namespaced resources use the empty string.
			return false
		}
	}
	if len(r.Resources) == 0 {
		return true
	}

	apiGroup := attrs.gvr.Group
	resource := attrs.gvr.Resource

	name := attrs.accessor.GetName()

	for _, gr := range r.Resources {
		if gr.Group == apiGroup {
			if len(gr.Resources) == 0 {
				return true
			}
			for _, res := range gr.Resources {
				if res == resource {
					if len(gr.ResourceNames) == 0 || hasString(gr.ResourceNames, name) {
						return true
					}
				}
			}
		}
	}
	return false
}

type attributes struct {
	gvr      schema.GroupVersionResource
	op       config.Operation
	accessor metav1.Object
}

// Utility function to check whether a string slice contains a string.
func hasString(slice []string, value string) bool {
	for _, s := range slice {
		if s == value {
			return true
		}
	}
	return false
}

func (f *EventForwarder) forward(op config.Operation, obj interface{}) error {
	if f.spec == nil {
		return nil
	}

	accessor, err := meta.Accessor(obj)
	if err != nil {
		return err
	}

	if !recentEvent(accessor.GetCreationTimestamp()) {
		return nil
	}

	resource, err := discovery.DetectResource(f.restmapper, obj)
	if err != nil {
		return err
	}
	gvk := resource.GroupVersion().WithKind(meta_util.GetKind(obj))

	attrs := attributes{
		gvr:      resource,
		op:       config.Create,
		accessor: accessor,
	}

	matches := false
	for _, rule := range f.spec.Rules {
		if ruleMatchesResource(rule, attrs) {
			matches = true
			break
		}
	}
	if !matches {
		return nil
	}

	data, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}
	switch e := obj.(type) {
	case *core.Event:
		host := ""
		if e.Source.Host != "" {
			host = "on host " + e.Source.Host
		}
		for _, receiver := range f.spec.Receivers {
			emailSub := fmt.Sprintf("[%s, %s]: %s %s/%s %s %s", stringz.Val(f.clusterName, "?"), e.Source.Component, e.InvolvedObject.Kind, e.InvolvedObject.Namespace, e.InvolvedObject.Name, e.Reason, host)
			chatSub := fmt.Sprintf("[%s, %s] %s %s/%s %s %s: %s", stringz.Val(f.clusterName, "?"), e.Source.Component, e.InvolvedObject.Kind, e.InvolvedObject.Namespace, e.InvolvedObject.Name, e.Reason, host, e.Message)
			if err := f.notify(emailSub, chatSub, string(data), receiver); err != nil {
				log.Errorln(err)
			}
		}
	default:
		for _, receiver := range f.spec.Receivers {
			gv, kind := gvk.ToAPIVersionAndKind()
			sub := fmt.Sprintf("[%s]: %s %s %s/%s %s", stringz.Val(f.clusterName, "?"), gv, kind, accessor.GetNamespace(), accessor.GetName(), op)
			if err := f.notify(sub, sub, string(data), receiver); err != nil {
				log.Errorln(err)
			}
		}
	}
	return nil
}
