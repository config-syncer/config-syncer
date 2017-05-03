package controller

import (
	"github.com/appscode/errors"
	"github.com/appscode/log"
	"github.com/appscode/searchlight/pkg/controller/host/extpoints"
)

func (b *IcingaController) Create(specificObject ...string) error {
	if !b.checkIcingaAvailability() {
		return errors.New("Icinga is down").External()
	}

	log.Debugln("Starting createing alert", b.ctx.Resource.ObjectMeta)

	object := ""
	if len(specificObject) > 0 {
		object = specificObject[0]
	}

	alertSpec := b.ctx.Resource.Spec
	command, found := b.ctx.IcingaData[alertSpec.CheckCommand]
	if !found {
		return errors.New().
			WithMessagef("check_command [%s] not found", alertSpec.CheckCommand).
			InvalidData()
	}
	hostType, found := command.HostType[b.ctx.ObjectType]
	if !found {
		return errors.New().
			WithMessagef("check_command [%s] is not applicable to %s", alertSpec.CheckCommand, b.ctx.ObjectType).
			InvalidData()
	}
	p := extpoints.IcingaHostTypes.Lookup(hostType)
	if p == nil {
		return errors.New().
			WithMessagef("IcingaHostType %v is unknown", hostType).
			NotFound()
	}
	return p.CreateAlert(b.ctx, object)
}
