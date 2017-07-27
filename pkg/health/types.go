package health

import "github.com/appscode/api/version"

type KubedHealth struct {
	OperatorNamespace   string           `json:"operatorNamespace,omitempty"`
	SearchEnabled       bool             `json:"searchEnabled"`
	ReverseIndexEnabled bool             `json:"reverseIndexEnabled"`
	AnalyticsEnabled    bool             `json:"analyticsEnabled"`
	Version             *version.Version `json:"version,omitempty"`
}
