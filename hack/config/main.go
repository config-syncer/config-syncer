package main

import (
	"io/ioutil"
	"time"

	"github.com/appscode/go/log"
	"github.com/appscode/go/runtime"
	apis "github.com/appscode/kubed/pkg/apis/v1alpha1"
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

func CreateClusterConfig() apis.ClusterConfig {
	return apis.ClusterConfig{
		ClusterName: "unicorn",
		APIServer: apis.APIServerSpec{
			Address:           ":8080",
			EnableSearchIndex: true,
		},
		Snapshotter: &apis.SnapshotSpec{
			Schedule: "@every 6h",
			Sanitize: true,
			Backend: apis.Backend{
				StorageSecretName: "snap-secret",
				GCS: &apis.GCSSpec{
					Bucket: "restic",
					Prefix: "minikube",
				},
			},
		},
		RecycleBin: &apis.RecycleBinSpec{
			Path:          "/tmp/kubed/trash",
			TTL:           metav1.Duration{Duration: 7 * 24 * time.Hour},
			HandleUpdates: false,
			Receivers: []apis.Receiver{{
				To:       []string{"ops@example.com"},
				Notifier: "Mailgun",
			},
			},
		},
		EnableConfigSyncer: true,
		EventForwarder: &apis.EventForwarderSpec{
			Rules: []apis.PolicyRule{
				{
					Operations: []apis.Operation{apis.Create},
					Resources: []apis.GroupResources{
						{
							Group: "",
							Resources: []string{
								"events",
							},
						},
					},
					Namespaces: []string{"kube-system"},
				},
				{
					Operations: []apis.Operation{apis.Create},
					Resources: []apis.GroupResources{
						{
							Group: "",
							Resources: []string{
								"nodes",
								"persistentvolumes",
								"persistentvolumeclaims",
							},
						},
					},
				},
				{
					Operations: []apis.Operation{apis.Create},
					Resources: []apis.GroupResources{
						{
							Group: "storage.k8s.io",
							Resources: []string{
								"storageclasses",
							},
						},
					},
				},
				{
					Operations: []apis.Operation{apis.Create},
					Resources: []apis.GroupResources{
						{
							Group: "extensions",
							Resources: []string{
								"ingresses",
							},
						},
					},
				},
				{
					Operations: []apis.Operation{apis.Create},
					Resources: []apis.GroupResources{
						{
							Group: "voyager.appscode.com",
							Resources: []string{
								"ingresses",
							},
						},
					},
				},
				{
					Operations: []apis.Operation{apis.Create},
					Resources: []apis.GroupResources{
						{
							Group: "certificates.k8s.io",
							Resources: []string{
								"certificatesigningrequests",
							},
						},
					},
				},
			},
			Receivers: []apis.Receiver{
				{
					To:       []string{"ops@example.com"},
					Notifier: "Mailgun",
				},
			},
		},
		Janitors: []apis.JanitorSpec{
			{
				Kind: apis.JanitorElasticsearch,
				TTL:  metav1.Duration{Duration: 90 * 24 * time.Hour},
				Elasticsearch: &apis.ElasticsearchSpec{
					Endpoint:       "https://elasticsearch-logging.kube-system:9200",
					LogIndexPrefix: "logstash-",
					SecretName:     "elasticsearch-logging-cert",
				},
			},
			{
				Kind: apis.JanitorInfluxDB,
				TTL:  metav1.Duration{Duration: 90 * 24 * time.Hour},
				InfluxDB: &apis.InfluxDBSpec{
					Endpoint: "https://monitoring-influxdb.kube-system:8086",
				},
			},
		},
		NotifierSecretName: "notifier-config",
	}
}
