package api

import (
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindNodeAlert = "NodeAlert"
	ResourceNameNodeAlert = "node-alert"
	ResourceTypeNodeAlert = "nodealerts"
)

// NodeAlert types for appscode.
type NodeAlert struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the desired state of the NodeAlert.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#spec-and-status
	Spec NodeAlertSpec `json:"spec,omitempty"`

	// Status is the current state of the NodeAlert.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#spec-and-status
	// Status AlertStatus `json:"status,omitempty"`
}

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

	// NotifierParams contains information to send notifications for Incident
	// State, UserUid, Method
	Receivers []Receiver `json:"receivers,omitempty"`

	// Vars contains Icinga Service variables to be used in CheckCommand
	Vars map[string]interface{} `json:"vars,omitempty"`
}

var _ Alert = &NodeAlert{}

func (a NodeAlert) GetName() string {
	return a.Name
}

func (a NodeAlert) GetNamespace() string {
	return a.Namespace
}

func (a NodeAlert) Command() string {
	return string(a.Spec.Check)
}

func (a NodeAlert) GetCheckInterval() time.Duration {
	return a.Spec.CheckInterval.Duration
}

func (a NodeAlert) GetAlertInterval() time.Duration {
	return a.Spec.AlertInterval.Duration
}

func (a NodeAlert) IsValid() (bool, error) {
	cmd, ok := NodeCommands[a.Spec.Check]
	if !ok {
		return false, fmt.Errorf("%s is not a valid node check command.", a.Spec.Check)
	}
	for k := range a.Spec.Vars {
		if _, ok := cmd.Vars[k]; !ok {
			return false, fmt.Errorf("Var %s is unsupported for check command %s.", k, a.Spec.Check)
		}
	}
	for _, rcv := range a.Spec.Receivers {
		found := false
		for _, state := range cmd.States {
			if state == rcv.State {
				found = true
				break
			}
		}
		if !found {
			return false, fmt.Errorf("State %s is unsupported for check command %s.", rcv.State, a.Spec.Check)
		}
	}
	return true, nil
}

func (a NodeAlert) GetReceivers() []Receiver {
	return a.Spec.Receivers
}
