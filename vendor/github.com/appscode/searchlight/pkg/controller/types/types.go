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
