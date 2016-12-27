package k8s

import (
	"os"

	"github.com/appscode/errors"
	_env "github.com/appscode/go/env"
	"github.com/appscode/go/io"
	_ "github.com/appscode/k8s-addons/api/install"
	acs "github.com/appscode/k8s-addons/client/clientset"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	rest "k8s.io/kubernetes/pkg/client/restclient"
)

var configDataPath = os.ExpandEnv("$GOPATH") + "/src/github.com/appscode/searchlight/pkg/client/k8s/config.ini"

const (
	host     = "host"
	username = "username"
	password = "password"
)

func NewClient() (*KubeClient, error) {
	var config *rest.Config
	var err error

	debugEnabled := _env.FromHost().DebugEnabled()
	if !debugEnabled {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		configData, err := io.ReadINIConfig(configDataPath)
		if err != nil {
			return nil, err
		}

		config = &rest.Config{
			Insecure: true,
		}
		if host, found := configData[host]; found {
			config.Host = host
		} else {
			return nil, errors.New().WithMessage("host address not fount").BadRequest()
		}

		if username, found := configData[username]; found {
			config.Username = username
		} else {
			return nil, errors.New().WithMessage("username not fount").BadRequest()
		}

		if password, found := configData[password]; found {
			config.Password = password
		} else {
			return nil, errors.New().WithMessage("password not fount").BadRequest()
		}
	}

	client, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, errors.New().WithCause(err).Internal()
	}

	appscodeClient, err := acs.NewACExtensionsForConfig(config)
	if err != nil {
		return nil, errors.New().WithCause(err).Internal()
	}

	return &KubeClient{
		config:                  config,
		Client:                  client,
		AppscodeExtensionClient: appscodeClient,
	}, nil
}
