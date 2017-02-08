package k8s

import (
	_env "github.com/appscode/go/env"
	_ "github.com/appscode/k8s-addons/api/install"
	acs "github.com/appscode/k8s-addons/client/clientset"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	rest "k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
)

type KubeClient struct {
	Client                  clientset.Interface
	AppscodeExtensionClient acs.AppsCodeExtensionInterface
}

func GetKubeConfig() (config *rest.Config, err error) {
	debugEnabled := _env.FromHost().DebugEnabled()
	if !debugEnabled {
		config, err = rest.InClusterConfig()
	} else {
		rules := clientcmd.NewDefaultClientConfigLoadingRules()
		rules.DefaultClientConfig = &clientcmd.DefaultClientConfig
		overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}
		config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides).ClientConfig()
	}
	return
}
