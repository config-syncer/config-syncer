package influx

import (
	"fmt"
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
	ttl := j.TTL
	// https://docs.influxdata.com/influxdb/v1.3/query_language/database_management/#create-retention-policies-with-create-retention-policy
	if ttl < 60*time.Minute {
		ttl = 60 * time.Minute
		log.Infof("influx janitor [%s]: resetting retention duration to minimum %s", j.Spec.Endpoint, ttl)
	}
	query := fmt.Sprintf("ALTER RETENTION POLICY default ON k8s DURATION %s SHARD DURATION 0s DEFAULT", ttl)
	log.Infof("influx janitor [%s]: %s", j.Spec.Endpoint, query)
	resp, err := client.Query(influxdb.Query{
		Command:  query,
		Database: "k8s",
	})
	if err == nil && resp.Err != nil {
		err = resp.Err
	}
	if err != nil {
		log.Warningf("failed to ALTER RETENTION POLICY for k8s database. Reason: %v", err)
		return err
	}

	log.Infoln("successfully altered retention policy for k8s database")
	return nil
}
