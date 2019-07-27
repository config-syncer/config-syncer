package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	ResourceCodeEtcdVersion     = "etcversion"
	ResourceKindEtcdVersion     = "EtcdVersion"
	ResourceSingularEtcdVersion = "etcdversion"
	ResourcePluralEtcdVersion   = "etcdversions"
)

// EtcdVersion defines a Etcd database version.

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=etcdversions,singular=etcdversion,scope=Cluster,shortName=etcversion,categories={datastore,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="DB_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type EtcdVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              EtcdVersionSpec `json:"spec,omitempty"`
}

// EtcdVersionSpec is the spec for postgres version
type EtcdVersionSpec struct {
	// Version
	Version string `json:"version"`
	// Database Image
	DB EtcdVersionDatabase `json:"db"`
	// Exporter Image
	Exporter EtcdVersionExporter `json:"exporter"`
	// Tools Image
	Tools EtcdVersionTools `json:"tools"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`
}

// EtcdVersionDatabase is the Etcd Database image
type EtcdVersionDatabase struct {
	Image string `json:"image"`
}

// EtcdVersionExporter is the image for the Etcd exporter
type EtcdVersionExporter struct {
	Image string `json:"image"`
}

// EtcdVersionTools is the image for the Etcd exporter
type EtcdVersionTools struct {
	Image string `json:"image"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EtcdVersionList is a list of EtcdVersions
type EtcdVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of EtcdVersion CRD objects
	Items []EtcdVersion `json:"items,omitempty"`
}
