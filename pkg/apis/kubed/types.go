package kubed

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +genclient
// +genclient:onlyVerbs=get
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Stuff struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Hits     []ResultEntry   `json:"hits,omitempty"`
	Total    uint64          `json:"totalHits"`
	MaxScore float64         `json:"maxScore"`
	Took     metav1.Duration `json:"took"`
}

var _ runtime.Object = &Stuff{}

type ResultEntry struct {
	Score  float64              `json:"score"`
	Object runtime.RawExtension `json:"object,omitempty"`
}
