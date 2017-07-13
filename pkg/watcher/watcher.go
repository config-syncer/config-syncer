package watcher

import (
	"sync"
	"time"

	"github.com/appscode/kubed/pkg/indexers"
	"github.com/appscode/kubed/pkg/recover"
	srch_cs "github.com/appscode/searchlight/client/clientset"
	scs "github.com/appscode/stash/client/clientset"
	vcs "github.com/appscode/voyager/client/clientset"
	pcm "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1alpha1"
	kcs "github.com/k8sdb/apimachinery/client/clientset"
	clientset "k8s.io/client-go/kubernetes"
)

type Controller struct {
	KubeClient        clientset.Interface
	VoyagerClient     vcs.ExtensionInterface
	SearchlightClient srch_cs.ExtensionInterface
	StashClient       scs.ExtensionInterface
	PromClient        pcm.MonitoringV1alpha1Interface
	KubeDBClient      kcs.ExtensionInterface

	Saver        *recover.RecoverStuff
	RunOptions   RunOptions
	Indexer      *indexers.ResourceIndexer
	ReverseIndex *indexers.ReverseIndexer
	SyncPeriod   time.Duration
	sync.Mutex
}

type RunOptions struct {
	Master                            string
	KubeConfig                        string
	ESEndpoint                        string
	InfluxSecretName                  string
	InfluxSecretNamespace             string
	ClusterName                       string
	ClusterKubedConfigSecretName      string
	ClusterKubedConfigSecretNamespace string
	Indexer                           string
	EnableReverseIndex                bool
	ServerAddress                     string
	NotifyOnCertSoonToBeExpired       bool
	NotifyVia                         string
}

func (c *Controller) Run() {
	c.WatchAlertmanagers()
	c.WatchClusterAlerts()
	c.WatchConfigMaps()
	c.WatchDaemonSets()
	c.WatchDeploymentApps()
	c.WatchDeploymentExtensions()
	c.WatchDormantDatabases()
	c.WatchElastics()
	c.WatchEvents()
	c.WatchIngresss()
	c.WatchJobs()
	c.watchNamespaces()
	c.WatchNodeAlerts()
	c.WatchPersistentVolumeClaims()
	c.WatchPersistentVolumes()
	c.WatchPodAlerts()
	c.WatchPostgreses()
	c.WatchPrometheuss()
	c.WatchReplicaSets()
	c.WatchReplicationControllers()
	c.WatchRestics()
	c.WatchSecrets()
	c.watchService()
	c.WatchServiceMonitors()
	c.WatchStatefulSets()
	c.WatchStorageClasss()
	c.WatchVoyagerCertificates()
	c.WatchVoyagerIngresses()
}
