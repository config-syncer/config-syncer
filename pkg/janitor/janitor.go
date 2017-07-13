package janitor

import (
	"strconv"
	"sync"
	"time"

	"github.com/appscode/kubed/pkg/config"
	es "github.com/appscode/kubed/pkg/janitor/elasticsearch"
	influx "github.com/appscode/kubed/pkg/janitor/influxdb"
	"github.com/appscode/log"
	influxdb "github.com/influxdata/influxdb/client"
	elastic "gopkg.in/olivere/elastic.v3"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

const (
	ESEndpoint                string = "es-endpoint"
	LogIndexPrefix            string = "log-index-prefix"
	LogStorageLifetime        string = "log-storage-lifetime"
	MonitoringStorageLifetime string = "monitoring-storage-lifetime"
)

type Janitor struct {
	KubeClient                        clientset.Interface

	ElasticConfig                     map[string]string
	InfluxConfig                      influxdb.Config
	ClusterKubedConfigSecretName      string
	ClusterKubedConfigSecretNamespace string

	once sync.Once
}

func (j *Janitor) Run() {
	j.once.Do(func() {
		// wait for the first time for starting up the other pods
		time.Sleep(time.Minute * 10)
	})

	cs, err := getClusterSettings(j.KubeClient, j.ClusterKubedConfigSecretName, j.ClusterKubedConfigSecretNamespace)
	if err != nil {
		log.Errorln(err)
		return
	}
	log.Infof("Cluster settings: %+v", cs)
	j.cleanES(cs)
	j.cleanInflux(cs)
}

func (j *Janitor) cleanES(k config.ClusterConfig) error {
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

func (j *Janitor) cleanInflux(k config.ClusterConfig) error {
	influxClient, err := influxdb.NewClient(j.InfluxConfig)
	if err != nil {
		log.Errorln(err)
		return err
	}
	return influx.UpdateRetentionPolicy(influxClient, k.MonitoringStorageLifetime)
}
