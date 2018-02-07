package resource

import (
	"github.com/appscode/go/log"
	"k8s.io/client-go/tools/cache"
)

func (ri *Indexer) EventHandler() cache.ResourceEventHandler {
	return ri
}

func (ri *Indexer) OnAdd(obj interface{}) {
	ri.cfgLock.RLock()
	defer ri.cfgLock.RUnlock()

	if !ri.enable {
		return
	}

	if err := ri.insert(obj); err != nil {
		log.Errorln(err)
	}
}

func (ri *Indexer) OnUpdate(oldObj, newObj interface{}) {
	ri.cfgLock.RLock()
	defer ri.cfgLock.RUnlock()

	if !ri.enable {
		return
	}

	if err := ri.insert(newObj); err != nil {
		log.Errorln(err)
	}
}

func (ri *Indexer) OnDelete(obj interface{}) {
	ri.cfgLock.RLock()
	defer ri.cfgLock.RUnlock()

	if !ri.enable {
		return
	}

	if err := ri.delete(obj); err != nil {
		log.Errorln(err)
		return
	}
}
