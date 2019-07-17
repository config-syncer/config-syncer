package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	ResourceCodePerconaVersion     = "pcversion"
	ResourceKindPerconaVersion     = "PerconaVersion"
	ResourceSingularPerconaVersion = "perconaversion"
	ResourcePluralPerconaVersion   = "perconaversions"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PerconaVersion defines a Percona database (percona variation for MySQL database) version.
type PerconaVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PerconaVersionSpec `json:"spec,omitempty"`
}

// PerconaVersionSpec is the spec for Percona version
type PerconaVersionSpec struct {
	// Version
	Version string `json:"version"`
	// Database Image
	DB PerconaVersionDatabase `json:"db"`
	// Proxysql Image
	Proxysql PerconaVersionProxysql `json:"proxysql"`
	// Exporter Image
	Exporter PerconaVersionExporter `json:"exporter"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`
	// Init container Image
	// TODO: remove if not needed
	InitContainer PerconaVersionInitContainer `json:"initContainer"`
	// PSP names
	PodSecurityPolicies PerconaVersionPodSecurityPolicy `json:"podSecurityPolicies"`
}

// PerconaVersionDatabase is the percona image
type PerconaVersionDatabase struct {
	Image string `json:"image"`
}

// PerconaVersionProxysql is the proxysql image
type PerconaVersionProxysql struct {
	Image string `json:"image"`
}

// PerconaVersionExporter is the image for the Percona exporter
type PerconaVersionExporter struct {
	Image string `json:"image"`
}

// PerconaVersionInitContainer is the Percona Container initializer
type PerconaVersionInitContainer struct {
	Image string `json:"image"`
}

// PerconaVersionPodSecurityPolicy is the Percona pod security policies
type PerconaVersionPodSecurityPolicy struct {
	DatabasePolicyName string `json:"databasePolicyName"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PerconaVersionList is a list of PerconaVersions
type PerconaVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of PerconaVersion CRD objects
	Items []PerconaVersion `json:"items,omitempty"`
}
