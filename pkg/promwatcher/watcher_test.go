package watcher

import (
	"testing"
	"time"

	acs "github.com/appscode/k8s-addons/client/clientset"
	acw "github.com/appscode/k8s-addons/pkg/watcher"
	"github.com/appscode/kubed/pkg/watcher"
	pcm "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	restclient "k8s.io/kubernetes/pkg/client/restclient"
	unversioned_clientcmd "k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
)

func getKubeConfig() (config *restclient.Config, err error) {
	rules := unversioned_clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &unversioned_clientcmd.DefaultClientConfig
	overrides := &unversioned_clientcmd.ConfigOverrides{ClusterDefaults: unversioned_clientcmd.ClusterDefaults}
	config, err = unversioned_clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides).ClientConfig()
	return
}

func getClientGoConfig() (config *rest.Config, err error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}
	config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides).ClientConfig()
	return
}

func TestWatchPrometheus(t *testing.T) {
	// get client for Prometheus TPR monitoring

	kubeConfig, err := getKubeConfig()
	assert.Nil(t, err)

	kubeWatcher := &watcher.KubedWatcher{
		Watcher: acw.Watcher{
			Client:                  clientset.NewForConfigOrDie(kubeConfig),
			AppsCodeExtensionClient: acs.NewACExtensionsForConfigOrDie(kubeConfig),
			SyncPeriod:              time.Minute * 2,
		},
	}

	go kubeWatcher.Run()

	clientGoConfig, err := getClientGoConfig()
	assert.Nil(t, err)

	client, err := pcm.NewForConfig(clientGoConfig)
	assert.Nil(t, err)

	watcher := &PromWatcher{
		Watcher:    kubeWatcher.Watcher,
		PromClient: client,
		SyncPeriod: time.Minute * 2,
	}
	go watcher.WatchPrometheus()

	time.Sleep(time.Hour)
}
