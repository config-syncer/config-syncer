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
		RecycleBin: &config.RecycleBinSpec{
			Path:          "/tmp/kubed/trash",
			TTL:           metav1.Duration{Duration: 7 * 24 * time.Hour},
			HandleUpdates: false,
			Receiver: &config.Receiver{
				To:       []string{"ops@example.com"},
				Notifier: mailgun.UID,
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
				}},
			Receiver: config.Receiver{
				To:       []string{"ops@example.com"},
				Notifier: mailgun.UID,
			},
		},
		Snapshotter: &config.SnapshotSpec{
			Schedule: "@every 6h",
			Sanitize: true,
			Storage: config.Backend{
				StorageSecretName: "snap-secret",
				GCS: &config.GCSSpec{
					Bucket: "restic",
					Prefix: "minikube",
				},
			},
		},
		NotifierSecretName: "kubed-notifier",
	}
}
