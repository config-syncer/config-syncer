package k8s

import (
	_ "github.com/appscode/k8s-addons/api/install"
	acs "github.com/appscode/k8s-addons/client/clientset"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	rest "k8s.io/kubernetes/pkg/client/restclient"
)

type KubeClient struct {
	Client                  clientset.Interface
	AppscodeExtensionClient acs.AppsCodeExtensionInterface

	config *rest.Config
}
