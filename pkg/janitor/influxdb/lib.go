package influx

import (
	"fmt"

	api "github.com/appscode/api/kubernetes/v1beta1"
	"github.com/appscode/log"
	influxdb "github.com/influxdata/influxdb/client"
)

func UpdateRetentionPolicy(influxClient *influxdb.Client, j *api.ClusterSettings) error {
	query := influxdb.Query{
		Command:  fmt.Sprintf("ALTER RETENTION POLICY default ON k8s DURATION %vs", j.MonitoringStorageLifetime),
		Database: "k8s",
	}
	if _, err := influxClient.Query(query); err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}
