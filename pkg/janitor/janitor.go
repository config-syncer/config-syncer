package janitor

import (
	"strconv"
	"sync"
	"time"

	es "github.com/appscode/kubed/pkg/janitor/elasticsearch"
	influx "github.com/appscode/kubed/pkg/janitor/influxdb"
	"github.com/appscode/log"
	"github.com/appscode/searchlight/pkg/client/icinga"
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
	ClusterName                  string
	ElasticConfig                map[string]string
	InfluxConfig                 influxdb.Config
	IcingaConfig                 map[string]string
	KubeClient                   clientset.Interface
	ClusterKubedConfigSecretName string
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

	cs, err := getClusterSettings(j.KubeClient, j.ClusterKubedConfigSecretName)
	if err != nil {
		log.Errorln(err)
		return
	}
	log.Infof("Cluster settings: %+v",  cs)
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

func getClusterSettings(client clientset.Interface, secretName string) (ClusterSettings, error) {
	clusterConf, err := client.Core().
		Secrets("kube-system").
		Get(secretName, meta_v1.GetOptions{})
	if err != nil {
		return ClusterSettings{}, err
	}
	return SecretToClusterSettings(*clusterConf)
}

func SecretToClusterSettings(cnf apiv1.Secret) (ClusterSettings, error) {
	cs := ClusterSettings{}
	var err error
	if d, ok := cnf.Data[LogIndexPrefix]; ok {
		cs.LogIndexPrefix = string(d)
	}
	if d, ok := cnf.Data[LogStorageLifetime]; ok {
		cs.LogStorageLifetime, err = strconv.ParseInt(string(d), 10, 64)
		if err != nil {
			return cs, err
		}
	}
	if d, ok := cnf.Data[MonitoringStorageLifetime]; ok {
		cs.MonitoringStorageLifetime, err = strconv.ParseInt(string(d), 10, 64)
		if err != nil {
			return cs, err
		}
	}
	return cs, nil
}
