package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

type AuthManagerType string

const (
	AuthManagerTypeVault AuthManagerType = "Vault"
)

type AuthManagerRef struct {
	Type AuthManagerType `json:"type"`

	// Optional AppRef fields

	// `namespace` is the namespace of the app.
	// +optional
	Namespace *string `json:"namespace,omitempty"`

	// `name` is the name of the app.
	// +optional
	Name *string `json:"name,omitempty"`

	// Parameters is a set of the parameters to be used to override default
	// parameters. The inline YAML/JSON payload to be translated into equivalent
	// JSON object.
	//
	// The Parameters field is NOT secret or secured in any way and should
	// NEVER be used to hold sensitive information.
	//
	// +optional
	Parameters *runtime.RawExtension `json:"parameters,omitempty"`
}

// Store specifies where to store credentials
type Store struct {
	// Specifies the name of the secret
	Secret string `json:"secret"`
}

// LeaseData contains lease info
type LeaseData struct {
	// lease id
	ID string `json:"id,omitempty"`

	// lease duration in seconds
	Duration int64 `json:"duration,omitempty"`

	// lease renew deadline in Unix time
	RenewDeadline int64 `json:"renewDeadline,omitempty"`
}
