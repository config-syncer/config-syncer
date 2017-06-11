package icinga

import (
	"errors"
	"fmt"
	"net"
	"strings"

	_env "github.com/appscode/go/env"
	_ "github.com/appscode/searchlight/api/install"
	"github.com/appscode/searchlight/pkg/dns"
	ini "github.com/vaughan0/go-ini"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
)

const (
	ENV string = ".env"

	IcingaAddress string = "ICINGA_ADDRESS"
	IcingaAPIUser string = "ICINGA_API_USER"
	IcingaAPIPass string = "ICINGA_API_PASSWORD"

	ConfigKeyPrefix = "ICINGA"

	IcingaDefaultPort string = "5665"
)

type authInfo struct {
	Endpoint string
	Username string
	Password string
}

func getIcingaSecretData(kubeClient clientset.Interface, secretName, secretNamespace string) (*authInfo, error) {
	secret, err := kubeClient.Core().Secrets(secretNamespace).Get(secretName)
	if err != nil {
		return nil, err
	}

	authData := new(authInfo)
	if data, found := secret.Data[ENV]; found {
		dataReader := strings.NewReader(string(data))
		secretData, err := ini.Load(dataReader)
		if err != nil {
			return nil, err
		}

		address, found := secretData.Get("", IcingaAddress)
		if !found {
			return nil, errors.New("No ICINGA_ADDRESS found")
		}

		parts := strings.Split(address, ":")
		host := parts[0]
		port := IcingaDefaultPort
		if len(parts) > 1 {
			port = parts[1]
		}

		hostIP := net.ParseIP(host)
		if hostIP == nil {
			if _env.InCluster() {
				host, err = dns.GetServiceClusterIP(kubeClient, ConfigKeyPrefix, host)
				if err != nil {
					return nil, err
				}
			}
		}

		authData.Endpoint = fmt.Sprintf("https://%v:%v/v1", host, port)

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

func getIcingaConfig(kubeClient clientset.Interface, secretName, secretNamespace string) (*IcingaConfig, error) {
	authData, err := getIcingaSecretData(kubeClient, secretName, secretNamespace)
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
