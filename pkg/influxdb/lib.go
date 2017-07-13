package influx

import (
	"fmt"
	"net/url"

	"github.com/appscode/kubed/pkg/config"
	influxdb "github.com/influxdata/influxdb/client"
)

type Janitor struct {
	Config config.ClusterConfig
}

func (j *Janitor) CleanInflux() error {
	u, err := url.Parse(j.Config.InfluxDB.Endpoint)
	if err != nil {
		return err
	}

	client, err := influxdb.NewClient(&influxdb.Config{
		URL:      *u,
		Username: j.Config.InfluxDB.Username,
		Password: j.Config.InfluxDB.Password,
	})
	if err != nil {
		return err
	}
	return j.UpdateRetentionPolicy(client)
}

func (j *Janitor) UpdateRetentionPolicy(influxClient *influxdb.Client) error {
	query := influxdb.Query{
		Command:  fmt.Sprintf("ALTER RETENTION POLICY default ON k8s DURATION %vs", j.Config.InfluxDB.MonitoringStorageLifetime),
		Database: "k8s",
	}
	_, err := influxClient.Query(query)
	return err
}
