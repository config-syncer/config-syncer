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

	bytes, err := yaml.Marshal(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	cfgmap := apiv1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kubed",
			Namespace: "kube-system",
		},
		Data: map[string]string{
			"config.yaml": string(bytes),
		},
	}
	bytes, err = yaml.Marshal(cfgmap)
	if err != nil {
		log.Fatalln(err)
	}
	p := runtime.GOPath() + "/src/github.com/appscode/kubed/hack/deploy/configmap.yaml"
	ioutil.WriteFile(p, bytes, 0644)
}

func CreateClusterConfig() config.ClusterConfig {
	return config.ClusterConfig{
		TrashCan: &config.TrashCanSpec{
			Path:         "/tmp/kubed",
			TTL:          metav1.Duration{7 * 24 * time.Hour},
			HandleUpdate: true,
			NotifyVia:    mailgun.UID,
		},
		EventForwarder: &config.EventForwarderSpec{
			NotifyOnStorageAdd:   true,
			NotifyOnIngressAdd:   true,
			ForwardWarningEvents: true,
			EventNamespaces:      []string{"kube-system"},
			NotifyVia:            mailgun.UID,
		},
		ClusterSnapshot: &config.SnapshotSpec{
			Schedule: "@every 5m",
			Sanitize: true,
			Storage: config.Backend{
				StorageSecretName: "snap-secret",
				GCS: &config.GCSSpec{
					Bucket: "restic",
					Prefix: "a/b/c",
				},
			},
		},
		NotifierSecretName: "kubed-notifier",
	}
}
