package controller

import (
	"github.com/appscode/log"
	"github.com/appscode/searchlight/pkg/controller/host"
)

func (b *IcingaController) parseAlertOptions() {
	if b.ctx.Resource == nil {
		log.Infoln("Config is nil, nothing to parse")
		return
	}
	log.Infoln("Parsing labels.")
	b.ctx.ObjectType, b.ctx.ObjectName = host.GetObjectInfo(b.ctx.Resource.ObjectMeta.Labels)
}
