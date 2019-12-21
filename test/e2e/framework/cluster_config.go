/*
Copyright The Kubed Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package framework

import (
	"time"

	api "github.com/appscode/kubed/apis/kubed/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	store "kmodules.xyz/objectstore-api/api/v1"
)

func SnapshotterClusterConfig(backend *store.Backend) api.ClusterConfig {
	return api.ClusterConfig{
		Snapshotter: &api.SnapshotSpec{
			Backend:  *backend,
			Sanitize: true,
			Schedule: "@every 1m",
		},
	}
}

func ConfigSyncClusterConfig() api.ClusterConfig {
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
