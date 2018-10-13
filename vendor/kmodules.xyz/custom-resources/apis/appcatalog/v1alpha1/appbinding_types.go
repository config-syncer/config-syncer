package v1alpha1

import (
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	ResourceKindApp = "AppBinding"
	ResourceApps    = "appbindings"
	ResourceApp     = "appbinding"
)

// +genclient
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppBinding defines a generic user application.
type AppBinding struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              AppBindingSpec `json:"spec,omitempty"`
}

// AppBindingSpec is the spec for app
type AppBindingSpec struct {
	// Used to facilitate programmatic handling of application.
	// +optional
	Type AppType `json:"type,omitempty"`

	// ClientConfig defines how to communicate with the app.
	// Required
	ClientConfig ClientConfig `json:"clientConfig"`

	// Secret is the name of the secret to create in the AppBinding's
	// namespace that will hold the credentials associated with the AppBinding.
	Secret *core.LocalObjectReference `json:"secret,omitempty"`

	// Parameters is a set of the parameters to be used to connect to the
	// app. The inline YAML/JSON payload to be translated into equivalent
	// JSON object.
	//
	// The Parameters field is NOT secret or secured in any way and should
	// NEVER be used to hold sensitive information. To set parameters that
	// contain secret information, you should ALWAYS store that information
	// in a Secret.
	//
	// +optional
	Parameters *runtime.RawExtension `json:"parameters,omitempty"`
}

type AppType string

const (
	// AppTypeOpaque is the default. A generic application.
	AppTypeOpaque AppType = "Opaque"
)

// ClientConfig contains the information to make a connection with an app
type ClientConfig struct {
	// `url` gives the location of the app, in standard URL form
	// (`[scheme://]host:port/path`). Exactly one of `url` or `service`
	// must be specified.
	//
	// The `host` should not refer to a service running in the cluster; use
	// the `service` field instead. The host might be resolved via external
	// DNS in some apiservers (e.g., `kube-apiserver` cannot resolve
	// in-cluster DNS as that would be a layering violation). `host` may
	// also be an IP address.
	//
	// A path is optional, and if present may be any string permissible in
	// a URL. You may use the path to pass an arbitrary string to the
	// app, for example, a cluster identifier.
	//
	// Attempting to use a user or basic auth e.g. "user:password@" is not
	// allowed. Fragments ("#...") and query parameters ("?...") are not
	// allowed, either.
	//
	// +optional
	URL *string `json:"url,omitempty"`

	// `service` is a reference to the service for this app. Either
	// `service` or `url` must be specified.
	//
	// If the webhook is running within the cluster, then you should use `service`.
	//
	// +optional
	Service *ServiceReference `json:"service"`

	// InsecureSkipTLSVerify disables TLS certificate verification when communicating with this app.
	// This is strongly discouraged.  You should use the CABundle instead.
	InsecureSkipTLSVerify bool `json:"insecureSkipTLSVerify,omitempty"`

	// CABundle is a PEM encoded CA bundle which will be used to validate the serving certificate of this app.
	// +optional
	CABundle []byte `json:"caBundle,omitempty"`
}

// ServiceReference holds a reference to Service.legacy.k8s.io
type ServiceReference struct {
	// Specifies which scheme to use, for example: http, https
	// If specified, then it will applied as prefix in this format: scheme://
	// If not specified, then nothing will be prefixed
	Scheme string `json:"scheme"`

	// `name` is the name of the service.
	// Required
	Name string `json:"name"`

	// The port that will be exposed by this app.
	Port int32 `json:"port"`

	// `path` is an optional URL path which will be sent in any request to
	// this service.
	// +optional
	Path *string `json:"path,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppBindingList is a list of Apps
type AppBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of AppBinding CRD objects
	Items []AppBinding `json:"items,omitempty"`
}

type AppReference struct {
	// `namespace` is the namespace of the app.
	// Required
	Namespace string `json:"namespace"`

	// `name` is the name of the app.
	// Required
	Name string `json:"name"`

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

type AppBindingMeta interface {
	Name() string
	Type() AppType
}
