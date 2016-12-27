package influxdb

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/appscode/searchlight/pkg/client/k8s"
	influxdb "github.com/influxdata/influxdb/client"
	ini "github.com/vaughan0/go-ini"
)

const (
	admin string = ".admin"

	influxDBHost     string = "INFLUX_HOST"
	influxDBDatabase string = "INFLUX_DB"
	influxDBReadUser string = "INFLUX_READ_USER"
	influxDBReadPass string = "INFLUX_READ_PASSWORD"

	influxDBDefaultDatabase string = "k8s"
	influxDBHostPort               = 8086
)

type AuthInfo struct {
	Host     string
	Username string
	Password string
	Database string
}

func GetInfluxDBSecretData(secretName string) (*AuthInfo, error) {
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}

	parts := strings.Split(secretName, ".")
	name := parts[0]
	namespace := "default"
	if len(parts) > 1 {
		namespace = parts[1]
	}

	secret, err := kubeClient.Client.Core().Secrets(namespace).Get(name)
	if err != nil {
		return nil, err
	}

	authData := new(AuthInfo)
	if data, found := secret.Data[admin]; found {
		dataReader := strings.NewReader(string(data))
		secretData, err := ini.Load(dataReader)
		if err != nil {
			return nil, err
		}

		if host, found := secretData.Get("", influxDBHost); found {
			authData.Host = fmt.Sprintf("%s:%d", host, influxDBHostPort)
		}

		if authData.Database, found = secretData.Get("", influxDBDatabase); !found {
			authData.Database = influxDBDefaultDatabase
		}
		if authData.Username, found = secretData.Get("", influxDBReadUser); !found {
			return nil, errors.New("No InfluxDB read user found")
		}
		if authData.Password, found = secretData.Get("", influxDBReadPass); !found {
			return nil, errors.New("No InfluxDB read password found")
		}
		return authData, nil
	}
	return nil, errors.New("Invalid InfluxDB secret")
}

func GetInfluxDBConfig(secretName string) (*influxdb.Config, error) {
	authData, err := GetInfluxDBSecretData(secretName)
	if err != nil {
		return nil, err
	}
	config := &influxdb.Config{
		URL: url.URL{
			Scheme: "http",
			Host:   authData.Host,
		},
		Username: authData.Username,
		Password: authData.Password,
	}
	return config, nil
}
