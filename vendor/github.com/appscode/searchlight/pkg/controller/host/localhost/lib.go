package localhost

import (
	"github.com/appscode/errors"
	"github.com/appscode/searchlight/pkg/controller/host"
	"github.com/appscode/searchlight/pkg/controller/host/extpoints"
	"github.com/appscode/searchlight/pkg/controller/types"
)

func init() {
	extpoints.IcingaHostTypes.Register(new(icingaHost), host.HostTypeLocalhost)
}

type icingaHost struct {
}

type biblio struct {
	*types.Context
}

func (p *icingaHost) CreateAlert(ctx *types.Context, specificObject string) error {
	if specificObject != "" {
		return nil
	}
	return (&biblio{ctx}).create()
}

func (p *icingaHost) UpdateAlert(ctx *types.Context) error {
	return (&biblio{ctx}).update()
}

func (p *icingaHost) DeleteAlert(ctx *types.Context, specificObject string) error {
	return (&biblio{ctx}).delete()
}

//-----------------------------------------------------

// set Alert in Icinga LocalHost
func (b *biblio) create() error {
	alertSpec := b.Resource.Spec
	if alertSpec.CheckCommand == "" {
		return errors.New().WithMessage("Invalid request").BadRequest()
	}

	// Get Icinga Host Info
	objectList, err := host.GetObjectList(b.KubeClient, alertSpec.CheckCommand, host.HostTypeLocalhost, b.Resource.Namespace, b.ObjectType, b.ObjectName, "")
	if err != nil {
		return errors.New().WithCause(err).Internal()
	}

	var has bool
	if has, err = host.CheckIcingaService(b.IcingaClient, b.Resource.Name, objectList); err != nil {
		return errors.New().WithCause(err).Internal()
	}
	if has {
		return nil
	}

	// Create Icinga Host
	if err := host.CreateIcingaHost(b.IcingaClient, objectList, b.Resource.Namespace); err != nil {
		return errors.New().WithCause(err).Internal()
	}

	if err := b.createIcingaService(objectList); err != nil {
		return errors.New().WithCause(err).Internal()
	}

	if err := host.CreateIcingaNotification(b.IcingaClient, b.Resource, objectList); err != nil {
		return errors.New().WithCause(err).Internal()
	}

	return nil
}

func (b *biblio) createIcingaService(objectList []*host.KubeObjectInfo) error {
	alertSpec := b.Resource.Spec

	mp := make(map[string]interface{})
	mp["check_command"] = alertSpec.CheckCommand
	if alertSpec.IcingaParam != nil && alertSpec.IcingaParam.CheckIntervalSec > int64(0) {
		mp["check_interval"] = alertSpec.IcingaParam.CheckIntervalSec
	}

	commandVars := b.IcingaData[alertSpec.CheckCommand].VarInfo

	for key, val := range alertSpec.Vars {
		if _, found := commandVars[key]; found {
			mp[host.IVar(key)] = val
		}
	}

	return host.CreateIcingaService(b.IcingaClient, mp, objectList[0], b.Resource.Name)
}

func (b *biblio) update() error {
	alertSpec := b.Resource.Spec

	// Get Icinga Host Info
	objectList, err := host.GetObjectList(b.KubeClient, alertSpec.CheckCommand, host.HostTypeLocalhost, b.Resource.Namespace, b.ObjectType, b.ObjectName, "")
	if err != nil {
		return errors.New().WithCause(err).Internal()
	}

	if err := b.updateIcingaService(objectList); err != nil {
		return errors.New().WithCause(err).Internal()
	}

	if err := host.UpdateIcingaNotification(b.IcingaClient, b.Resource, objectList); err != nil {
		return errors.New().WithCause(err).Internal()
	}
	return nil
}

func (b *biblio) updateIcingaService(objectList []*host.KubeObjectInfo) error {
	alertSpec := b.Resource.Spec

	mp := make(map[string]interface{})
	if alertSpec.IcingaParam != nil && alertSpec.IcingaParam.CheckIntervalSec > int64(0) {
		mp["check_interval"] = alertSpec.IcingaParam.CheckIntervalSec
	}

	commandVars := b.IcingaData[alertSpec.CheckCommand].VarInfo
	for key, val := range alertSpec.Vars {
		if _, found := commandVars[key]; found {
			mp[host.IVar(key)] = val
		}
	}

	for _, object := range objectList {
		if err := host.UpdateIcingaService(b.IcingaClient, mp, object, b.Resource.Name); err != nil {
			return errors.New().WithCause(err).Internal()
		}
	}
	return nil
}

func (b *biblio) delete() error {
	alertSpec := b.Resource.Spec

	objectList, err := host.GetObjectList(b.KubeClient, alertSpec.CheckCommand, host.HostTypeLocalhost, b.Resource.Namespace, b.ObjectType, b.ObjectName, "")
	if err != nil {
		return errors.New().WithCause(err).Internal()
	}

	if err := host.DeleteIcingaService(b.IcingaClient, objectList, b.Resource.Name); err != nil {
		return errors.New().WithCause(err).Internal()
	}

	for _, object := range objectList {
		if err := host.DeleteIcingaHost(b.IcingaClient, object.Name); err != nil {
			return errors.New().WithCause(err).Internal()
		}
	}
	return nil
}
