package eventer

import (
	"fmt"
	"strings"

	"github.com/appscode/envconfig"
	"github.com/appscode/go-notify"
	"github.com/appscode/go-notify/unified"
	stringz "github.com/appscode/go/strings"
	"github.com/appscode/kubed/pkg/config"
	"github.com/appscode/log"
	"github.com/ghodss/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

type EventForwarder struct {
	ClusterName string
	Receivers   []config.Receiver
	Loader      envconfig.LoaderFunc
}

func (f *EventForwarder) ForwardEvent(e *apiv1.Event) error {
	bytes, err := yaml.Marshal(e)
	if err != nil {
		return err
	}
	host := ""
	if e.Source.Host != "" {
		host = "on host " + e.Source.Host
	}
	for _, receiver := range f.Receivers {
		emailSub := fmt.Sprintf("[%s, %s]: %s %s/%s %s %s", stringz.Val(f.ClusterName, "?"), e.Source.Component, e.InvolvedObject.Kind, e.InvolvedObject.Namespace, e.InvolvedObject.Name, e.Reason, host)
		chatSub := fmt.Sprintf("[%s, %s] %s %s/%s %s %s: %s", stringz.Val(f.ClusterName, "?"), e.Source.Component, e.InvolvedObject.Kind, e.InvolvedObject.Namespace, e.InvolvedObject.Name, e.Reason, host, e.Message)
		if err := f.send(emailSub, chatSub, string(bytes), receiver); err != nil {
			log.Errorln(err)
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
		sub := fmt.Sprintf("[%s]: %s %s %s/%s added", stringz.Val(f.ClusterName, "?"), t.APIVersion, t.Kind, meta.Namespace, meta.Name)
		if err := f.send(sub, sub, string(bytes), receiver); err != nil {
			log.Errorln(err)
		}
	}
	return nil
}

func (f *EventForwarder) send(emailSub, chatSub, body string, receiver config.Receiver) error {
	notifier, err := unified.LoadVia(strings.ToLower(receiver.Notifier), f.Loader)
	if err != nil {
		return err
	}
	switch n := notifier.(type) {
	case notify.ByEmail:
		return n.To(receiver.To[0], receiver.To[1:]...).
			WithSubject(emailSub).
			WithBody(body).
			WithNoTracking().
			Send()
	case notify.BySMS:
		return n.To(receiver.To[0], receiver.To[1:]...).
			WithBody(emailSub).
			Send()
	case notify.ByChat:
		return n.To(receiver.To[0], receiver.To[1:]...).
			WithBody(chatSub).
			Send()
	}
	return nil
}
