package util

import (
	"github.com/appscode/errors"
	"github.com/appscode/log"
	_ "github.com/appscode/searchlight/api/install"
	acs "github.com/appscode/searchlight/client/clientset"
	clientset "k8s.io/client-go/kubernetes"
)

// NewClient() should only be used to create kube client for plugins.
func NewClient() (*KubeClient, error) {

	config, err := GetKubeConfig()
	if err != nil {
		return nil, errors.FromErr(err).Err()
	}
	log.Debugln("Using cluster:", config.Host)

	client, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, errors.FromErr(err).Err()
	}

	extClient, err := acs.NewForConfig(config)
	if err != nil {
		return nil, errors.FromErr(err).Err()
	}

	return &KubeClient{
		Client:    client,
		ExtClient: extClient,
	}, nil
}
