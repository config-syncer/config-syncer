package monitoring

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ClusterAlert struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the desired state of the ClusterAlert.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#spec-and-status
	Spec ClusterAlertSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterAlertList is a collection of ClusterAlert.
type ClusterAlertList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is the list of ClusterAlert.
	Items []ClusterAlert `json:"items"`
}

// ClusterAlertSpec describes the ClusterAlert the user wishes to create.
type ClusterAlertSpec struct {
	// Icinga CheckCommand name
	Check CheckCluster `json:"check,omitempty"`

	// How frequently Icinga Service will be checked
	CheckInterval metav1.Duration `json:"checkInterval,omitempty"`

	// How frequently notifications will be send
	AlertInterval metav1.Duration `json:"alertInterval,omitempty"`

	// Secret containing notifier credentials
	NotifierSecretName string `json:"notifierSecretName,omitempty"`

	// NotifierParams contains information to send notifications for Incident
	// State, UserUid, Method
	Receivers []Receiver `json:"receivers,omitempty"`

	// Vars contains Icinga Service variables to be used in CheckCommand
	Vars map[string]string `json:"vars,omitempty"`
}

type CheckCluster string

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type NodeAlert struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the desired state of the NodeAlert.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#spec-and-status
	Spec NodeAlertSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeAlertList is a collection of NodeAlert.
type NodeAlertList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is the list of NodeAlert.
	Items []NodeAlert `json:"items"`
}

// NodeAlertSpec describes the NodeAlert the user wishes to create.
type NodeAlertSpec struct {
	Selector map[string]string `json:"selector,omitempty"`

	NodeName string `json:"nodeName,omitempty"`

	// Icinga CheckCommand name
	Check CheckNode `json:"check,omitempty"`

	// How frequently Icinga Service will be checked
	CheckInterval metav1.Duration `json:"checkInterval,omitempty"`

	// How frequently notifications will be send
	AlertInterval metav1.Duration `json:"alertInterval,omitempty"`

	// Secret containing notifier credentials
	NotifierSecretName string `json:"notifierSecretName,omitempty"`

	// NotifierParams contains information to send notifications for Incident
	// State, UserUid, Method
	Receivers []Receiver `json:"receivers,omitempty"`

	// Vars contains Icinga Service variables to be used in CheckCommand
	Vars map[string]string `json:"vars,omitempty"`
}

type CheckNode string

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PodAlert struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the desired state of the PodAlert.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#spec-and-status
	Spec PodAlertSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodAlertList is a collection of PodAlert.
type PodAlertList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is the list of PodAlert.
	Items []PodAlert `json:"items"`
}

// PodAlertSpec describes the PodAlert the user wishes to create.
type PodAlertSpec struct {
	Selector metav1.LabelSelector `json:"selector,omitempty"`

	PodName string `json:"podName,omitempty"`

	// Icinga CheckCommand name
	Check CheckPod `json:"check,omitempty"`

	// How frequently Icinga Service will be checked
	CheckInterval metav1.Duration `json:"checkInterval,omitempty"`

	// How frequently notifications will be send
	AlertInterval metav1.Duration `json:"alertInterval,omitempty"`

	// Secret containing notifier credentials
	NotifierSecretName string `json:"notifierSecretName,omitempty"`

	// NotifierParams contains information to send notifications for Incident
	// State, UserUid, Method
	Receivers []Receiver `json:"receivers,omitempty"`

	// Vars contains Icinga Service variables to be used in CheckCommand
	Vars map[string]string `json:"vars,omitempty"`
}

type CheckPod string

type Receiver struct {
	// For which state notification will be sent
	State string `json:"state,omitempty"`

	// To whom notification will be sent
	To []string `json:"to,omitempty"`

	// How this notification will be sent
	Notifier string `json:"notifier,omitempty"`
}
