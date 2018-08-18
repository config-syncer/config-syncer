package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	ResourceCodeMySQLVersion     = "myversion"
	ResourceKindMySQLVersion     = "MySQLVersion"
	ResourceSingularMySQLVersion = "mysqlversion"
	ResourcePluralMySQLVersion   = "mysqlversions"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MySQLVersion defines a MySQL database version.
type MySQLVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MySQLVersionSpec `json:"spec,omitempty"`
}

// MySQLVersionSpec is the spec for postgres version
type MySQLVersionSpec struct {
	// Version
	Version string `json:"version,omitempty"`
	// Database Image
	DB MySQLVersionDatabase `json:"db"`
	// Exporter Image
	Exporter MySQLVersionExporter `json:"exporter"`
	// Tools Image
	Tools MySQLVersionTools `json:"tools"`
}

// MySQLVersionDatabase is the MySQL Database image
type MySQLVersionDatabase struct {
	Image string `json:"image"`
}

// MySQLVersionExporter is the image for the MySQL exporter
type MySQLVersionExporter struct {
	Image string `json:"image"`
}

// MySQLVersionTools is the image for the postgres tools
type MySQLVersionTools struct {
	Image string `json:"image"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MySQLVersionList is a list of MySQLVersions
type MySQLVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of MySQLVersion CRD objects
	Items []MySQLVersion `json:"items,omitempty"`
}
