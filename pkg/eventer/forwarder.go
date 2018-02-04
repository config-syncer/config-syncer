package eventer

import (
	"strings"
	"sync"

	"github.com/appscode/envconfig"
	"github.com/appscode/go-notify"
	"github.com/appscode/go-notify/unified"
	"github.com/appscode/kubed/pkg/config"
	"github.com/appscode/kutil/discovery"
)

type EventForwarder struct {
	clusterName  string
	spec         *config.EventForwarderSpec
	notifierCred envconfig.LoaderFunc
	restmapper   *discovery.DefaultRESTMapper

	lock sync.RWMutex
}

func (f *EventForwarder) Configure(clusterName string, spec *config.EventForwarderSpec, notifierCred envconfig.LoaderFunc) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.clusterName = clusterName
	f.spec = spec
	f.notifierCred = notifierCred
}

func (f *EventForwarder) notify(emailSub, chatSub, body string, receiver config.Receiver) error {
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
