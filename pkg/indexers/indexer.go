package indexers

import (
	"encoding/json"
	"path/filepath"

	"github.com/appscode/log"
	"github.com/blevesearch/bleve"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/appscode/errors"
)

type ResourceIndexer struct {
	// Full text indexer client
	client bleve.Index
}

func NewResourceIndexer(dst string) (*ResourceIndexer, error) {
	c, err := ensureIndex(filepath.Join(dst, "resource.indexer"), "search")
	if err != nil {
		return nil, err
	}
	return &ResourceIndexer{
		client: c,
	}, nil
}

func ensureIndex(dst, doctype string) (bleve.Index, error) {
	c, err := bleve.Open(dst)
	if err != nil {
		documentMapping := bleve.NewDocumentMapping()
		mapping := bleve.NewIndexMapping()
		mapping.AddDocumentMapping(doctype, documentMapping)
		c, err := bleve.New(dst, mapping)
		if err != nil {
			return nil, err
		}
		return c, nil
	}
	return c, nil
}

func (ri *ResourceIndexer) HandleAdd(obj interface{}) error {
	return ri.indexDocument(obj)
}

func (ri *ResourceIndexer) HandleDelete(obj interface{}) error {
	key := keyFunction(obj)
	err := ri.client.Delete(key)
	if err != nil {
		return err
	}
	return ri.client.DeleteInternal([]byte(key))
}

func (ri *ResourceIndexer) HandleUpdate(oldObj, newObj interface{}) error {
	return ri.indexDocument(newObj)
}

func (ri *ResourceIndexer) indexDocument(obj interface{}) error {
	key := keyFunction(obj)
	err := ri.client.Index(key, obj)
	if err != nil {
		return errors.FromErr(err).WithMessage("Failed to index document").Err()
	}

	internalData, err := json.Marshal(obj)
	if err != nil {
		return errors.FromErr(err).WithMessage("Failed to marshal internal document").Err()
	}

	err = ri.client.SetInternal([]byte(key), internalData)
	if err != nil {
		return errors.FromErr(err).WithMessage("Failed store internal document").Err()
	}
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
