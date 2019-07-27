package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	ResourceCodeMemcachedVersion     = "mcversion"
	ResourceKindMemcachedVersion     = "MemcachedVersion"
	ResourceSingularMemcachedVersion = "memcachedversion"
	ResourcePluralMemcachedVersion   = "memcachedversions"
)

// MemcachedVersion defines a Memcached database version.

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=memcachedversions,singular=memcachedversion,scope=Cluster,shortName=mcversion,categories={datastore,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="DB_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type MemcachedVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MemcachedVersionSpec `json:"spec,omitempty"`
}

// MemcachedVersionSpec is the spec for memcached version
type MemcachedVersionSpec struct {
	// Version
	Version string `json:"version"`
	// Database Image
	DB MemcachedVersionDatabase `json:"db"`
	// Exporter Image
	Exporter MemcachedVersionExporter `json:"exporter"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`
	// PSP names
	PodSecurityPolicies MemcachedVersionPodSecurityPolicy `json:"podSecurityPolicies"`
}

// MemcachedVersionDatabase is the Memcached Database image
type MemcachedVersionDatabase struct {
	Image string `json:"image"`
}

// MemcachedVersionExporter is the image for the Memcached exporter
type MemcachedVersionExporter struct {
	Image string `json:"image"`
}

// MemcachedVersionPodSecurityPolicy is the Memcached pod security policies
type MemcachedVersionPodSecurityPolicy struct {
	DatabasePolicyName string `json:"databasePolicyName"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MemcachedVersionList is a list of MemcachedVersions
type MemcachedVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of MemcachedVersion CRD objects
	Items []MemcachedVersion `json:"items,omitempty"`
}
