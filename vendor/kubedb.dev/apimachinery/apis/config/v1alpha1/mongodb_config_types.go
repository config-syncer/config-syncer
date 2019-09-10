package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	ResourceKindMongoConfiguration = "MongoConfiguration"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MongoDBConfiguration defines a MongoDB app configuration.
// https://www.vaultproject.io/api/secret/databases/index.html
// https://www.vaultproject.io/api/secret/databases/mongodb.html#configure-connection
type MongoDBConfiguration struct {
	metav1.TypeMeta `json:",inline,omitempty"`

	// ConfigServer is the dsn of config server of mongodb sharding. The dsn includes the port no too.
	ConfigServer string `json:"configServer,omitempty"`

	// ReplicaSets contains the dsn of each replicaset of sharding. The DSNs are in key-value pair, where
	// the keys are host-0, host-1 etc, and the values are DSN of each replicaset. If there is no sharding
	// but only one replicaset, then ReplicaSets field contains only one key-value pair where the key is
	// host-0 and the value is dsn of that replicaset.
	ReplicaSets map[string]string `json:"replicaSets,omitempty"`
}
