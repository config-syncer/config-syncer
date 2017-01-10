package k8s

import (
	"log"

	"github.com/appscode/errors"
	_ "github.com/appscode/k8s-addons/api/install"
	acs "github.com/appscode/k8s-addons/client/clientset"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
)

// NewClient() should only be used to create kube client for plugins.
func NewClient() (*KubeClient, error) {

	config, err := GetKubeConfig()
	if err != nil {
		return nil, errors.New().WithCause(err).Internal()
	}
	log.Println("Using cluster:", config.Host)

	client, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, errors.New().WithCause(err).Internal()
	}

	appscodeClient, err := acs.NewACExtensionsForConfig(config)
	if err != nil {
		return nil, errors.New().WithCause(err).Internal()
	}

	return &KubeClient{
		Client:                  client,
		AppscodeExtensionClient: appscodeClient,
	}, nil
}
