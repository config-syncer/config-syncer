package api

type KubedMetadata struct {
	OperatorNamespace   string      `json:"operatorNamespace,omitempty"`
	SearchEnabled       bool        `json:"searchEnabled"`
	ReverseIndexEnabled bool        `json:"reverseIndexEnabled"`
	AnalyticsEnabled    bool        `json:"analyticsEnabled"`
	Version             interface{} `json:"version,omitempty"`
}
