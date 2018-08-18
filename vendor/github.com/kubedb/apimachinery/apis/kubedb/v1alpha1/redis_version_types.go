package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	ResourceCodeRedisVersion     = "rdversion"
	ResourceKindRedisVersion     = "RedisVersion"
	ResourceSingularRedisVersion = "redisversion"
	ResourcePluralRedisVersion   = "redisversions"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RedisVersion defines a Redis database version.
type RedisVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              RedisVersionSpec `json:"spec,omitempty"`
}

// RedisVersionSpec is the spec for redis version
type RedisVersionSpec struct {
	// Version
	Version string `json:"version,omitempty"`
	// Database Image
	DB RedisVersionDatabase `json:"db"`
	// Exporter Image
	Exporter RedisVersionExporter `json:"exporter"`
}

// RedisVersionDatabase is the Redis Database image
type RedisVersionDatabase struct {
	Image string `json:"image"`
}

// RedisVersionExporter is the image for the Redis exporter
type RedisVersionExporter struct {
	Image string `json:"image"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RedisVersionList is a list of RedisVersions
type RedisVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of RedisVersion CRD objects
	Items []RedisVersion `json:"items,omitempty"`
}
