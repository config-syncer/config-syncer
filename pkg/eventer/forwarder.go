package eventer

import (
	"strings"
	"sync"

	api "github.com/appscode/kubed/apis/kubed/v1alpha1"
	"gomodules.xyz/envconfig"
	notify "gomodules.xyz/notify"
	"gomodules.xyz/notify/unified"
	"k8s.io/client-go/discovery"
	discovery_util "kmodules.xyz/client-go/discovery"
)

type EventForwarder struct {
	Client discovery.DiscoveryInterface

	clusterName  string
	spec         *api.EventForwarderSpec
	notifierCred envconfig.LoaderFunc
	restmapper   *discovery_util.DefaultRESTMapper

	lock sync.RWMutex
}

func (f *EventForwarder) Configure(clusterName string, spec *api.EventForwarderSpec, notifierCred envconfig.LoaderFunc) error {
	f.lock.Lock()
	defer f.lock.Unlock()

	var err error

	f.clusterName = clusterName
	f.spec = spec
	f.notifierCred = notifierCred
	f.restmapper, err = discovery_util.LoadRestMapper(f.Client)

	return err
}

func (f *EventForwarder) notify(emailSub, chatSub, body string, receiver api.Receiver) error {
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
