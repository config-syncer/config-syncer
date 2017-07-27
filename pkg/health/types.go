package health

import "github.com/appscode/api/version"

type KubedHealth struct {
	OperatorNamespace   string           `json:"operator_namespace,omitempty"`
	SearchEnabled       bool             `json:"search_enabled"`
	ReverseIndexEnabled bool             `json:"reverse_index_enabled"`
	AnalyticsEnabled    bool             `json:"analytics_enabled"`
	Version             *version.Version `json:"version,omitempty"`
}
