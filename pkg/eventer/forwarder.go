package eventer

import (
	"strings"
	"sync"

	"github.com/appscode/envconfig"
	"github.com/appscode/go-notify"
	"github.com/appscode/go-notify/unified"
	apis "github.com/appscode/kubed/pkg/apis/v1alpha1"
	discovery_util "github.com/appscode/kutil/discovery"
	"k8s.io/client-go/discovery"
)

type EventForwarder struct {
	Client discovery.DiscoveryInterface

	clusterName  string
	spec         *apis.EventForwarderSpec
	notifierCred envconfig.LoaderFunc
	restmapper   *discovery_util.DefaultRESTMapper

	lock sync.RWMutex
}

func (f *EventForwarder) Configure(clusterName string, spec *apis.EventForwarderSpec, notifierCred envconfig.LoaderFunc) error {
	f.lock.Lock()
	defer f.lock.Unlock()

	var err error

	f.clusterName = clusterName
	f.spec = spec
	f.notifierCred = notifierCred
	f.restmapper, err = discovery_util.LoadRestMapper(f.Client)

	return err
}

func (f *EventForwarder) notify(emailSub, chatSub, body string, receiver apis.Receiver) error {
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
