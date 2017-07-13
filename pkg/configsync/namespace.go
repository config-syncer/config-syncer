package configsync

import (
	"github.com/appscode/errors"
	"github.com/appscode/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientset "k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

const (
	ConfigMaps string = "configmaps"
	Secrets    string = "secrets"

	KubeSecretAPIToken    = "appscode-api-token"
	KubeSecretIcinga      = "searchlight-icinga"
	KubeConfigMapMetadata = "cluster-metadata"
)

type NamespaceHandler struct {
	KubeClient clientset.Interface
}

func NewHandler(k clientset.Interface) *NamespaceHandler {
	return &NamespaceHandler{
		KubeClient: k,
	}
}

func NewType(t string) (runtime.Object, error) {
	switch t {
	case Secrets:
		return &apiv1.Secret{}, nil
	case ConfigMaps:
		return &apiv1.ConfigMap{}, nil
	}
	return nil, errors.New("Resource type: " + t + " not found").Err()
}

func setObjectMeta(o interface{}, namespace string, t string) {
	var objectMeta *metav1.ObjectMeta
	switch t {
	case Secrets:
		objectMeta = &o.(*apiv1.Secret).ObjectMeta
	case ConfigMaps:
		objectMeta = &o.(*apiv1.ConfigMap).ObjectMeta
	}
	objectMeta.SetNamespace(namespace)
	objectMeta.SetResourceVersion("")
}

func (h *NamespaceHandler) Handle(n interface{}) {
	ns, ok := n.(*apiv1.Namespace)
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
	obj, err := NewType(t)
	if err != nil {
		log.Errorln(err)
		return false
	}
	err = h.KubeClient.CoreV1().RESTClient().Get().
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
	result, err := NewType(t)
	if err != nil {
		log.Errorln(err)
		return
	}
	err = h.KubeClient.CoreV1().RESTClient().Get().
		Namespace(metav1.NamespaceSystem).
		Resource(t).
		Name(name).
		Do().
		Into(result)
	if err != nil {
		log.Errorln(err)
		return
	}
	setObjectMeta(result, namespace, t)
	err = h.KubeClient.CoreV1().RESTClient().Post().
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
