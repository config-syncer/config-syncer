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
	ResourceKindPodAlert     = "PodAlert"
	ResourcePluralPodAlert   = "podalerts"
	ResourceSingularPodAlert = "podalert"
)

// +genclient
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
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
	Selector *metav1.LabelSelector `json:"selector,omitempty"`

	PodName *string `json:"podName,omitempty"`

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

var _ Alert = &PodAlert{}

func (a PodAlert) GetName() string {
	return a.Name
}

func (a PodAlert) GetNamespace() string {
	return a.Namespace
}

func (a PodAlert) Command() string {
	return string(a.Spec.Check)
}

func (a PodAlert) GetCheckInterval() time.Duration {
	return a.Spec.CheckInterval.Duration
}

func (a PodAlert) GetAlertInterval() time.Duration {
	return a.Spec.AlertInterval.Duration
}

func (a PodAlert) IsValid(kc kubernetes.Interface) error {
	if a.Spec.Paused {
		return nil
	}

	if a.Spec.PodName != nil && a.Spec.Selector != nil {
		return fmt.Errorf("can't specify both pod name and selector")
	}
	if a.Spec.PodName == nil && a.Spec.Selector == nil {
		return fmt.Errorf("specify either pod name or selector")
	}
	if a.Spec.Selector != nil {
		_, err := metav1.LabelSelectorAsSelector(a.Spec.Selector)
		if err != nil {
			return err
		}
	}

	cmd, ok := PodCommands.Get(a.Spec.Check)
	if !ok {
		return fmt.Errorf("%s is not a valid pod check command", a.Spec.Check)
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

func (a PodAlert) GetNotifierSecretName() string {
	return a.Spec.NotifierSecretName
}

func (a PodAlert) GetReceivers() []Receiver {
	return a.Spec.Receivers
}

func (a PodAlert) ObjectReference() *core.ObjectReference {
	return &core.ObjectReference{
		APIVersion:      SchemeGroupVersion.String(),
		Kind:            ResourceKindPodAlert,
		Namespace:       a.Namespace,
		Name:            a.Name,
		UID:             a.UID,
		ResourceVersion: a.ResourceVersion,
	}
}
