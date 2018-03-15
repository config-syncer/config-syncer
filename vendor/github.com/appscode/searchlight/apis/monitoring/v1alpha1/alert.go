package v1alpha1

import (
	"time"

	core "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type Alert interface {
	GetName() string
	GetNamespace() string
	Command() string
	GetCheckInterval() time.Duration
	GetAlertInterval() time.Duration
	IsValid(kc kubernetes.Interface) error
	GetNotifierSecretName() string
	GetReceivers() []Receiver
	ObjectReference() *core.ObjectReference
}
