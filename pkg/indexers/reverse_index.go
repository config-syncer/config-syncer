package indexers

import (
	"path/filepath"

	"github.com/appscode/kutil/meta"
	"github.com/blevesearch/bleve"
	prom "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"
	"k8s.io/client-go/kubernetes"
)

type ReverseIndexer struct {
	// kubeClient to access kube api server
	kubeClient kubernetes.Interface
	promClient prom.MonitoringV1Interface
	index      bleve.Index

	Service        ServiceIndexer
	ServiceMonitor ServiceMonitorIndexer
	Prometheus     PrometheusIndexer
}

func NewReverseIndexer(cl kubernetes.Interface, pm prom.MonitoringV1Interface, dst string) (*ReverseIndexer, error) {
	index, err := ensureIndex(filepath.Join(dst, "reverse.indexer"), "indexer")
	if err != nil {
		return nil, err
	}
	ri := ReverseIndexer{
		kubeClient: cl,
		promClient: pm,
		index:      index,
	}
	ri.Service = &ServiceIndexerImpl{kubeClient: cl, index: index}
	if meta.IsPreferredAPIResource(cl, prom.Group+"/"+prom.Version, prom.ServiceMonitorsKind) {
		// Add Indexer only if Server support this resource
		ri.ServiceMonitor = &ServiceMonitorIndexerImpl{kubeClient: cl, index: index}
	}
	if meta.IsPreferredAPIResource(cl, prom.Group+"/"+prom.Version, prom.PrometheusesKind) {
		// Add Indexer only if Server support this resource
		ri.Prometheus = &PrometheusIndexerImpl{kubeClient: cl, promClient: pm, index: index}
	}
	return &ri, nil
}
