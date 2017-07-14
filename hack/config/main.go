package main

import (
	"github.com/appscode/kubed/pkg/config"
	"github.com/appscode/go/runtime"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
	"github.com/appscode/go-notify/mailgun"
)

func main() {
	cfg := config.ClusterConfig{
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
				StorageSecretName: "",
				Local: &config.LocalSpec{
					Path: "/tmp/csnap",
				},
			},
		},
		NotifierSecretName: "",
	}
	cfg.Save(runtime.GOPath() + "/src/github.com/appscode/kubed/hack/config/kubed.yaml")
}
