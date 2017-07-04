package janitor

import (
	"strconv"
	"sync"
	"time"

	es "github.com/appscode/kubed/pkg/janitor/elasticsearch"
	influx "github.com/appscode/kubed/pkg/janitor/influxdb"
	"github.com/appscode/kubed/pkg/util"
	"github.com/appscode/log"
	"github.com/appscode/searchlight/pkg/client/icinga"
	influxdb "github.com/influxdata/influxdb/client"
	elastic "gopkg.in/olivere/elastic.v3"
)

const (
	ESEndpoint                string = "es-endpoint"
	LogIndexPrefix            string = "log-index-prefix"
	LogStorageLifetime        string = "log-storage-lifetime"
	MonitoringStorageLifetime string = "monitoring-storage-lifetime"
)

type Janitor struct {
	ClusterName                         string
	ElasticConfig                       map[string]string
	InfluxConfig                        influxdb.Config
	IcingaConfig                        map[string]string
	ClusterKubedConfigSecretMountedPath string

	// Icinga Client
	IcingaClient *icinga.IcingaClient

	once sync.Once
}

type ClusterSettings struct {
	LogIndexPrefix            string `json:"log_index_prefix"`
	LogStorageLifetime        int64  `json:"log_storage_lifetime"`
	MonitoringStorageLifetime int64  `json:"monitoring_storage_lifetime"`
}

func (j *Janitor) Run() {
	j.once.Do(func() {
		// wait for the first time for starting up the other pods
		time.Sleep(time.Minute * 10)
	})

	cs, err := getClusterSettings(j.ClusterKubedConfigSecretMountedPath)
	if err != nil {
		log.Errorln(err)
		return
	}
	log.Infof("Cluster settings: %+v", cs)
	j.cleanES(cs)
	j.cleanInflux(cs)
}

func (j *Janitor) cleanES(k ClusterSettings) error {
	if value, ok := j.ElasticConfig[ESEndpoint]; ok {
		esClient, err := elastic.NewClient(
			// elastic.SetSniff(false),
			elastic.SetURL(value),
		)
		if err != nil {
			log.Errorln(err)
			return err
		}
		return es.DeleteIndices(esClient, k.LogIndexPrefix, k.LogStorageLifetime)
	} else {
		log.Infoln("elastic config url not set, ignoring elastic clean")
	}
	return nil
}

func (j *Janitor) cleanInflux(k ClusterSettings) error {
	influxClient, err := influxdb.NewClient(j.InfluxConfig)
	if err != nil {
		log.Errorln(err)
		return err
	}
	return influx.UpdateRetentionPolicy(influxClient, k.MonitoringStorageLifetime)
}

func getClusterSettings(secretMountedPath string) (ClusterSettings, error) {
	cs := ClusterSettings{}
	var err error
	m, err := util.MountedSecretToMap(secretMountedPath)
	if err != nil {
		return cs, err
	}
	if d, ok := m[LogIndexPrefix]; ok {
		cs.LogIndexPrefix = string(d)
	}
	if d, ok := m[LogStorageLifetime]; ok {
		cs.LogStorageLifetime, err = strconv.ParseInt(string(d), 10, 64)
		if err != nil {
			return cs, err
		}
	}
	if d, ok := m[MonitoringStorageLifetime]; ok {
		cs.MonitoringStorageLifetime, err = strconv.ParseInt(string(d), 10, 64)
		if err != nil {
			return cs, err
		}
	}
	return cs, nil

}
