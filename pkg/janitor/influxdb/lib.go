package influx

import (
	"fmt"

	"github.com/appscode/log"
	influxdb "github.com/influxdata/influxdb/client"
)

func UpdateRetentionPolicy(influxClient *influxdb.Client, monitoringStorageLifetime int64) error {
	query := influxdb.Query{
		Command:  fmt.Sprintf("ALTER RETENTION POLICY default ON k8s DURATION %vs", monitoringStorageLifetime),
		Database: "k8s",
	}
	if _, err := influxClient.Query(query); err != nil {
		log.Errorln(err)
		return err
	}
	return nil
}
