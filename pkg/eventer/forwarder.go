package eventer

import (
	"fmt"
	"strings"
	"sync"

	"github.com/appscode/envconfig"
	"github.com/appscode/go-notify"
	"github.com/appscode/go-notify/unified"
	"github.com/appscode/go/log"
	stringz "github.com/appscode/go/strings"
	"github.com/appscode/kubed/pkg/config"
	"github.com/ghodss/yaml"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type EventForwarder struct {
	clusterName  string
	spec         *config.EventForwarderSpec
	notifierCred envconfig.LoaderFunc

	lock sync.RWMutex
}

func (f *EventForwarder) Configure(clusterName string, spec *config.EventForwarderSpec, notifierCred envconfig.LoaderFunc) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.clusterName = clusterName
	f.spec = spec
	f.notifierCred = notifierCred
}

func (f *EventForwarder) ForwardEvent(e *core.Event) error {
	bytes, err := yaml.Marshal(e)
	if err != nil {
		return err
	}
	host := ""
	if e.Source.Host != "" {
		host = "on host " + e.Source.Host
	}
	for _, receiver := range f.spec.Receivers {
		emailSub := fmt.Sprintf("[%s, %s]: %s %s/%s %s %s", stringz.Val(f.clusterName, "?"), e.Source.Component, e.InvolvedObject.Kind, e.InvolvedObject.Namespace, e.InvolvedObject.Name, e.Reason, host)
		chatSub := fmt.Sprintf("[%s, %s] %s %s/%s %s %s: %s", stringz.Val(f.clusterName, "?"), e.Source.Component, e.InvolvedObject.Kind, e.InvolvedObject.Namespace, e.InvolvedObject.Name, e.Reason, host, e.Message)
		if err := f.send(emailSub, chatSub, string(bytes), receiver); err != nil {
			log.Errorln(err)
		}
	}
	return nil
}

func (f *EventForwarder) Forward(t metav1.TypeMeta, meta metav1.ObjectMeta, eventType string, v interface{}) error {
	bytes, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	for _, receiver := range f.spec.Receivers {
		sub := fmt.Sprintf("[%s]: %s %s %s/%s %s", stringz.Val(f.clusterName, "?"), t.APIVersion, t.Kind, meta.Namespace, meta.Name, eventType)
		if err := f.send(sub, sub, string(bytes), receiver); err != nil {
			log.Errorln(err)
		}
	}
	return nil
}

func (f *EventForwarder) send(emailSub, chatSub, body string, receiver config.Receiver) error {
	notifier, err := unified.LoadVia(strings.ToLower(receiver.Notifier), f.notifierCred)
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
	case notify.ByPush:
		return n.To(receiver.To...).
			WithBody(chatSub).
			Send()
	}
	return nil
}
