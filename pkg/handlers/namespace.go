package handlers

import (
	"github.com/appscode/errors"
	"github.com/appscode/kubed/pkg/events"
	"github.com/appscode/log"
	"k8s.io/kubernetes/pkg/api"
	kapi "k8s.io/kubernetes/pkg/api"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/runtime"
)

const (
	ConfigMaps string = "configmaps"
	Secrets    string = "secrets"

	KubeSecretAPIToken    = "appscode-api-token"
	KubeSecretIcinga      = "appscode-icinga"
	KubeConfigMapMetadata = "cluster-metadata"
)

type NamespaceHandler struct {
	KubeClient clientset.Interface
}

func New(t string) (runtime.Object, error) {
	switch t {
	case Secrets:
		return &kapi.Secret{}, nil
	case ConfigMaps:
		return &kapi.ConfigMap{}, nil
	}
	return nil, errors.New("Resource type: " + t + " not found").Err()
}

func setObjectMeta(o interface{}, namespace string, t string) {
	var objectMeta *kapi.ObjectMeta
	switch t {
	case Secrets:
		objectMeta = &o.(*kapi.Secret).ObjectMeta
	case ConfigMaps:
		objectMeta = &o.(*kapi.ConfigMap).ObjectMeta
	}
	objectMeta.SetNamespace(namespace)
	objectMeta.SetResourceVersion("")
}

func (h *NamespaceHandler) Handle(e *events.Event) {
	if !e.EventType.IsAdded() {
		return
	}
	ns, ok := e.RuntimeObj[0].(*api.Namespace)
	if ok {
		h.ensureTypes(ns.Name)
	}
}

func (h *NamespaceHandler) ensureTypes(namespace string) {
	if !h.isFound(namespace, Secrets, KubeSecretAPIToken) {
		h.copyObjectFromKubeSystemNS(namespace, Secrets, KubeSecretAPIToken)
	}
	if !h.isFound(namespace, Secrets, KubeSecretIcinga) {
		h.copyObjectFromKubeSystemNS(namespace, Secrets, KubeSecretIcinga)
	}
	if !h.isFound(namespace, ConfigMaps, KubeConfigMapMetadata) {
		h.copyObjectFromKubeSystemNS(namespace, ConfigMaps, KubeConfigMapMetadata)
	}
}

func (h *NamespaceHandler) isFound(namespace string, t string, name string) bool {
	var err error
	obj, err := New(t)
	if err != nil {
		log.Errorln(err)
		return false
	}
	err = h.KubeClient.Core().RESTClient().Get().
		Namespace(namespace).
		Resource(t).
		Name(name).
		Do().Into(obj)
	if err != nil {
		return false
	}
	return true
}

func (h *NamespaceHandler) copyObjectFromKubeSystemNS(namespace string, t string, name string) {
	result, err := New(t)
	if err != nil {
		log.Errorln(err)
		return
	}
	err = h.KubeClient.Core().RESTClient().Get().
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
	err = h.KubeClient.Core().RESTClient().Post().
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
