package indexers

import (
	"github.com/blevesearch/bleve"
)

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
