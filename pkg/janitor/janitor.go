package janitor

import (
	"sync"
	"time"

	api "github.com/appscode/api/kubernetes/v1beta1"
	"github.com/appscode/client"
	es "github.com/appscode/kubed/pkg/janitor/elasticsearch"
	influx "github.com/appscode/kubed/pkg/janitor/influxdb"
	"github.com/appscode/log"
	"github.com/appscode/searchlight/pkg/client/icinga"
	influxdb "github.com/influxdata/influxdb/client"
	elastic "gopkg.in/olivere/elastic.v3"
)

const (
	ESEndpoint string = "es-endpoint"
)

type Janitor struct {
	ClusterName   string
	ElasticConfig map[string]string
	InfluxConfig  influxdb.Config
	IcingaConfig  map[string]string

	// appscode api server client
	APIClientOptions *client.ClientOption

	// Icinga Client
	IcingaClient *icinga.IcingaClient

	once sync.Once
}

func (j *Janitor) Run() {
	j.once.Do(func() {
		// wait for the first time for starting up the other pods
		time.Sleep(time.Minute * 10)
	})

	conn, err := client.New(j.APIClientOptions)
	if err != nil {
		log.Errorln(err)
		return
	}
	defer conn.Close()

	clusterDescribeReq := &api.ClusterDescribeRequest{
		Uid: j.ClusterName,
	}
	clusterDescribeResp, err := conn.Kubernetes().V1beta1().Cluster().Describe(conn.Context(), clusterDescribeReq)
	if err != nil {
		log.Errorln(err)
		return
	}

	// This ensures not panic if there were some communication
	// error with apiserver.
	if clusterDescribeResp != nil {
		if clusterDescribeResp.Cluster != nil {
			if clusterDescribeResp.Cluster.Settings == nil {
				log.Warningln("failed to get cluster settings informations")
				return
			}

			j.cleanES(clusterDescribeResp.Cluster.Settings)
			j.cleanInflux(clusterDescribeResp.Cluster.Settings)
			//j.syncAlert(conn)
		}
	}
}

func (j *Janitor) cleanES(k *api.ClusterSettings) error {
	if value, ok := j.ElasticConfig[ESEndpoint]; ok {
		esClient, err := elastic.NewClient(
			// elastic.SetSniff(false),
			elastic.SetURL(value),
		)
		if err != nil {
			log.Errorln(err)
			return err
		}
		return es.DeleteIndices(esClient, k)
	} else {
		log.Infoln("elastic config url not set, ignoring elastic clean")
	}
	return nil
}

func (j *Janitor) cleanInflux(k *api.ClusterSettings) error {
	influxClient, err := influxdb.NewClient(j.InfluxConfig)
	if err != nil {
		log.Errorln(err)
		return err
	}
	return influx.UpdateRetentionPolicy(influxClient, k)
}
