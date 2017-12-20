package clientcmd

import (
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func BuildConfigFromContext(kubeconfigPath, contextName string) (*rest.Config, error) {
	var loader clientcmd.ClientConfigLoader
	if kubeconfigPath == "" {
		rules := clientcmd.NewDefaultClientConfigLoadingRules()
		rules.DefaultClientConfig = &clientcmd.DefaultClientConfig
		loader = rules
	} else {
		loader = &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath}
	}
	overrides := &clientcmd.ConfigOverrides{
		CurrentContext:  contextName,
		ClusterDefaults: clientcmd.ClusterDefaults,
	}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loader, overrides).ClientConfig()
}
