package v1alpha1

import (
	"fmt"
	"strings"
	"time"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	ResourceKindNodeAlert     = "NodeAlert"
	ResourcePluralNodeAlert   = "nodealerts"
	ResourceSingularNodeAlert = "nodealert"
)

// +genclient
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
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

	NodeName *string `json:"nodeName,omitempty"`

	// Icinga CheckCommand name
	Check string `json:"check,omitempty"`

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

	// Indicates that Check is paused
	// Icinga Services are removed
	Paused bool `json:"paused,omitempty"`
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

func (a NodeAlert) IsValid(kc kubernetes.Interface) error {
	if a.Spec.Paused {
		return nil
	}

	if a.Spec.NodeName != nil && len(a.Spec.Selector) > 0 {
		return fmt.Errorf("can't specify both node name and selector")
	}

	cmd, ok := NodeCommands.Get(a.Spec.Check)
	if !ok {
		return fmt.Errorf("%s is not a valid node check command", a.Spec.Check)
	}

	if err := validateVariables(cmd.Vars, a.Spec.Vars); err != nil {
		return err
	}

	for _, rcv := range a.Spec.Receivers {
		found := false
		for _, state := range cmd.States {
			if strings.EqualFold(state, rcv.State) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("state %s is unsupported for check command %s", rcv.State, a.Spec.Check)
		}
	}

	return checkNotifiers(kc, a)
}

func (a NodeAlert) GetNotifierSecretName() string {
	return a.Spec.NotifierSecretName
}

func (a NodeAlert) GetReceivers() []Receiver {
	return a.Spec.Receivers
}

func (a NodeAlert) ObjectReference() *core.ObjectReference {
	return &core.ObjectReference{
		APIVersion:      SchemeGroupVersion.String(),
		Kind:            ResourceKindNodeAlert,
		Namespace:       a.Namespace,
		Name:            a.Name,
		UID:             a.UID,
		ResourceVersion: a.ResourceVersion,
	}
}
