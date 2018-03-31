package kubed

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +genclient
// +genclient:onlyVerbs=get
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SearchResult struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Hits     []ResultEntry
	Total    uint64
	MaxScore float64
	Took     metav1.Duration
}

var _ runtime.Object = &SearchResult{}

type ResultEntry struct {
	Score  float64
	Object runtime.RawExtension
}
