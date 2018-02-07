package resource

import (
	"github.com/appscode/go/log"
	"k8s.io/client-go/tools/cache"
)

func (ri *ResourceIndexer) EventHandler() cache.ResourceEventHandler {
	return ri
}

func (ri *ResourceIndexer) OnAdd(obj interface{}) {
	if err := ri.insert(obj); err != nil {
		log.Errorln(err)
	}
}

func (ri *ResourceIndexer) OnUpdate(oldObj, newObj interface{}) {
	if err := ri.insert(newObj); err != nil {
		log.Errorln(err)
	}
}

func (ri *ResourceIndexer) OnDelete(obj interface{}) {
	if err := ri.delete(obj); err != nil {
		log.Errorln(err)
		return
	}
}
