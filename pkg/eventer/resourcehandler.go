package eventer

import (
	"time"

	"github.com/appscode/go/log"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

const (
	MaxSyncInterval = 10 * time.Minute
)

var _ cache.ResourceEventHandler = &EventForwarder{}

func (f *EventForwarder) OnAdd(obj interface{}) {
	f.lock.RLock()
	defer f.lock.RUnlock()

	if f.spec == nil {
		return
	}

	accessor, err := meta.Accessor(obj)
	if err != nil {
		log.Errorln(err)
		return
	}

	if !recentEvent(accessor.GetCreationTimestamp()) {
		return
	}

}

func (f *EventForwarder) OnUpdate(oldObj, newObj interface{}) {
	f.lock.RLock()
	defer f.lock.RUnlock()

	if f.spec == nil {
		return
	}

	accessor, err := meta.Accessor(newObj)
	if err != nil {
		log.Errorln(err)
		return
	}

	if !recentEvent(accessor.GetCreationTimestamp()) {
		return
	}

}

func (f *EventForwarder) OnDelete(obj interface{}) {
	f.lock.RLock()
	defer f.lock.RUnlock()

	if f.spec == nil {
		return
	}

	accessor, err := meta.Accessor(obj)
	if err != nil {
		log.Errorln(err)
		return
	}

	if !recentEvent(accessor.GetCreationTimestamp()) {
		return
	}

}

func recentEvent(t metav1.Time) bool {
	return time.Now().Sub(t.Time) < MaxSyncInterval
}
