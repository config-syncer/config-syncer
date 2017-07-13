package operator

import (
	"net/http"
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
	"github.com/appscode/kubed/pkg/storage"
	"github.com/appscode/log"
	"github.com/appscode/pat"
	srch_cs "github.com/appscode/searchlight/client/clientset"
	scs "github.com/appscode/stash/client/clientset"
	vcs "github.com/appscode/voyager/client/clientset"
	shell "github.com/codeskyblue/go-sh"
	pcm "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1alpha1"
	kcs "github.com/k8sdb/apimachinery/client/clientset"
	"gopkg.in/robfig/cron.v2"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Options struct {
	Master     string
	KubeConfig string

	ConfigPath         string
	ServerAddress      string
	Indexer            string
	EnableReverseIndex bool

	EnableConfigSync  bool
	ScratchDir        string
	OperatorNamespace string

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

func (op *Operator) Setup() error {
	cfg, err := config.LoadConfig(op.Opt.ConfigPath)
	if err != nil {
		return err
	}
	op.Config = *cfg

	op.Cron = cron.New()
	op.Cron.Start()
	//Saver: &recyclebin.RecoverStuff{
	//	Opt: *cfg.RecycleBin,
	//},
	op.SyncPeriod = time.Minute * 2

	return nil
}

func (op *Operator) RunWatchers() {
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
}

func (op *Operator) ListenAndServe() {
	// router is default HTTP request multiplexer for kubed. It matches the URL of each
	// incoming request against a list of registered patterns with their associated
	// methods and calls the handler for the pattern that most closely matches the
	// URL.
	//
	// Pattern matching attempts each pattern in the order in which they were
	// registered.
	router := pat.New()

	// Enable full text indexing to have search feature
	if len(op.Opt.Indexer) > 0 {
		indexer, err := indexers.NewResourceIndexer(op.Opt.Indexer)
		if err != nil {
			log.Errorln(err)
		} else {
			indexer.RegisterRouters(router)
			op.Indexer = indexer
		}
	}

	// Enable pod -> service, service -> serviceMonitor indexing
	if op.Opt.EnableReverseIndex {
		ri, err := indexers.NewReverseIndexer(op.KubeClient, op.Opt.Indexer)
		if err != nil {
			log.Errorln("Failed to create indexer", err)
		} else {
			ri.RegisterRouters(router)
			op.ReverseIndex = ri
		}
	}

	http.Handle("/", router)
	log.Fatalln(http.ListenAndServe(op.Opt.ServerAddress, nil))
}

func (op *Operator) StartCron() {
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
		if config, err := clientcmd.BuildConfigFromFlags(op.Opt.Master, op.Opt.KubeConfig); err == nil {
			t := time.Now().UTC()
			backupDir := filepath.Join(op.Opt.ScratchDir, "snapshot", t.Format(time.RFC3339))

			err := backup.SnapshotCluster(config, backupDir, op.Config.ClusterSnapshot.Sanitize)
			if err != nil {
				log.Errorln(err)
			}

			// upload to cloud
		}

		// run backup
	})
}

func (op *Operator) RunClusterSnapshotter() error {
	if op.Config.ClusterSnapshot == nil {
		return nil
	}

	osmconfigPath := filepath.Join(op.Opt.ScratchDir, "osm", "config.yaml")
	err := storage.WriteOSMConfig(op.KubeClient, op.Config.ClusterSnapshot.Storage, op.Opt.OperatorNamespace, osmconfigPath)
	if err != nil {
		return err
	}

	sh := shell.NewSession()
	sh.SetDir(op.Opt.ScratchDir)
	sh.ShowCMD = true

	container, err := op.Config.ClusterSnapshot.Storage.Container()
	if err != nil {
		return err
	}

	snapshotter := func() error {
		config, err := clientcmd.BuildConfigFromFlags(op.Opt.Master, op.Opt.KubeConfig)
		if err != nil {
			return err
		}

		t := time.Now().UTC()
		snapshotDir := filepath.Join(op.Opt.ScratchDir, "snapshot", t.Format(time.RFC3339))
		err = backup.SnapshotCluster(config, snapshotDir, op.Config.ClusterSnapshot.Sanitize)
		if err != nil {
			return err
		}

		dest, err := op.Config.ClusterSnapshot.Storage.Location(t)
		if err != nil {
			return err
		}

		return sh.Command("osm", "push", "-c", container, snapshotDir, dest).Run()
	}

	_, err = op.Cron.AddFunc(op.Config.ClusterSnapshot.Schedule, func() {
		err := snapshotter()
		if err != nil {
			log.Errorln(err)
			return
		}
	})
	return err
}

func (op *Operator) RunAndHold() {
	op.StartCron()
	op.RunWatchers()
	op.ListenAndServe()
}
