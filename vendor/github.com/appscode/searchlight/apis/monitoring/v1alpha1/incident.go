package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindIncident     = "Incident"
	ResourcePluralIncident   = "incidents"
	ResourceSingularIncident = "incident"
)

// +genclient
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Incident struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Derived information about the incident.
	// +optional
	Status IncidentStatus `json:"status,omitempty"`
}

type IncidentStatus struct {
	// Type of last notification, such as problem, acknowledgement, recovery or custom
	LastNotificationType IncidentNotificationType `json:"lastNotificationType"`

	// Notifications for the incident, such as problem or acknowledgement.
	// +optional
	Notifications []IncidentNotification `json:"notifications,omitempty"`
}

type IncidentNotificationType string

// These are the possible notifications for an incident.
const (
	NotificationProblem         IncidentNotificationType = "Problem"
	NotificationAcknowledgement IncidentNotificationType = "Acknowledgement"
	NotificationRecovery        IncidentNotificationType = "Recovery"
	NotificationCustom          IncidentNotificationType = "Custom"
)

type IncidentNotification struct {
	// incident notification type.
	Type IncidentNotificationType `json:"type"`
	// brief output of check command for the incident
	// +optional
	CheckOutput string `json:"checkOutput"`
	// name of user making comment
	// +optional
	Author *string `json:"author,omitempty"`
	// comment made by user
	// +optional
	Comment *string `json:"comment,omitempty"`
	// The time at which this notification was first recorded. (Time of server receipt is in TypeMeta.)
	// +optional
	FirstTimestamp metav1.Time `json:"firstTimestamp,omitempty"`
	// The time at which the most recent occurrence of this notification was recorded.
	// +optional
	LastTimestamp metav1.Time `json:"lastTimestamp,omitempty"`
	// state of incident, such as Critical, Warning, OK, Unknown
	LastState string `json:"state"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IncidentList is a collection of Incident.
type IncidentList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is the list of Incident.
	Items []Incident `json:"items"`
}
