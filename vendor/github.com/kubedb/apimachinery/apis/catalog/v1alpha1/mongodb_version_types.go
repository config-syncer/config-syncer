package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	ResourceCodeMongoDBVersion     = "mgversion"
	ResourceKindMongoDBVersion     = "MongoDBVersion"
	ResourceSingularMongoDBVersion = "mongodbversion"
	ResourcePluralMongoDBVersion   = "mongodbversions"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MongoDBVersion defines a MongoDB database version.
type MongoDBVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MongoDBVersionSpec `json:"spec,omitempty"`
}

// MongoDBVersionSpec is the spec for mongodb version
type MongoDBVersionSpec struct {
	// Version
	Version string `json:"version"`
	// Database Image
	DB MongoDBVersionDatabase `json:"db"`
	// Exporter Image
	Exporter MongoDBVersionExporter `json:"exporter"`
	// Tools Image
	Tools MongoDBVersionTools `json:"tools"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`
	// Init container Image
	InitContainer MongoDBVersionInitContainer `json:"initContainer"`
	// PSP names
	PodSecurityPolicies MongoDBVersionPodSecurityPolicy `json:"podSecurityPolicies"`
}

// MongoDBVersionDatabase is the MongoDB Database image
type MongoDBVersionDatabase struct {
	Image string `json:"image"`
}

// MongoDBVersionExporter is the image for the MongoDB exporter
type MongoDBVersionExporter struct {
	Image string `json:"image"`
}

// MongoDBVersionTools is the image for the mongodb tools
type MongoDBVersionTools struct {
	Image string `json:"image"`
}

// MongoDBVersionInitContainer is the Elasticsearch Container initializer
type MongoDBVersionInitContainer struct {
	Image string `json:"image"`
}

// MongoDBVersionPodSecurityPolicy is the MongoDB pod security policies
type MongoDBVersionPodSecurityPolicy struct {
	DatabasePolicyName    string `json:"databasePolicyName"`
	SnapshotterPolicyName string `json:"snapshotterPolicyName"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MongoDBVersionList is a list of MongoDBVersions
type MongoDBVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of MongoDBVersion CRD objects
	Items []MongoDBVersion `json:"items,omitempty"`
}
