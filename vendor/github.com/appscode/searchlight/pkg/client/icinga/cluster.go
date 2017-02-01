package icinga

import (
	"fmt"
	"os"
	"strings"

	_ "github.com/appscode/k8s-addons/api/install"
	ini "github.com/vaughan0/go-ini"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
)

const (
	ConfigKeyPrefix = "ICINGA"
	IcingaSecret    = "E2E_ICINGA_SECRET"
	IcingaURL       = "E2E_ICINGA_URL"
)

func NewInClusterIcingaClient(kubeClient clientset.Interface, secretName string) (*IcingaClient, error) {
	config, err := getIcingaConfig(kubeClient, secretName)
	if err != nil {
		return nil, err
	}
	c := NewClient(config)
	return c, nil
}

func NewIcingaClient(kubeClient clientset.Interface) (*IcingaClient, error) {
	secretName := os.Getenv(IcingaSecret)

	if secretName == "" {
		fmt.Println("Set E2E_ICINGA_SECRET ENV to kubernetes secret name")
		os.Exit(1)
	}
	config, _ := getIcingaConfig(kubeClient, secretName)

	icinga_url := os.Getenv(IcingaURL)
	if icinga_url == "" {
		fmt.Println("Getting Icinga2 API URL from LoadBalancer.Ingress")
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
		if data, found := secret.Data[env]; found {
			dataReader := strings.NewReader(string(data))
			secretData, err := ini.Load(dataReader)
			if err != nil {
				return nil, err
			}
			if serviceName, found := secretData.Get("", IcingaService); found {
				parts := strings.Split(serviceName, ".")
				name := parts[0]
				namespace := "default"
				if len(parts) > 1 {
					namespace = parts[1]
				}

				service, err := kubeClient.Core().Services(namespace).Get(name)
				if err != nil {
					return nil, err
				}

				if len(service.Status.LoadBalancer.Ingress) > 0 {
					icinga_url = service.Status.LoadBalancer.Ingress[0].Hostname
				}
			}
		}

	}

	if icinga_url == "" {
		fmt.Println("Set E2E_ICINGA_URL ENV to Icinga API address")
		os.Exit(1)
	}

	config.Endpoint = fmt.Sprintf("https://%v:5665/v1", icinga_url)
	fmt.Println("Using Icinga2: ", config.Endpoint)
	c := NewClient(config)
	return c, nil
}
