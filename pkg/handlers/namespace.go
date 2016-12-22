package handlers

import (
	"github.com/appscode/k8s-addons/pkg/events"
	"github.com/appscode/log"
	"k8s.io/kubernetes/pkg/api"
)

const (
	KubeSecretAPIToken    = "appscode-api-token"
	KubeSecretIcinga      = "appscode-icinga"
	KubeConfigMapMetadata = "cluster-metadata"
)

type NamespaceHandler struct {
	*Handler
}

func (a *NamespaceHandler) Handle(e *events.Event) {
	if !e.EventType.IsAdded() {
		return
	}
	ns, ok := e.RuntimeObj[0].(*api.Namespace)
	if ok {
		a.ensureTypes(ns.Namespace)
	}
}

func (a *NamespaceHandler) ensureTypes(namespace string) {
	if !a.isFound(namespace, Secrets, KubeSecretAPIToken) {
		a.copyObjectFromKubeSystemNS(namespace, Secrets, KubeSecretAPIToken)
	}
	if !a.isFound(namespace, Secrets, KubeSecretIcinga) {
		a.copyObjectFromKubeSystemNS(namespace, Secrets, KubeSecretIcinga)
	}
	if !a.isFound(namespace, ConfigMaps, KubeConfigMapMetadata) {
		a.copyObjectFromKubeSystemNS(namespace, ConfigMaps, KubeConfigMapMetadata)
	}
}

func (a *NamespaceHandler) isFound(namespace string, t string, name string) bool {
	var err error
	obj, err := New(t)
	if err != nil {
		log.Errorln(err)
		return false
	}
	err = a.Kube.Core().RESTClient().Get().
		Namespace(namespace).
		Resource(t).
		Name(name).
		Do().Into(obj)
	if err != nil {
		return false
	}
	return true
}

func (a *NamespaceHandler) copyObjectFromKubeSystemNS(namespace string, t string, name string) {
	result, err := New(t)
	if err != nil {
		log.Errorln(err)
		return
	}
	err = a.Kube.Core().RESTClient().Get().
		Namespace(api.NamespaceSystem).
		Resource(t).
		Name(name).
		Do().
		Into(result)
	if err != nil {
		log.Errorln(err)
		return
	}
	setObjectMeta(result, namespace, t)
	err = a.Kube.Core().RESTClient().Post().
		Namespace(namespace).
		Resource(t).
		Body(result).
		Do().
		Into(result)

	if err != nil {
		log.Errorln(err)
		return
	} else {
		log.Infof("%s `%s` copied to namespace `%s` from kube-system.", t, name, namespace)
	}
}
