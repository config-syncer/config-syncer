package influx

import (
	"fmt"
	"math"
	"net/url"
	"time"

	"github.com/appscode/go/log"
	"github.com/appscode/kubed/pkg/config"
	influxdb "github.com/influxdata/influxdb/client"
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

	/*
		// ref: https://docs.influxdata.com/influxdb/v1.4/query_language/schema_exploration/#example-2-run-a-show-retention-policies-query-without-the-on-clause
		$ curl -G "http://localhost:8086/query?db=k8s&pretty=true" --data-urlencode "q=SHOW RETENTION POLICIES"
	*/
	query := influxdb.Query{
		Command:  fmt.Sprintf("ALTER RETENTION POLICY default ON k8s DURATION %vs", int(math.Ceil(j.TTL.Seconds()))),
		Database: "k8s",
	}
	_, err = client.Query(query)
	if err != nil {
		log.Warningf("failed to ALTER RETENTION POLICY for k8s database. Reason: %v", err)
	} else {
		log.Infoln("successfully ALTER-ed RETENTION POLICY for k8s database")
	}
	return err
}
