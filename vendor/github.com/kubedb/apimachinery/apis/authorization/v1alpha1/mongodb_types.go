package v1alpha1

import (
	"github.com/appscode/go/encoding/json/types"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
)

const (
	ResourceKindMongoDBRole = "MongoDBRole"
	ResourceMongoDBRole     = "mongodbrole"
	ResourceMongoDBRoles    = "mongodbroles"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MongoDBRole
type MongoDBRole struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MongoDBRoleSpec   `json:"spec,omitempty"`
	Status            MongoDBRoleStatus `json:"status,omitempty"`
}

// MongoDBRoleSpec contains connection information, Mongodb role info etc
type MongoDBRoleSpec struct {
	AuthManagerRef *appcat.AppReference `json:"authManagerRef,omitempty"`

	DatabaseRef *core.LocalObjectReference `json:"databaseRef"`

	// links:
	// 	- https://www.vaultproject.io/api/secret/databases/index.html
	//	- https://www.vaultproject.io/api/secret/databases/mongodb.html

	// Specifies the TTL for the leases associated with this role.
	// Accepts time suffixed strings ("1h") or an integer number of seconds.
	// Defaults to system/engine default TTL time
	DefaultTTL string `json:"defaultTTL,omitempty"`

	// Specifies the maximum TTL for the leases associated with this role.
	// Accepts time suffixed strings ("1h") or an integer number of seconds.
	// Defaults to system/engine default TTL time.
	MaxTTL string `json:"maxTTL,omitempty"`

	// https://www.vaultproject.io/api/secret/databases/Mongodb-maria.html#creation_statements
	// Specifies the database statements executed to create and configure a user.
	CreationStatements []string `json:"creationStatements"`

	// https://www.vaultproject.io/api/secret/databases/Mongodb-maria.html#revocation_statements
	// Specifies the database statements to be executed to revoke a user.
	RevocationStatements []string `json:"revocationStatements,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type MongoDBRoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is a list of MongoDBRole objects
	Items []MongoDBRole `json:"items,omitempty"`
}

type MongoDBRolePhase string

type MongoDBRoleStatus struct {
	Phase MongoDBRolePhase `json:"phase,omitempty"`

	// observedGeneration is the most recent generation observed for this MongoDBRole. It corresponds to the
	// MongoDBRole's generation, which is updated on mutation by the API Server.
	ObservedGeneration *types.IntHash `json:"observedGeneration,omitempty"`

	// Represents the latest available observations of a MongoDBRole current state.
	Conditions []MongoDBRoleCondition `json:"conditions,omitempty"`
}

// MongoDBRoleCondition describes the state of a MongoDBRole at a certain point.
type MongoDBRoleCondition struct {
	// Type of MongoDBRole condition.
	Type string `json:"type,omitempty"`

	// Status of the condition, one of True, False, Unknown.
	Status core.ConditionStatus `json:"status,omitempty"`

	// The reason for the condition's.
	Reason string `json:"reason,omitempty"`

	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
}
