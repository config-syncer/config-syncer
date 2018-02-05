package main

import (
	"io/ioutil"
	"time"

	"github.com/appscode/go/log"
	"github.com/appscode/go/runtime"
	"github.com/appscode/kubed/pkg/api"
	"github.com/ghodss/yaml"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	cfg := CreateClusterConfig()
	cfg.Save(runtime.GOPath() + "/src/github.com/appscode/kubed/hack/deploy/config.yaml")

	cfgBytes, err := yaml.Marshal(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	cfgmap := core.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
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

func CreateClusterConfig() api.ClusterConfig {
	return api.ClusterConfig{
		ClusterName: "unicorn",
		APIServer: api.APIServerSpec{
			Address:           ":8080",
			EnableSearchIndex: true,
		},
		Snapshotter: &api.SnapshotSpec{
			Schedule: "@every 6h",
			Sanitize: true,
			Backend: api.Backend{
				StorageSecretName: "snap-secret",
				GCS: &api.GCSSpec{
					Bucket: "restic",
					Prefix: "minikube",
				},
			},
		},
		RecycleBin: &api.RecycleBinSpec{
			Path:          "/tmp/kubed/trash",
			TTL:           metav1.Duration{Duration: 7 * 24 * time.Hour},
			HandleUpdates: false,
			Receivers: []api.Receiver{{
				To:       []string{"ops@example.com"},
				Notifier: "Mailgun",
			},
			},
		},
		EnableConfigSyncer: true,
		EventForwarder: &api.EventForwarderSpec{
			NodeAdded: api.ForwarderSpec{
				Handle: true,
			},
			StorageAdded: api.ForwarderSpec{
				Handle: true,
			},
			IngressAdded: api.ForwarderSpec{
				Handle: true,
			},
			WarningEvents: api.ForwarderSpec{
				Handle: true,
				Namespaces: []string{
					"kube-system",
				},
			},
			Receivers: []api.Receiver{{
				To:       []string{"ops@example.com"},
				Notifier: "Mailgun",
			},
			},
		},
		Janitors: []api.JanitorSpec{
			{
				Kind: api.JanitorElasticsearch,
				TTL:  metav1.Duration{Duration: 90 * 24 * time.Hour},
				Elasticsearch: &api.ElasticsearchSpec{
					Endpoint:       "https://elasticsearch-logging.kube-system:9200",
					LogIndexPrefix: "logstash-",
					SecretName:     "elasticsearch-logging-cert",
				},
			},
			{
				Kind: api.JanitorInfluxDB,
				TTL:  metav1.Duration{Duration: 90 * 24 * time.Hour},
				InfluxDB: &api.InfluxDBSpec{
					Endpoint: "https://monitoring-influxdb.kube-system:8086",
				},
			},
		},
		NotifierSecretName: "notifier-config",
	}
}
