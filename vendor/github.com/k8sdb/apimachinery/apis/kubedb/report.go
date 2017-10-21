package kubedb

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen=false
// +k8s:gen-deepcopy=false
type Report struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Summary           ReportSummary `json:"summary,omitempty"`
	Status            ReportStatus  `json:"status,omitempty"`
}

// +k8s:deepcopy-gen=false
// +k8s:gen-deepcopy=false
type ReportSummary struct {
	Postgres      map[string]*PostgresSummary      `json:"postgres,omitempty"`
	Elasticsearch map[string]*ElasticsearchSummary `json:"elastic,omitempty"`
}

// +k8s:deepcopy-gen=false
// +k8s:gen-deepcopy=false
type ReportStatus struct {
	StartTime      *metav1.Time `json:"startTime,omitempty"`
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`
}
