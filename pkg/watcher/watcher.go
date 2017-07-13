package watcher

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/appscode/kubed/pkg/backup"
	"github.com/appscode/kubed/pkg/config"
	"github.com/appscode/kubed/pkg/elasticsearch"
	"github.com/appscode/kubed/pkg/indexers"
	"github.com/appscode/kubed/pkg/influxdb"
	"github.com/appscode/kubed/pkg/recyclebin"
	"github.com/appscode/log"
	srch_cs "github.com/appscode/searchlight/client/clientset"
	scs "github.com/appscode/stash/client/clientset"
	vcs "github.com/appscode/voyager/client/clientset"
	pcm "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1alpha1"
	kcs "github.com/k8sdb/apimachinery/client/clientset"
	"gopkg.in/robfig/cron.v2"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Options struct {
	Master     string
	KubeConfig string

	ConfigPath         string
	ServerAddress      string
	Indexer            string
	EnableReverseIndex bool

	EnableConfigSync bool
	ScratchDir       string

	EnableAnalytics bool
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
	Saver        *recyclebin.RecoverStuff
	Indexer      *indexers.ResourceIndexer
	ReverseIndex *indexers.ReverseIndexer

	Cron       *cron.Cron
	SyncPeriod time.Duration
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

	go w.StartCron()
}

func (w *Watchers) StartCron() {
	w.Cron.Start()

	w.Cron.AddFunc("@every 24h", func() {
		janitor := influx.Janitor{Config: w.Config}
		janitor.CleanInflux()
	})
	w.Cron.AddFunc("@every 24h", func() {
		janitor := es.Janitor{Config: w.Config}
		janitor.CleanES()
	})
	w.Cron.AddFunc("@every 24h", func() {
		err := filepath.Walk(w.Config.RecycleBin.Path, func(path string, info os.FileInfo, err error) error {
			// delete old objects
			return nil
		})
		if err != nil {
			log.Errorln(err)
		}
		// expire saver
	})
	w.Cron.AddFunc(w.Config.Backup.Schedule, func() {
		if config, err := rest.InClusterConfig(); err == nil {
			err := backup.Backup(config, backup.Options{
				BackupDir: "/tmp/abc",
				Sanitize:  w.Config.Backup.Sanitize,
			})
			if err != nil {
				log.Errorln(err)
			}

			// upload to cloud
		}

		// run backup
	})
}
