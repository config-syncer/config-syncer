package influxdb

import (
	"fmt"
	"net/url"

	"github.com/appscode/go/io"
	"github.com/appscode/k8s-addons/pkg/dns"
	"github.com/influxdata/influxdb/client"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
)

const (
	ConfigFile      = "/srv/influxdb/secrets/.admin"
	ConfigKeyPrefix = "INFLUX"
	InfluxDBHost    = "INFLUX_HOST"
	InfluxDBPort    = "INFLUX_API_PORT"
	InfluxDBUser    = "INFLUX_WRITE_USER"
	InfluxDBPass    = "INFLUX_WRITE_PASSWORD"
)

func LoadConfig(kubeClient clientset.Interface) (*client.Config, error) {
	m, err := io.ReadINIConfig(ConfigFile)
	if err != nil {
		return nil, err
	}
	serviceIP, err := dns.GetServiceClusterIP(kubeClient, ConfigKeyPrefix, m[InfluxDBHost])
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(fmt.Sprintf("http://%v:%v", serviceIP, m[InfluxDBPort]))
	if err != nil {
		return nil, err
	}
	return &client.Config{
		URL:       *u,
		Username:  m[InfluxDBUser],
		Password:  m[InfluxDBPass],
		UserAgent: fmt.Sprintf("%v/%v", "searchlight", 1.0),
	}, nil
}
