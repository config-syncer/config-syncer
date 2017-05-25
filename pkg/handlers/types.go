package handlers

import (
	"github.com/appscode/client"
	"github.com/appscode/errors"
	"github.com/appscode/k8s-addons/pkg/stash"
	kapi "k8s.io/kubernetes/pkg/api"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/runtime"
)

const (
	ConfigMaps string = "configmaps"
	Secrets    string = "secrets"
)

type Handler struct {
	ClientOptions *client.ClientOption
	ClusterName   string
	Kube          clientset.Interface
	Storage       *stash.Storage
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
