package v1alpha1

import (
	"github.com/appscode/go/encoding/json/types"
	"k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
)

const (
	ResourceKindMongoDBRole = "MongoDBRole"
	ResourceMongoDBRole     = "mongodbrole"
	ResourceMongoDBRoles    = "mongodbroles"

	ResourceKindMongoDBRoleBinding = "MongoDBRoleBinding"
	ResourceMongoDBRoleBinding     = "mongodbrolebinding"
	ResourceMongoDBRoleBindings    = "mongodbrolebindings"
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
	AuthManagerRef AuthManagerRef `json:"authManagerRef"`

	DatabaseRef appcat.AppReference `json:"databaseRef"`

	// links:
	// 	- https://www.vaultproject.io/api/secret/databases/index.html
	//	- https://www.vaultproject.io/api/secret/databases/mongodb.html

	// The name of the database connection to use for this role.
	DBName string `json:"dbName"`

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
	Status v1.ConditionStatus `json:"status,omitempty"`

	// The reason for the condition's.
	Reason string `json:"reason,omitempty"`

	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
}

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MongoDBRoleBinding binds mongodb credential to user
type MongoDBRoleBinding struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MongoDBRoleBindingSpec   `json:"spec,omitempty"`
	Status            MongoDBRoleBindingStatus `json:"status,omitempty"`
}

type MongoDBRoleBindingSpec struct {
	// Specifies the name of the MongoDBRole
	RoleRef string `json:"roleRef"`

	Subjects []rbacv1.Subject `json:"subjects"`

	Store Store `json:"store"`
}

type MongoDBRoleBindingPhase string

type MongoDBRoleBindingStatus struct {
	// observedGeneration is the most recent generation observed for this MongoDBRoleBinding. It corresponds to the
	// MongoDBRoleBinding's generation, which is updated on mutation by the API Server.
	ObservedGeneration *types.IntHash `json:"observedGeneration,omitempty"`

	// contains lease info of the credentials
	Lease LeaseData `json:"lease,omitempty"`

	// Specifies the phase of the MongoDBRoleBinding
	Phase MongoDBRoleBindingPhase `json:"phase,omitempty"`

	// Represents the latest available observations of a MongoDBRoleBinding current state.
	Conditions []MongoDBRoleBindingCondition `json:"conditions,omitempty"`
}

// MongoDBRoleBindingCondition describes the state of a MongoDBRoleBinding at a certain point.
type MongoDBRoleBindingCondition struct {
	// Type of MongoDBRoleBinding condition.
	Type string `json:"type,omitempty"`

	// Status of the condition, one of True, False, Unknown.
	Status v1.ConditionStatus `json:"status,omitempty"`

	// The reason for the condition's.
	Reason string `json:"reason,omitempty"`

	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type MongoDBRoleBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is a list of MongoDBRoleBinding objects
	Items []MongoDBRoleBinding `json:"items,omitempty"`
}
