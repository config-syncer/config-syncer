package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	ResourceCodePostgresVersion     = "pgversion"
	ResourceKindPostgresVersion     = "PostgresVersion"
	ResourceSingularPostgresVersion = "postgresversion"
	ResourcePluralPostgresVersion   = "postgresversions"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PostgresVersion defines a Postgres database version.
type PostgresVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PostgresVersionSpec `json:"spec,omitempty"`
}

// PostgresVersionSpec is the spec for postgres version
type PostgresVersionSpec struct {
	// Version
	Version string `json:"version,omitempty"`
	// Database Image
	DB PostgresVersionDatabase `json:"db"`
	// Exporter Image
	Exporter PostgresVersionExporter `json:"exporter"`
	// Tools Image
	Tools PostgresVersionTools `json:"tools"`
}

// PostgresVersionDatabase is the Postgres Database image
type PostgresVersionDatabase struct {
	Image string `json:"image"`
}

// PostgresVersionExporter is the image for the Postgres exporter
type PostgresVersionExporter struct {
	Image string `json:"image"`
}

// PostgresVersionTools is the image for the postgres tools
type PostgresVersionTools struct {
	Image string `json:"image"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PostgresVersionList is a list of PostgresVersions
type PostgresVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of PostgresVersion CRD objects
	Items []PostgresVersion `json:"items,omitempty"`
}
