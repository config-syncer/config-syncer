package operator

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

type Operator struct {
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

func (op *Operator) Run() {
	go op.WatchAlertmanagers()
	go op.WatchClusterAlerts()
	go op.WatchConfigMaps()
	go op.WatchDaemonSets()
	go op.WatchDeploymentApps()
	go op.WatchDeploymentExtensions()
	go op.WatchDormantDatabases()
	go op.WatchElastics()
	go op.WatchEvents()
	go op.WatchIngresss()
	go op.WatchJobs()
	go op.watchNamespaces()
	go op.WatchNodeAlerts()
	go op.WatchPersistentVolumeClaims()
	go op.WatchPersistentVolumes()
	go op.WatchPodAlerts()
	go op.WatchPostgreses()
	go op.WatchPrometheuss()
	go op.WatchReplicaSets()
	go op.WatchReplicationControllers()
	go op.WatchRestics()
	go op.WatchSecrets()
	go op.watchService()
	go op.WatchServiceMonitors()
	go op.WatchStatefulSets()
	go op.WatchStorageClasss()
	go op.WatchVoyagerCertificates()
	go op.WatchVoyagerIngresses()

	go op.StartCron()
}

func (op *Operator) StartCron() {
	op.Cron.Start()

	op.Cron.AddFunc("@every 24h", func() {
		janitor := influx.Janitor{Config: op.Config}
		janitor.CleanInflux()
	})
	op.Cron.AddFunc("@every 24h", func() {
		janitor := es.Janitor{Config: op.Config}
		janitor.CleanES()
	})
	op.Cron.AddFunc("@every 24h", func() {
		err := filepath.Walk(op.Config.RecycleBin.Path, func(path string, info os.FileInfo, err error) error {
			// delete old objects
			return nil
		})
		if err != nil {
			log.Errorln(err)
		}
		// expire saver
	})
	op.Cron.AddFunc(op.Config.ClusterSnapshot.Schedule, func() {
		if config, err := rest.InClusterConfig(); err == nil {
			err := backup.Backup(config, backup.Options{
				BackupDir: "/tmp/abc",
				Sanitize:  op.Config.ClusterSnapshot.Sanitize,
			})
			if err != nil {
				log.Errorln(err)
			}

			// upload to cloud
		}

		// run backup
	})
}
