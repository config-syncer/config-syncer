package framework

import (
	api "github.com/appscode/kubed/apis/kubed/v1alpha1"
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
