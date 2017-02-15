package types

import (
	"sync"

	aci "github.com/appscode/k8s-addons/api"
	acs "github.com/appscode/k8s-addons/client/clientset"
	"github.com/appscode/k8s-addons/pkg/stash"
	"github.com/appscode/searchlight/data"
	"github.com/appscode/searchlight/pkg/client/icinga"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
)

type EventReason string

const (
	// Icinga objects create event list
	CreatingIcingaObjects       EventReason = "Creating"
	FailedToCreateIcingaObjects EventReason = "FailedToCreate"
	NoIcingaObjectCreated       EventReason = "NoIcingaObjectCreated"
	CreatedIcingaObjects        EventReason = "Created"

	// Icinga objects update event list
	UpdatingIcingaObjects       EventReason = "Updating"
	FailedToUpdateIcingaObjects EventReason = "FailedToUpdate"
	UpdatedIcingaObjects        EventReason = "Updated"

	// Icinga objects delete event list
	DeletingIcingaObjects       EventReason = "Deleting"
	FailedToDeleteIcingaObjects EventReason = "FailedToDelete"
	DeletedIcingaObjects        EventReason = "Deleted"

	// Icinga objects sync event list
	SyncIcingaObjects         EventReason = "Sync"
	FailedToSyncIcingaObjects EventReason = "FailedToSync"
	SyncedIcingaObjects       EventReason = "Synced"
)

func (r EventReason) String() string {
	return string(r)
}

const (
	AcknowledgeTimestamp string = "acknowledgement_timestamp"
)

type IcingaData struct {
	HostType map[string]string
	VarInfo  map[string]data.CommandVar
}

type Context struct {
	// kubernetes client
	KubeClient              clientset.Interface
	AppsCodeExtensionClient acs.AppsCodeExtensionInterface

	IcingaClient *icinga.IcingaClient
	IcingaData   map[string]*IcingaData

	Resource   *aci.Alert
	ObjectType string
	ObjectName string

	Storage *stash.Storage
	sync.Mutex
}

type KubeOptions struct {
	ObjectType string
	ObjectName string
}

type Ancestors struct {
	Type  string   `json:"type,omitempty"`
	Names []string `json:"names,omitempty"`
}

type AlertEventMessage struct {
	IncidentEventId int64  `json:"incident_event_id,omitempty"`
	Comment         string `json:"comment,omitempty"`
	UserName        string `json:"username,omitempty"`
}
