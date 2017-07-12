package api

import "time"

type Alert interface {
	GetName() string
	GetNamespace() string
	Command() string
	GetCheckInterval() time.Duration
	GetAlertInterval() time.Duration
	IsValid() (bool, error)
}
