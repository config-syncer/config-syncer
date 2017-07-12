package config

type ClusterSettings struct {
	LogIndexPrefix            string `json:"log_index_prefix"`
	LogStorageLifetime        int64  `json:"log_storage_lifetime"`
	MonitoringStorageLifetime int64  `json:"monitoring_storage_lifetime"`
}
