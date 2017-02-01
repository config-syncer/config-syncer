package icinga

import (
	"errors"
	"fmt"
	"strings"

	_env "github.com/appscode/go/env"
	_ "github.com/appscode/k8s-addons/api/install"
	"github.com/appscode/k8s-addons/pkg/dns"
	ini "github.com/vaughan0/go-ini"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
)

const (
	env string = ".env"

	IcingaService string = "ICINGA_K8S_SERVICE"
	IcingaAPIUser string = "ICINGA_API_USER"
	IcingaAPIPass string = "ICINGA_API_PASSWORD"
)

type authInfo struct {
	Endpoint string
	Username string
	Password string
}

func getIcingaSecretData(kubeClient clientset.Interface, secretName string) (*authInfo, error) {
	parts := strings.Split(secretName, ".")
	name := parts[0]
	namespace := "default"
	if len(parts) > 1 {
		namespace = parts[1]
	}

	secret, err := kubeClient.Core().Secrets(namespace).Get(name)
	if err != nil {
		return nil, err
	}

	authData := new(authInfo)
	if data, found := secret.Data[env]; found {
		dataReader := strings.NewReader(string(data))
		secretData, err := ini.Load(dataReader)
		if err != nil {
			return nil, err
		}

		if _env.InCluster() {
			if host, found := secretData.Get("", IcingaService); found {
				serviceIP, err := dns.GetServiceClusterIP(kubeClient, ConfigKeyPrefix, host)
				if err != nil {
					return nil, err
				}
				authData.Endpoint = fmt.Sprintf("https://%v:5665/v1", serviceIP)
			}
		}

		if authData.Username, found = secretData.Get("", IcingaAPIUser); !found {
			return nil, errors.New("No ICINGA_API_USER found")
		}

		if authData.Password, found = secretData.Get("", IcingaAPIPass); !found {
			return nil, errors.New("No ICINGA_API_PASSWORD found")
		}
		return authData, nil
	}
	return nil, errors.New("Invalid Icinga secret")
}

func getIcingaConfig(kubeClient clientset.Interface, secretName string) (*IcingaConfig, error) {
	authData, err := getIcingaSecretData(kubeClient, secretName)
	if err != nil {
		return nil, err
	}

	icingaConfig := &IcingaConfig{
		Endpoint: authData.Endpoint,
		CaCert:   nil,
	}
	icingaConfig.BasicAuth.Username = authData.Username
	icingaConfig.BasicAuth.Password = authData.Password

	return icingaConfig, nil
}
