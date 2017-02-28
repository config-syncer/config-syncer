package watcher

import (
	"time"

	"github.com/appscode/go/wait"
	acw "github.com/appscode/k8s-addons/pkg/watcher"
	"github.com/appscode/log"
	pcm "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1alpha1"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/tools/cache"
)

type PromWatcher struct {
	acw.Watcher

	// client for prometheus monitoring
	PromClient *pcm.MonitoringV1alpha1Client

	// sync time to sync the list.
	SyncPeriod time.Duration
}

func (p *PromWatcher) WatchPrometheus() {
	handler := cache.ResourceEventHandlerFuncs{
		AddFunc:    p.add,
		DeleteFunc: p.delete,
	}

	lw := &cache.ListWatch{
		ListFunc:  p.PromClient.Prometheuses(api.NamespaceAll).List,
		WatchFunc: p.PromClient.Prometheuses(api.NamespaceAll).Watch,
	}

	_, controller := cache.NewInformer(lw, &pcm.Prometheus{}, p.SyncPeriod, handler)
	go controller.Run(wait.NeverStop)
}

const keyPrometheus string = "prometheus"

func (p *PromWatcher) add(obj interface{}) {
	prometheus := obj.(*pcm.Prometheus)
	namespace := prometheus.Namespace

	backendService, err := p.createGoverningService(prometheus)
	if err != nil {
		log.Errorln(err)
		return
	}

	deployment, err := p.createProxyDeployment(prometheus.Name, backendService.Name, namespace)
	if err != nil {
		log.Errorln(err)
		return
	}

	if err := p.createProxyService(prometheus.Name, deployment); err != nil {
		log.Errorln(err)
		return
	}

	if err := p.createIngressRule(prometheus.Name, namespace); err != nil {
		log.Errorln(err)
		return
	}
}

func (p *PromWatcher) delete(obj interface{}) {
	prometheus := obj.(*pcm.Prometheus)
	namespace := prometheus.Namespace

	if err := p.deleteProxyDeployment(prometheus.Name, namespace); err != nil {
		log.Errorln(err)
		return
	}

	if err := p.deleteProxyService(prometheus.Name, namespace); err != nil {
		log.Errorln(err)
		return
	}

	if err := p.deleteIngressRule(prometheus.Name, namespace); err != nil {
		log.Errorln(err)
		return
	}
}
