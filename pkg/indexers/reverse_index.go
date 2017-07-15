package indexers

import (
	"path/filepath"

	"github.com/blevesearch/bleve"
	clientset "k8s.io/client-go/kubernetes"
)

type ReverseIndexer struct {
	// kubeClient to access kube api server
	kubeClient clientset.Interface
	index      bleve.Index

	Service ServiceIndexer
}

func NewReverseIndexer(cl clientset.Interface, dst string) (*ReverseIndexer, error) {
	index, err := ensureIndex(filepath.Join(dst, "reverse.indexer"), "indexer")
	if err != nil {
		return nil, err
	}
	ri := ReverseIndexer{
		kubeClient: cl,
		index:      index,
	}
	ri.Service = &ServiceIndexerImpl{kubeClient: cl, index: index}

	return &ri, nil
}
