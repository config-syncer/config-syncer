package icinga

import (
	"fmt"

	"github.com/appscode/go/io"
	_ "github.com/appscode/k8s-addons/api/install"
	"github.com/appscode/k8s-addons/pkg/dns"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
)

const (
	ConfigFile             = "/srv/icinga2/secrets/.env"
	ConfigKeyPrefix        = "ICINGA"
	IcingaService   string = "ICINGA_K8S_SERVICE"
	IcingaAPIUser   string = "ICINGA_API_USER"
	IcingaAPIPass   string = "ICINGA_API_PASSWORD"
)

func NewInClusterClient(kubeClient clientset.Interface) (*IcingaClient, error) {
	m, err := io.ReadINIConfig(ConfigFile)
	if err != nil {
		return nil, err
	}

	serviceName := m[IcingaService]
	if serviceName == "" {
		serviceName = "appscode-alert"
	}
	serviceIP, err := dns.GetServiceClusterIP(kubeClient, ConfigKeyPrefix, serviceName)
	if err != nil {
		return nil, err
	}
	c := NewClient(fmt.Sprintf("https://%v:5665/v1", serviceIP), m[IcingaAPIUser], m[IcingaAPIPass], nil)
	return c, nil
}
