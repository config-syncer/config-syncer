package framework

import (
	"time"

	api "github.com/appscode/kubed/apis/kubed/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func SnapshotterClusterConfig(backend *api.Backend) api.ClusterConfig {
	return api.ClusterConfig{
		Snapshotter: &api.SnapshotSpec{
			Backend:  *backend,
			Sanitize: true,
			Schedule: "@every 1m",
		},
	}
}

func ConfigMapSyncClusterConfig() api.ClusterConfig {
	return api.ClusterConfig{
		EnableConfigSyncer: true,
	}
}

func (f *Invocation) EventForwarderClusterConfig() api.ClusterConfig {
	return api.ClusterConfig{
		EventForwarder: &api.EventForwarderSpec{
			Rules: []api.PolicyRule{
				{
					Operations: []api.Operation{api.Create},
					Namespaces: []string{f.namespace},
					Resources: []api.GroupResources{
						{
							Group: "",
							Resources: []string{
								"events",
							},
						},
					},
				},
				{
					Operations: []api.Operation{api.Create},
					Namespaces: []string{f.namespace},
					Resources: []api.GroupResources{
						{
							Group: "",
							Resources: []string{
								"persistentvolumeclaims",
							},
						},
					},
				},
			},
		},
	}
}

func APIServerClusterConfig() api.ClusterConfig {
	return api.ClusterConfig{
		ClusterName: "minikube",
	}
}

func RecycleBinClusterConfig() api.ClusterConfig {
	return api.ClusterConfig{
		RecycleBin: &api.RecycleBinSpec{
			Path: "/tmp/kubed/trash",
			TTL: metav1.Duration{
				Duration: time.Hour,
			},
			HandleUpdates: false,
		},
	}
}

func WebhookReceiver() []api.Receiver {
	return []api.Receiver{
		{
			To:       []string{"ops-alerts"},
			Notifier: "Webhook",
		},
	}
}
