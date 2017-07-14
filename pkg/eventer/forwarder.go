package eventer

import (
	"fmt"

	"github.com/appscode/envconfig"
	"github.com/appscode/go-notify"
	"github.com/appscode/go-notify/unified"
	stringz "github.com/appscode/go/strings"
	"github.com/appscode/kubed/pkg/config"
	"github.com/ghodss/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

type EventForwarder struct {
	Spec   config.EventForwarderSpec
	Loader envconfig.LoaderFunc
}

func (f *EventForwarder) ForwardEvent(e *apiv1.Event) error {
	if e.Type == apiv1.EventTypeWarning &&
		(len(f.Spec.EventNamespaces) == 0 || stringz.Contains(f.Spec.EventNamespaces, e.Namespace)) {

		if f.Spec.NotifyVia != "" {
			sub := fmt.Sprintf("%s %s/%s: %s", e.InvolvedObject.Kind, e.InvolvedObject.Namespace, e.InvolvedObject.Name, e.Reason)
			if notifier, err := unified.LoadVia(f.Spec.NotifyVia, f.Loader); err == nil {
				switch n := notifier.(type) {
				case notify.ByEmail:
					bytes, err := yaml.Marshal(e)
					if err != nil {
						return err
					}
					n.WithSubject(sub).WithBody(string(bytes)).Send()
				case notify.BySMS:
					n.WithBody(sub).Send()
				case notify.ByChat:
					n.WithBody(sub).Send()
				}
			}
		}
	}
	return nil
}

func (f *EventForwarder) Forward(t metav1.TypeMeta, meta metav1.ObjectMeta, v interface{}) error {
	bytes, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	if f.Spec.NotifyVia != "" {
		sub := fmt.Sprintf("%s %s %s/%s added", t.APIVersion, t.Kind, meta.Namespace, meta.Name)
		if notifier, err := unified.LoadVia(f.Spec.NotifyVia, f.Loader); err == nil {
			switch n := notifier.(type) {
			case notify.ByEmail:
				n.WithSubject(sub).WithBody(string(bytes)).Send()
			case notify.BySMS:
				n.WithBody(sub).Send()
			case notify.ByChat:
				n.WithBody(sub).Send()
			}
		}
	}
	return nil
}
