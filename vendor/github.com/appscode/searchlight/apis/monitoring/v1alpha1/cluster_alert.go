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
	ResourceKindClusterAlert     = "ClusterAlert"
	ResourcePluralClusterAlert   = "clusteralerts"
	ResourceSingularClusterAlert = "clusteralert"
)

// +genclient
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
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

var _ Alert = &ClusterAlert{}

func (a ClusterAlert) GetName() string {
	return a.Name
}

func (a ClusterAlert) GetNamespace() string {
	return a.Namespace
}

func (a ClusterAlert) Command() string {
	return string(a.Spec.Check)
}

func (a ClusterAlert) GetCheckInterval() time.Duration {
	return a.Spec.CheckInterval.Duration
}

func (a ClusterAlert) GetAlertInterval() time.Duration {
	return a.Spec.AlertInterval.Duration
}

func (a ClusterAlert) IsValid(kc kubernetes.Interface) error {
	if a.Spec.Paused {
		return nil
	}

	cmd, ok := ClusterCommands.Get(a.Spec.Check)
	if !ok {
		return fmt.Errorf("'%s' is not a valid cluster check command", a.Spec.Check)
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
			return fmt.Errorf("state '%s' is unsupported for check command %s", rcv.State, a.Spec.Check)
		}
	}

	return checkNotifiers(kc, a)
}

func (a ClusterAlert) GetNotifierSecretName() string {
	return a.Spec.NotifierSecretName
}

func (a ClusterAlert) GetReceivers() []Receiver {
	return a.Spec.Receivers
}

func (a ClusterAlert) ObjectReference() *core.ObjectReference {
	return &core.ObjectReference{
		APIVersion:      SchemeGroupVersion.String(),
		Kind:            ResourceKindClusterAlert,
		Namespace:       a.Namespace,
		Name:            a.Name,
		UID:             a.UID,
		ResourceVersion: a.ResourceVersion,
	}
}
