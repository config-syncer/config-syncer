package watcher

import (
	"github.com/appscode/voyager/pkg/controller/certificates"
	aci "github.com/appscode/k8s-addons/api"
	"github.com/appscode/k8s-addons/pkg/events"
	acw "github.com/appscode/k8s-addons/pkg/watcher"
	"github.com/appscode/log"
	"k8s.io/kubernetes/pkg/client/cache"
	"k8s.io/kubernetes/pkg/util/wait"
)

func (k *KubedWatcher) Certificate() {
	log.Debugln("watching", events.Certificate.String())
	lw := &cache.ListWatch{
		ListFunc:  acw.CertificateListFunc(k.AppsCodeExtensionClient),
		WatchFunc: acw.CertificateWatchFunc(k.AppsCodeExtensionClient),
	}
	_, controller := k.Cache(events.Certificate, &aci.Certificate{}, lw)
	go controller.Run(wait.NeverStop)

	go certificates.NewCertificateSyncer(k.Client, k.AppsCodeExtensionClient).RunSync()
}
