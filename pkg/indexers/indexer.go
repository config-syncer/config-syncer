package indexers

import (
	"encoding/json"
	"strings"

	"github.com/appscode/log"
	"github.com/blevesearch/bleve"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ResourceIndexer struct {
	// Full text indexer client
	client bleve.Index
}

func NewResourceIndexer(dst string) (*ResourceIndexer, error) {
	c, err := ensureIndex(strings.TrimRight(dst, "/")+"/resource.indexer", "search")
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

func (ri *ResourceIndexer) HandleAdd(obj interface{}) {
	ri.indexDocument(obj)
}

func (ri *ResourceIndexer) HandleDelete(obj interface{}) {
	key := keyFunction(obj)
	ri.client.Delete(key)
	ri.client.DeleteInternal([]byte(key))
}

func (ri *ResourceIndexer) HandleUpdate(oldObj, newObj interface{}) {
	ri.indexDocument(newObj)
}

func (ri *ResourceIndexer) indexDocument(obj interface{}) {
	key := keyFunction(obj)
	err := ri.client.Index(key, obj)
	if err != nil {
		log.Errorln("Failed to index document", err)
		return
	}

	internalData, err := json.Marshal(obj)
	if err != nil {
		log.Errorln("Failed to marshal internal document")
		return
	}

	err = ri.client.SetInternal([]byte(key), internalData)
	if err != nil {
		log.Errorln("Failed store internal document")
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
