package eventer

import (
	"fmt"
	"strings"

	"github.com/appscode/envconfig"
	"github.com/appscode/go-notify"
	"github.com/appscode/go-notify/unified"
	"github.com/appscode/kubed/pkg/config"
	"github.com/ghodss/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"github.com/tamalsaha/go-oneliners"
)

type EventForwarder struct {
	Receivers []config.Receiver
	Loader    envconfig.LoaderFunc
}

func (f *EventForwarder) ForwardEvent(e *apiv1.Event) error {
	oneliners.FILE()
	if e.Type == apiv1.EventTypeWarning {
		oneliners.FILE()
		for _, receiver := range f.Receivers {
			oneliners.FILE()
			if len(receiver.To) > 0 {
				oneliners.FILE()
				sub := fmt.Sprintf("%s %s/%s: %s", e.InvolvedObject.Kind, e.InvolvedObject.Namespace, e.InvolvedObject.Name, e.Reason)
				if notifier, err := unified.LoadVia(strings.ToLower(receiver.Notifier), f.Loader); err == nil {
					switch n := notifier.(type) {
					case notify.ByEmail:
						bytes, err := yaml.Marshal(e)
						if err != nil {
							oneliners.FILE(err)
							return err
						}
						n.To(receiver.To[0], receiver.To[1:]...).
							WithSubject(sub).
							WithBody(string(bytes)).
							WithNoTracking().
							Send()
					case notify.BySMS:
						n.To(receiver.To[0], receiver.To[1:]...).
							WithBody(sub).
							Send()
					case notify.ByChat:
						n.To(receiver.To[0], receiver.To[1:]...).
							WithBody(sub).
							Send()
					}
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
	for _, receiver := range f.Receivers {
		if len(receiver.To) > 0 {
			sub := fmt.Sprintf("%s %s %s/%s added", t.APIVersion, t.Kind, meta.Namespace, meta.Name)
			if notifier, err := unified.LoadVia(strings.ToLower(receiver.Notifier), f.Loader); err == nil {
				switch n := notifier.(type) {
				case notify.ByEmail:
					n.To(receiver.To[0], receiver.To[1:]...).
						WithSubject(sub).
						WithBody(string(bytes)).
						WithNoTracking().
						Send()
				case notify.BySMS:
					n.To(receiver.To[0], receiver.To[1:]...).
						WithBody(sub).
						Send()
				case notify.ByChat:
					n.To(receiver.To[0], receiver.To[1:]...).
						WithBody(sub).
						Send()
				}
			}
		}
	}
	return nil
}
