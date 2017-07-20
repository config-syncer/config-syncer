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
	promClient prom.MonitoringV1alpha1Interface
	index      bleve.Index

	Service        ServiceIndexer
	ServiceMonitor ServiceMonitorIndexer
	Prometheus     PrometheusIndexer
}

func NewReverseIndexer(cl clientset.Interface, pm prom.MonitoringV1alpha1Interface, dst string) (*ReverseIndexer, error) {
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
	if util.IsPreferredAPIResource(cl, prom.TPRGroup+"/"+prom.TPRVersion, prom.TPRServiceMonitorsKind) {
		// Add Indexer only if Server support this resource
		ri.ServiceMonitor = &ServiceMonitorIndexerImpl{kubeClient: cl, index: index}
	}
	if util.IsPreferredAPIResource(cl, prom.TPRGroup+"/"+prom.TPRVersion, prom.TPRPrometheusesKind) {
		// Add Indexer only if Server support this resource
		ri.Prometheus = &PrometheusIndexerImpl{kubeClient: cl, promClient: pm, index: index}
	}
	return &ri, nil
}
