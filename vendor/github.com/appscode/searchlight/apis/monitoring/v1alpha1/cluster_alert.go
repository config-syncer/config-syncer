package v1alpha1

import (
	"fmt"
	"time"

	apiv1 "k8s.io/client-go/pkg/api/v1"
)

const (
	ResourceKindClusterAlert = "ClusterAlert"
	ResourceNameClusterAlert = "cluster-alert"
	ResourceTypeClusterAlert = "clusteralerts"
)

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

func (a ClusterAlert) IsValid() (bool, error) {
	cmd, ok := ClusterCommands[a.Spec.Check]
	if !ok {
		return false, fmt.Errorf("'%s' is not a valid cluster check command.", a.Spec.Check)
	}
	for k := range a.Spec.Vars {
		if _, ok := cmd.Vars[k]; !ok {
			return false, fmt.Errorf("Var '%s' is unsupported for check command %s.", k, a.Spec.Check)
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
			return false, fmt.Errorf("State '%s' is unsupported for check command %s.", rcv.State, a.Spec.Check)
		}
	}
	return true, nil
}

func (a ClusterAlert) GetNotifierSecretName() string {
	return a.Spec.NotifierSecretName
}

func (a ClusterAlert) GetReceivers() []Receiver {
	return a.Spec.Receivers
}

func (a ClusterAlert) ObjectReference() *apiv1.ObjectReference {
	return &apiv1.ObjectReference{
		APIVersion:      SchemeGroupVersion.String(),
		Kind:            ResourceKindClusterAlert,
		Namespace:       a.Namespace,
		Name:            a.Name,
		UID:             a.UID,
		ResourceVersion: a.ResourceVersion,
	}
}
