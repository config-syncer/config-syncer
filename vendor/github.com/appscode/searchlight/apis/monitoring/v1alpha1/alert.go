package v1alpha1

import (
	"time"

	apiv1 "k8s.io/client-go/pkg/api/v1"
)

type Alert interface {
	GetName() string
	GetNamespace() string
	Command() string
	GetCheckInterval() time.Duration
	GetAlertInterval() time.Duration
	IsValid() (bool, error)
	GetNotifierSecretName() string
	GetReceivers() []Receiver
	ObjectReference() *apiv1.ObjectReference
}
