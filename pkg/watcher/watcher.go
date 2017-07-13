package watcher

import (
	"sync"
	"time"

	"github.com/appscode/kubed/pkg/config"
	"github.com/appscode/kubed/pkg/indexers"
	"github.com/appscode/kubed/pkg/recover"
	srch_cs "github.com/appscode/searchlight/client/clientset"
	scs "github.com/appscode/stash/client/clientset"
	vcs "github.com/appscode/voyager/client/clientset"
	pcm "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1alpha1"
	kcs "github.com/k8sdb/apimachinery/client/clientset"
	clientset "k8s.io/client-go/kubernetes"
)

type Options struct {
	Master             string
	KubeConfig         string
	EnableAnalytics    bool
	Indexer            string
	EnableReverseIndex bool
	ServerAddress      string
	ConfigPath         string
}

type Watchers struct {
	KubeClient        clientset.Interface
	VoyagerClient     vcs.ExtensionInterface
	SearchlightClient srch_cs.ExtensionInterface
	StashClient       scs.ExtensionInterface
	PromClient        pcm.MonitoringV1alpha1Interface
	KubeDBClient      kcs.ExtensionInterface

	Opt          Options
	Config       config.ClusterConfig
	Saver        *recover.RecoverStuff
	Indexer      *indexers.ResourceIndexer
	ReverseIndex *indexers.ReverseIndexer
	SyncPeriod   time.Duration
	sync.Mutex
}

func (w *Watchers) Run() {
	go w.WatchAlertmanagers()
	go w.WatchClusterAlerts()
	go w.WatchConfigMaps()
	go w.WatchDaemonSets()
	go w.WatchDeploymentApps()
	go w.WatchDeploymentExtensions()
	go w.WatchDormantDatabases()
	go w.WatchElastics()
	go w.WatchEvents()
	go w.WatchIngresss()
	go w.WatchJobs()
	go w.watchNamespaces()
	go w.WatchNodeAlerts()
	go w.WatchPersistentVolumeClaims()
	go w.WatchPersistentVolumes()
	go w.WatchPodAlerts()
	go w.WatchPostgreses()
	go w.WatchPrometheuss()
	go w.WatchReplicaSets()
	go w.WatchReplicationControllers()
	go w.WatchRestics()
	go w.WatchSecrets()
	go w.watchService()
	go w.WatchServiceMonitors()
	go w.WatchStatefulSets()
	go w.WatchStorageClasss()
	go w.WatchVoyagerCertificates()
	go w.WatchVoyagerIngresses()
}
