package resource

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/blevesearch/bleve"
	"k8s.io/apimachinery/pkg/api/meta"
)

type Indexer struct {
	indices map[string]bleve.Index
	dir     string

	enable bool

	idxLock sync.RWMutex
	cfgLock sync.RWMutex
}

func NewIndexer(dir string) *Indexer {
	return &Indexer{
		indices: map[string]bleve.Index{},
		dir:     dir,
	}
}

func (ri *Indexer) Configure(enable bool) error {
	ri.cfgLock.Lock()
	defer ri.cfgLock.Unlock()

	ri.enable = enable
	return nil
}

func (ri *Indexer) indexFor(ns string) (bleve.Index, error) {
	ri.idxLock.RLock()
	if idx, ok := ri.indices[ns]; ok {
		ri.idxLock.RUnlock()
		return idx, nil
	}

	ri.idxLock.Lock()
	defer ri.idxLock.Unlock()

	indexDir := filepath.Join(ri.dir, ns)
	idx, err := bleve.Open(indexDir)
	if err != nil {
		mapping := bleve.NewIndexMapping()
		mapping.AddDocumentMapping("search", bleve.NewDocumentMapping())
		idx, err := bleve.New(indexDir, mapping)
		if err != nil {
			return nil, fmt.Errorf("failed to create index for namespace %s at dir: %s", ns, indexDir)
		}
		ri.indices[ns] = idx
		return idx, nil
	}
	return idx, nil
}

func (ri *Indexer) insert(obj interface{}) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return err
	}

	index, err := ri.indexFor(accessor.GetNamespace())
	if err != nil {
		return err
	}
	id := accessor.GetUID()

	err = index.Index(string(id), obj)
	if err != nil {
		return fmt.Errorf("failed to index id %s. Reason: %s", id, err)
	}

	data, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("failed to serialize to json document for id %s. Reason: %s", id, err)
	}

	err = index.SetInternal([]byte(id), data)
	if err != nil {
		return fmt.Errorf("failed to store document for id %s. Reason: %s", id, err)
	}
	return nil
}

func (ri *Indexer) delete(obj interface{}) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return err
	}

	index, err := ri.indexFor(accessor.GetNamespace())
	if err != nil {
		return err
	}
	id := accessor.GetUID()

	if err := index.Delete(string(id)); err != nil {
		return fmt.Errorf("failed to delete id %s. Reason: %s", id, err)
	}
	if err := index.DeleteInternal([]byte(id)); err != nil {
		return fmt.Errorf("failed to delete document for id %s. Reason: %s", id, err)
	}
	return nil
}
