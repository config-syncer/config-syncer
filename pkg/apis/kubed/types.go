package kubed

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SearchResult struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Hits     []ResultEntry   `json:"hits,omitempty"`
	Total    uint64          `json:"totalHits"`
	MaxScore float64         `json:"maxScore"`
	Took     metav1.Duration `json:"took"`
}

var _ runtime.Object = &SearchResult{}

type ResultEntry struct {
	Score  float64              `json:"score"`
	Object runtime.RawExtension `json:"object,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SearchResultList is a list of SearchResult objects.
type SearchResultList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []SearchResult
}
