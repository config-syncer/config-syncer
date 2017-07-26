package influx

import (
	"fmt"
	"math"
	"net/url"
	"time"

	"github.com/appscode/kubed/pkg/config"
	influxdb "github.com/influxdata/influxdb/client"
)

const (
	Kind = "InfluxDB"
)

type Janitor struct {
	Spec config.InfluxDBSpec
	TTL  time.Duration
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
		Command:  fmt.Sprintf("ALTER RETENTION POLICY default ON k8s DURATION %vs", int(math.Ceil(j.TTL.Seconds()))),
		Database: "k8s",
	}
	_, err = client.Query(query)
	return err
}
