package watcher

import (
	"reflect"

	"appscode.com/kubed/pkg/handlers"
	"github.com/appscode/client"
	"github.com/appscode/k8s-addons/pkg/events"
	"github.com/appscode/k8s-addons/pkg/stash"
	acw "github.com/appscode/k8s-addons/pkg/watcher"
	"github.com/appscode/log"
	"github.com/appscode/searchlight/pkg/client/icinga"
	aac "github.com/appscode/searchlight/pkg/controller"
	"github.com/appscode/voyager/pkg/controller/certificates"
	lbc "github.com/appscode/voyager/pkg/controller/ingress"
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
	watch.Pod()
	watch.StatefulSet()
	watch.DaemonSet()
	watch.ReplicaSet()
	watch.Namespace()
	watch.Node()
	watch.Service()
	watch.RC()
	watch.Endpoint()

	watch.ExtendedIngress()
	watch.Ingress()
	watch.Alert()
	watch.Certificate()
}

func (k *KubedWatcher) setup() {
	lbc.SetLoadbalancerImage(k.LoadbalancerImage)
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

	if e.ResourceType == events.Ingress || e.ResourceType == events.ExtendedIngress {
		// Handle Ingress first
		err := lbc.NewEngressController(k.ClusterName,
			k.ProviderName,
			k.Client,
			k.AppsCodeExtensionClient,
			k.Storage, k.IngressClass).Handle(e)

		// Check the Ingress or Extended Ingress Annotations. To Work for auto certificate
		// operations.
		if err == nil {
			certController := certificates.NewController(k.Client, k.AppsCodeExtensionClient)
			certController.Handle(e)
		}
		return err
	}

	if e.ResourceType == events.Certificate {
		if e.EventType.IsAdded() || e.EventType.IsUpdated() {
			certController := certificates.NewController(k.Client, k.AppsCodeExtensionClient)
			certController.Handle(e)
		}
	}

	if e.ResourceType == events.Service {
		if e.EventType.IsAdded() || e.EventType.IsDeleted() {
			return lbc.UpgradeAllEngress(e.MetaData.Name+"."+e.MetaData.Namespace,
				k.ClusterName,
				k.ProviderName,
				k.Client,
				k.AppsCodeExtensionClient,
				k.Storage, k.IngressClass)
		}
	}

	if e.ResourceType == events.Endpoint {
		// Checking if this endpoint have a service or not. If
		// this do not have a Service we do not want to update our ingress
		_, err := k.Client.Core().Services(e.MetaData.Namespace).Get(e.MetaData.Name)
		if err == nil {
			log.Infoln("Endpoint has an service with name", e.MetaData.Name, e.MetaData.Namespace, "Event type", e.EventType.String())
			// Service exists. So we should process.
			if e.EventType.IsUpdated() {
				return lbc.UpgradeAllEngress(e.MetaData.Name+"."+e.MetaData.Namespace,
					k.ClusterName,
					k.ProviderName,
					k.Client,
					k.AppsCodeExtensionClient,
					k.Storage, k.IngressClass)
			}
		}
	}

	if e.ResourceType == events.Alert || e.ResourceType == events.Node || e.ResourceType == events.Pod || e.ResourceType == events.Service {
		return aac.New(k.Client, k.IcingaClient, k.AppsCodeExtensionClient, k.Storage).Handle(e)
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
