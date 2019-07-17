package v1alpha1

import (
	stringz "github.com/appscode/go/strings"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	store "kmodules.xyz/objectstore-api/api/v1"
)

const (
	TimestampFormat    = "20060102T150405"
	ConfigSyncKey      = "kubed.appscode.com/sync"
	ConfigOriginKey    = "kubed.appscode.com/origin"
	ConfigSyncContexts = "kubed.appscode.com/sync-contexts"

	JanitorElasticsearch = "Elasticsearch"
	JanitorInfluxDB      = "InfluxDB"

	OriginNameLabelKey      = "kubed.appscode.com/origin.name"
	OriginNamespaceLabelKey = "kubed.appscode.com/origin.namespace"
	OriginClusterLabelKey   = "kubed.appscode.com/origin.cluster"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ClusterConfig struct {
	metav1.TypeMeta `json:",inline"`

	ClusterName           string              `json:"clusterName,omitempty"`
	ConfigSourceNamespace string              `json:"configSourceNamespace,omitempty"`
	Snapshotter           *SnapshotSpec       `json:"snapshotter,omitempty"`
	RecycleBin            *RecycleBinSpec     `json:"recycleBin,omitempty"`
	EventForwarder        *EventForwarderSpec `json:"eventForwarder,omitempty"`
	EnableConfigSyncer    bool                `json:"enableConfigSyncer"`
	NotifierSecretName    string              `json:"notifierSecretName,omitempty"`
	Janitors              []JanitorSpec       `json:"janitors,omitempty"`
	KubeConfigFile        string              `json:"kubeConfigFile,omitempty"`
}

type JanitorSpec struct {
	Kind          string             `json:"kind"`
	TTL           metav1.Duration    `json:"ttl"`
	Elasticsearch *ElasticsearchSpec `json:"elasticsearch,omitempty,omitempty"`
	InfluxDB      *InfluxDBSpec      `json:"influxdb,omitempty"`
}

type ElasticsearchSpec struct {
	Endpoint       string `json:"endpoint,omitempty"`
	LogIndexPrefix string `json:"logIndexPrefix,omitempty"`
	SecretName     string `json:"secretName,omitempty"`
}

type InfluxDBSpec struct {
	Endpoint string `json:"endpoint,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type RecycleBinSpec struct {
	Path          string          `json:"path,omitempty"`
	TTL           metav1.Duration `json:"ttl,omitempty"`
	HandleUpdates bool            `json:"handleUpdates"`
}

type EventForwarderSpec struct {
	// Deprecated
	NodeAdded *ForwarderSpec `json:"nodeAdded,omitempty"`
	// Deprecated
	StorageAdded *ForwarderSpec `json:"storageAdded,omitempty"`
	// Deprecated
	IngressAdded *ForwarderSpec `json:"ingressAdded,omitempty"`
	// Deprecated
	WarningEvents *ForwarderSpec `json:"warningEvents,omitempty"`
	// Deprecated
	CSREvents *ForwarderSpec `json:"csrEvents,omitempty"`

	Rules     []PolicyRule `json:"rules"`
	Receivers []Receiver   `json:"receivers,omitempty"`
}

type PolicyRule struct {
	// Operation is the operation being performed
	Operations []Operation `json:"operations"`

	// Resources that this rule matches. An empty list implies all kinds in all API groups.
	// +optional
	Resources []GroupResources `json:"resources,omitempty"`

	// Namespaces that this rule matches.
	// The empty string "" matches non-namespaced resources.
	// An empty list implies every namespace.
	// +optional
	Namespaces []string `json:"namespaces,omitempty"`
}

// GroupResources represents resource kinds in an API group.
type GroupResources struct {
	// Group is the name of the API group that contains the resources.
	// The empty string represents the core API group.
	// +optional
	Group string `json:"group"`
	// Resources is a list of resources within the API group. Subresources are
	// matched using a "/" to indicate the subresource. For example, "pods/log"
	// would match request to the log subresource of pods. The top level resource
	// does not match subresources, "pods" doesn't match "pods/log".
	// +optional
	Resources []string `json:"resources,omitempty"`
	// ResourceNames is a list of resource instance names that the policy matches.
	// Using this field requires Resources to be specified.
	// An empty list implies that every instance of the resource is matched.
	// +optional
	ResourceNames []string `json:"resourceNames,omitempty"`
}

// Operation is the type of resource operation being checked for admission control
type Operation string

// Operation constants
const (
	Create Operation = "CREATE"
	Delete Operation = "DELETE"
)

type NoNamespacedForwarderSpec struct {
	Namespaces []string `json:"namespaces,omitempty"`
}

type ForwarderSpec struct {
	Handle     bool     `json:"handle"`
	Namespaces []string `json:"namespaces,omitempty"`
}

func (f ForwarderSpec) IsAllowed(ns string) bool {
	return len(f.Namespaces) == 0 || stringz.Contains(f.Namespaces, ns)
}

type Receiver struct {
	// To whom notification will be sent
	To []string `json:"to,omitempty"`

	// How this notification will be sent
	Notifier string `json:"notifier,omitempty"`
}

// For periodic full cluster backup
// https://github.com/appscode/kubed/issues/16
type SnapshotSpec struct {
	Schedule      string `json:"schedule,omitempty"`
	Sanitize      bool   `json:"sanitize,omitempty"`
	Overwrite     bool   `json:"overwrite,omitempty"`
	store.Backend `json:",inline,omitempty"`
}

type JanitorAuthInfo struct {
	CACertData         []byte `envconfig:"CA_CERT_DATA"`
	ClientCertData     []byte `envconfig:"CLIENT_CERT_DATA"`
	ClientKeyData      []byte `envconfig:"CLIENT_KEY_DATA"`
	InsecureSkipVerify bool   `envconfig:"INSECURE_SKIP_VERIFY"`

	Username string `envconfig:"USERNAME"`
	Password string `envconfig:"PASSWORD"`
	Token    string `envconfig:"TOKEN"`
}

type KubedMetadata struct {
	OperatorNamespace string `json:"operatorNamespace,omitempty"`
	SearchEnabled     bool   `json:"searchEnabled"`
}

// +genclient
// +genclient:onlyVerbs=get
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SearchResult struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Hits     []ResultEntry   `json:"hits,omitempty"`
	Total    uint64          `json:"totalHits"`
	MaxScore float64         `json:"maxScore"`
	Took     metav1.Duration `json:"took"`
}

var _ runtime.Object = &SearchResult{}

type ResultEntry struct {
	Score  float64              `json:"score"`
	Object runtime.RawExtension `json:"object,omitempty"`
}
