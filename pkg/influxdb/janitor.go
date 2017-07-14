package influx

import (
	"fmt"
	"net/url"

	"github.com/appscode/kubed/pkg/config"
	influxdb "github.com/influxdata/influxdb/client"
)

type Janitor struct {
	Spec config.InfluxDBSpec
}

func (j *Janitor) Cleanup() error {
	u, err := url.Parse(j.Spec.Endpoint)
	if err != nil {
		return err
	}

	client, err := influxdb.NewClient(influxdb.Config{
		URL:      *u,
		Username: j.Spec.Username,
		Password: j.Spec.Password,
	})
	if err != nil {
		return err
	}

	query := influxdb.Query{
		Command:  fmt.Sprintf("ALTER RETENTION POLICY default ON k8s DURATION %vs", int(j.Spec.TTL.Seconds())),
		Database: "k8s",
	}
	_, err = client.Query(query)
	return err
}
