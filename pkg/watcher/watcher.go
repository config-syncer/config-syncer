package watcher

import (
	"reflect"

	"github.com/appscode/client"
	"github.com/appscode/k8s-addons/pkg/events"
	"github.com/appscode/k8s-addons/pkg/stash"
	acw "github.com/appscode/k8s-addons/pkg/watcher"
	"github.com/appscode/kubed/pkg/handlers"
	"github.com/appscode/log"
	"github.com/appscode/searchlight/pkg/client/icinga"
)

type KubedWatcher struct {
	acw.Watcher

	// name of the cloud provider
	ProviderName string

	// name of the cluster the daemon running.
	ClusterName string

	// appscode api server client
	AppsCodeApiClientOptions *client.ClientOption

	// Icinga Client
	IcingaClient *icinga.IcingaClient

	// Loadbalancer image name that will be used to create
	// the loadbalancer.
	LoadbalancerImage string

	IngressClass string
}

func (watch *KubedWatcher) Run() {
	watch.setup()
	watch.Storage = &stash.Storage{}
	watch.Namespace()
}

func (k *KubedWatcher) setup() {
	k.Watcher.Dispatch = k.Dispatch
}

func (k *KubedWatcher) Dispatch(e *events.Event) error {
	if ignoreAble(e) {
		return nil
	}
	log.Debugln("Dispatching event with resource", e.ResourceType, "event", e.EventType)
	if e.ResourceType == events.Namespace && e.EventType == events.Added {
		h := &handlers.NamespaceHandler{
			Handler: &handlers.Handler{
				ClientOptions: k.AppsCodeApiClientOptions,
				ClusterName:   k.ClusterName,
				Kube:          k.Client,
				Storage:       k.Storage,
			},
		}
		h.Handle(e)
	}
	return nil
}

func ignoreAble(e *events.Event) bool {
	if e.EventType == events.None {
		return true
	}

	if e.EventType == events.Updated {
		// updated called but only old object is present.
		if len(e.RuntimeObj) <= 1 {
			return true
		}

		// updated but both are equal. no changes
		if reflect.DeepEqual(e.RuntimeObj[0], e.RuntimeObj[1]) {
			return true
		}
	}
	return false
}
