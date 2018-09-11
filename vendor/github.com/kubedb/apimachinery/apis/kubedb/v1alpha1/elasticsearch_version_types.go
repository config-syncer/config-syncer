package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	ResourceCodeElasticsearchVersion     = "esversion"
	ResourceKindElasticsearchVersion     = "ElasticsearchVersion"
	ResourceSingularElasticsearchVersion = "elasticsearchversion"
	ResourcePluralElasticsearchVersion   = "elasticsearchversions"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ElasticsearchVersion defines a Elasticsearch database version.
type ElasticsearchVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ElasticsearchVersionSpec `json:"spec,omitempty"`
}

// ElasticsearchVersionSpec is the spec for elasticsearch version
type ElasticsearchVersionSpec struct {
	// Version
	Version string `json:"version"`
	// Database Image
	DB ElasticsearchVersionDatabase `json:"db"`
	// Exporter Image
	Exporter ElasticsearchVersionExporter `json:"exporter"`
	// Tools Image
	Tools ElasticsearchVersionTools `json:"tools"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`
}

// ElasticsearchVersionDatabase is the Elasticsearch Database image
type ElasticsearchVersionDatabase struct {
	Image string `json:"image"`
}

// ElasticsearchVersionExporter is the image for the Elasticsearch exporter
type ElasticsearchVersionExporter struct {
	Image string `json:"image"`
}

// ElasticsearchVersionTools is the image for the elasticsearch tools
type ElasticsearchVersionTools struct {
	Image string `json:"image"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ElasticsearchVersionList is a list of ElasticsearchVersions
type ElasticsearchVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of ElasticsearchVersion CRD objects
	Items []ElasticsearchVersion `json:"items,omitempty"`
}
