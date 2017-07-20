package indexers

import (
	"path/filepath"

	"github.com/appscode/kubed/pkg/util"
	"github.com/blevesearch/bleve"
	prom "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1alpha1"
	clientset "k8s.io/client-go/kubernetes"
)

type ReverseIndexer struct {
	// kubeClient to access kube api server
	kubeClient clientset.Interface
	index      bleve.Index

	Service        ServiceIndexer
	ServiceMonitor ServiceMonitorIndexer
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
	if util.IsPreferredAPIResource(cl, prom.TPRGroup+"/"+prom.TPRVersion, prom.TPRServiceMonitorsKind) {
		// Add Indexer only if Server support this resource
		ri.ServiceMonitor = &ServiceMonitorIndexerImpl{kubeClient: cl, index: index}
	}
	return &ri, nil
}
