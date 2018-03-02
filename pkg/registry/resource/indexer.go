package resource

import (
	"path/filepath"
	"sync"

	"github.com/blevesearch/bleve"
	"github.com/json-iterator/go"
	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type ResourceIndexer struct {
	indices map[string]bleve.Index
	dir     string

	idxLock sync.Mutex
}

func NewIndexer(dir string) *ResourceIndexer {
	return &ResourceIndexer{
		indices: map[string]bleve.Index{},
		dir:     dir,
	}
}

func (ri *ResourceIndexer) indexFor(ns string) (bleve.Index, error) {
	ri.idxLock.Lock()
	defer ri.idxLock.Unlock()

	if idx, ok := ri.indices[ns]; ok {
		return idx, nil
	}

	indexDir := filepath.Join(ri.dir, ns)
	idx, err := bleve.Open(indexDir)
	if err != nil {
		mapping := bleve.NewIndexMapping()
		mapping.AddDocumentMapping("search", bleve.NewDocumentMapping())
		idx, err := bleve.New(indexDir, mapping)
		if err != nil {
			return nil, errors.Errorf("failed to create index for namespace %s at dir: %s", ns, indexDir)
		}
		ri.indices[ns] = idx
		return idx, nil
	}
	return idx, nil
}

func (ri *ResourceIndexer) insert(obj interface{}) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return err
	}

	if annotations := accessor.GetAnnotations(); annotations != nil {
		delete(annotations, "kubectl.kubernetes.io/last-applied-configuration")
		accessor.SetAnnotations(annotations)
	}

	if s, ok := obj.(*core.Secret); ok {
		for key := range s.Data {
			s.Data[key] = []byte("")
		}
	}

	index, err := ri.indexFor(accessor.GetNamespace())
	if err != nil {
		return err
	}
	id := accessor.GetUID()

	err = index.Index(string(id), obj)
	if err != nil {
		return errors.Errorf("failed to index id %s. Reason: %s", id, err)
	}

	data, err := json.Marshal(obj)
	if err != nil {
		return errors.Errorf("failed to serialize to json document for id %s. Reason: %s", id, err)
	}

	err = index.SetInternal([]byte(id), data)
	if err != nil {
		return errors.Errorf("failed to store document for id %s. Reason: %s", id, err)
	}
	return nil
}

func (ri *ResourceIndexer) delete(obj interface{}) error {
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
		return errors.Errorf("failed to delete id %s. Reason: %s", id, err)
	}
	if err := index.DeleteInternal([]byte(id)); err != nil {
		return errors.Errorf("failed to delete document for id %s. Reason: %s", id, err)
	}
	return nil
}
