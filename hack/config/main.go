package main

import (
	"io/ioutil"
	"time"

	"github.com/appscode/go-notify/mailgun"
	"github.com/appscode/go/runtime"
	"github.com/appscode/kubed/pkg/config"
	"github.com/appscode/log"
	"github.com/ghodss/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func main() {
	cfg := CreateClusterConfig()
	cfg.Save(runtime.GOPath() + "/src/github.com/appscode/kubed/hack/deploy/config.yaml")

	cfgBytes, err := yaml.Marshal(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	cfgmap := apiv1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kubed-config",
			Namespace: "kube-system",
			Labels: map[string]string{
				"app": "kubed",
			},
		},
		Data: map[string][]byte{
			"config.yaml": cfgBytes,
		},
	}
	bytes, err := yaml.Marshal(cfgmap)
	if err != nil {
		log.Fatalln(err)
	}
	p := runtime.GOPath() + "/src/github.com/appscode/kubed/hack/deploy/kubed-config.yaml"
	ioutil.WriteFile(p, bytes, 0644)
}

func CreateClusterConfig() config.ClusterConfig {
	return config.ClusterConfig{
		APIServer: config.APIServerSpec{
			Address:            ":8080",
			EnableSearchIndex:  true,
			EnableReverseIndex: true,
		},
		Snapshotter: &config.SnapshotSpec{
			Schedule: "@every 6h",
			Sanitize: true,
			Backend: config.Backend{
				StorageSecretName: "snap-secret",
				GCS: &config.GCSSpec{
					Bucket: "restic",
					Prefix: "minikube",
				},
			},
		},
		RecycleBin: &config.RecycleBinSpec{
			Path:          "/tmp/kubed/trash",
			TTL:           metav1.Duration{Duration: 7 * 24 * time.Hour},
			HandleUpdates: false,
			Receivers: []config.Receiver{{
				To:       []string{"ops@example.com"},
				Notifier: mailgun.UID,
			},
			},
		},
		EnableConfigSyncer: true,
		EventForwarder: &config.EventForwarderSpec{
			NodeAdded: config.ForwarderSpec{
				Handle: true,
			},
			StorageAdded: config.ForwarderSpec{
				Handle: true,
			},
			IngressAdded: config.ForwarderSpec{
				Handle: true,
			},
			WarningEvents: config.ForwarderSpec{
				Handle: true,
				Namespaces: []string{
					"kube-system",
				},
			},
			Receivers: []config.Receiver{{
				To:       []string{"ops@example.com"},
				Notifier: mailgun.UID,
			},
			},
		},
		Janitors: []config.JanitorSpec{
			{
				Kind: config.JanitorElasticsearch,
				TTL:  metav1.Duration{Duration: 90 * 24 * time.Hour},
				Elasticsearch: &config.ElasticSearchSpec{
					Endpoint:       "http://elasticsearch-logging.kube-system:9200",
					LogIndexPrefix: "logstash-",
				},
			},
			{
				Kind: config.JanitorInfluxDB,
				TTL:  metav1.Duration{Duration: 90 * 24 * time.Hour},
				InfluxDB: &config.InfluxDBSpec{
					Endpoint: "https://monitoring-influxdb.kube-system:8086",
				},
			},
		},
		NotifierSecretName: "notifier-config",
	}
}
