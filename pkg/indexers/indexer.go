package indexers

import (
	"encoding/json"
	"path/filepath"
	"sync"

	"github.com/appscode/go/errors"
	"github.com/appscode/go/log"
	"github.com/blevesearch/bleve"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

type ResourceIndexer struct {
	// Full text indexer client
	index bleve.Index

	enable bool
	lock   sync.RWMutex
}

var _ cache.ResourceEventHandler = &ResourceIndexer{}

func NewResourceIndexer(dst string) (*ResourceIndexer, error) {
	c, err := ensureIndex(filepath.Join(dst, "resource.indexer"), "search")
	if err != nil {
		return nil, err
	}
	return &ResourceIndexer{
		index: c,
	}, nil
}

func (ri *ResourceIndexer) Configure(enable bool) {
	ri.lock.Lock()
	defer ri.lock.Unlock()
	ri.enable = enable
}

func (ri *ResourceIndexer) OnAdd(obj interface{}) {
	ri.lock.RLock()
	defer ri.lock.RUnlock()

	if !ri.enable {
		return
	}

	if err := ri.indexDocument(obj); err != nil {
		log.Errorln(err)
	}
}

func (ri *ResourceIndexer) OnDelete(obj interface{}) {
	ri.lock.RLock()
	defer ri.lock.RUnlock()

	if !ri.enable {
		return
	}

	key := keyFunction(obj)
	if err := ri.index.Delete(key); err != nil {
		log.Errorln(err)
		return
	}
	if err := ri.index.DeleteInternal([]byte(key)); err != nil {
		log.Errorln(err)
		return
	}
}

func (ri *ResourceIndexer) OnUpdate(oldObj, newObj interface{}) {
	ri.lock.RLock()
	defer ri.lock.RUnlock()

	if !ri.enable {
		return
	}

	if err := ri.indexDocument(newObj); err != nil {
		log.Errorln(err)
	}
}

func (ri *ResourceIndexer) indexDocument(obj interface{}) error {
	key := keyFunction(obj)
	err := ri.index.Index(key, obj)
	if err != nil {
		return errors.FromErr(err).WithMessage("Failed to index document").Err()
	}

	data, err := json.Marshal(obj)
	if err != nil {
		return errors.FromErr(err).WithMessage("Failed to marshal internal document").Err()
	}

	err = ri.index.SetInternal([]byte(key), data)
	if err != nil {
		return errors.FromErr(err).WithMessage("Failed store internal document").Err()
	}
	return nil
}

func keyFunction(obj interface{}) string {
	meta := metaAccessor(obj)
	if meta != nil {
		return string(meta.GetUID())
	}
	return ""
}

func metaAccessor(obj interface{}) metav1.Object {
	switch t := obj.(type) {
	case metav1.Object:
		return t
	case metav1.ObjectMetaAccessor:
		if m := t.GetObjectMeta(); m != nil {
			return m
		}
	}
	return nil
}
