/*
Copyright The Kubed Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	ClusterName           string `json:"clusterName,omitempty"`
	ConfigSourceNamespace string `json:"configSourceNamespace,omitempty"`
	EnableConfigSyncer    bool   `json:"enableConfigSyncer"`
	KubeConfigFile        string `json:"kubeConfigFile,omitempty"`
}
