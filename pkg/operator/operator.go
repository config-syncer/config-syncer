package operator

import (
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/appscode/envconfig"
	"github.com/appscode/kubed/pkg/backup"
	"github.com/appscode/kubed/pkg/config"
	"github.com/appscode/kubed/pkg/elasticsearch"
	"github.com/appscode/kubed/pkg/eventer"
	"github.com/appscode/kubed/pkg/indexers"
	"github.com/appscode/kubed/pkg/influxdb"
	rbin "github.com/appscode/kubed/pkg/recyclebin"
	"github.com/appscode/kubed/pkg/storage"
	"github.com/appscode/kubed/pkg/syncer"
	"github.com/appscode/kubed/pkg/util"
	"github.com/appscode/log"
	"github.com/appscode/pat"
	srch_cs "github.com/appscode/searchlight/client/clientset"
	scs "github.com/appscode/stash/client/clientset"
	vcs "github.com/appscode/voyager/client/clientset"
	shell "github.com/codeskyblue/go-sh"
	pcm "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1alpha1"
	prom "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1alpha1"
	kcs "github.com/k8sdb/apimachinery/client/clientset"
	"gopkg.in/robfig/cron.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Options struct {
	Master     string
	KubeConfig string

	ConfigPath         string
	Address            string
	EnableSearchIndex  bool
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

	Opt    Options
	Config config.ClusterConfig

	SearchIndex    *indexers.ResourceIndexer
	ReverseIndex   *indexers.ReverseIndexer
	TrashCan       *rbin.RecycleBin
	Eventer        *eventer.EventForwarder
	Cron           *cron.Cron
	NotifierLoader envconfig.LoaderFunc
	ConfigSyncer   *syncer.ConfigSyncer

	syncPeriod time.Duration
	sync.Mutex
}

func (op *Operator) Setup() error {
	cfg, err := config.LoadConfig(op.Opt.ConfigPath)
	if err != nil {
		return err
	}
	op.Config = *cfg

	op.NotifierLoader, err = op.getLoader()
	if err != nil {
		return err
	}

	if op.Config.TrashCan != nil {
		if op.Config.TrashCan.Path == "" {
			op.Config.TrashCan.Path = filepath.Join(op.Opt.ScratchDir, "transhcan")
		}
		op.TrashCan = &rbin.RecycleBin{
			Spec:   *op.Config.TrashCan,
			Loader: op.NotifierLoader,
		}
	}

	if op.Config.EventForwarder != nil {
		op.Eventer = &eventer.EventForwarder{
			Spec:   *op.Config.EventForwarder,
			Loader: op.NotifierLoader,
		}
	}

	op.ConfigSyncer = &syncer.ConfigSyncer{KubeClient: op.KubeClient}

	op.Cron = cron.New()
	op.Cron.Start()

	if op.Config.InfluxDB != nil {
		janitor := influx.Janitor{Spec: *op.Config.InfluxDB}
		err = janitor.Cleanup()
		if err != nil {
			return err
		}
	}

	op.syncPeriod = time.Minute * 2
	return nil
}

func (op *Operator) getLoader() (envconfig.LoaderFunc, error) {
	cfg, err := op.KubeClient.CoreV1().
		Secrets(op.Opt.OperatorNamespace).
		Get(op.Config.NotifierSecretName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return func(key string) (value string, found bool) {
		var bytes []byte
		bytes, found = cfg.Data[key]
		value = string(bytes)
		return
	}, nil
}

func (op *Operator) RunWatchers() {
	go op.WatchAlertmanagers()
	go op.WatchClusterAlerts()
	go op.WatchClusterRoleBindingV1alpha1()
	go op.WatchClusterRoleBindingV1beta1()
	go op.WatchClusterRoleV1alpha1()
	go op.WatchClusterRoleV1beta1()
	go op.WatchConfigMaps()
	go op.WatchDaemonSets()
	go op.WatchDeploymentApps()
	go op.WatchDeploymentExtensions()
	go op.WatchDormantDatabases()
	go op.WatchElasticsearches()
	go op.WatchEvents()
	go op.WatchIngresses()
	go op.WatchJobs()
	go op.watchNamespaces()
	go op.WatchNodeAlerts()
	go op.WatchNodes()
	go op.WatchPersistentVolumeClaims()
	go op.WatchPersistentVolumes()
	go op.WatchPodAlerts()
	go op.WatchPostgreses()
	go op.WatchPrometheuss()
	go op.WatchReplicaSets()
	go op.WatchReplicationControllers()
	go op.WatchRestics()
	go op.WatchRoleBindingV1alpha1()
	go op.WatchRoleBindingV1beta1()
	go op.WatchRoleV1alpha1()
	go op.WatchRoleV1beta1()
	go op.WatchSecrets()
	go op.watchService()
	go op.WatchEndpoints()
	go op.WatchServiceMonitors()
	go op.WatchStatefulSets()
	go op.WatchStorageClassV1()
	go op.WatchStorageClassV1beta1()
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
	indexDir := filepath.Join(op.Opt.ScratchDir, "bleve")

	if op.Opt.EnableSearchIndex {
		si, err := indexers.NewResourceIndexer(indexDir)
		if err != nil {
			log.Errorln(err)
		} else {
			si.RegisterRouters(router)
			op.SearchIndex = si
		}
	}

	// Enable pod -> service, service -> serviceMonitor indexing
	if op.Opt.EnableReverseIndex {
		ri, err := indexers.NewReverseIndexer(op.KubeClient, op.PromClient, indexDir)
		if err != nil {
			log.Errorln("Failed to create indexer", err)
		} else {
			router.Get("/api/v1/namespaces/:namespace/:resource/:name/services", http.HandlerFunc(ri.Service.ServeHTTP))
			if util.IsPreferredAPIResource(op.KubeClient, prom.TPRGroup+"/"+prom.TPRVersion, prom.TPRServiceMonitorsKind) {
				// Add Indexer only if Server support this resource
				router.Get("/apis/"+prom.TPRGroup+"/"+prom.TPRVersion+"/namespaces/:namespace/:resource/:name/"+prom.TPRServiceMonitorName, http.HandlerFunc(ri.ServiceMonitor.ServeHTTP))
			}
			if util.IsPreferredAPIResource(op.KubeClient, prom.TPRGroup+"/"+prom.TPRVersion, prom.TPRPrometheusesKind) {
				// Add Indexer only if Server support this resource
				router.Get("/apis/"+prom.TPRGroup+"/"+prom.TPRVersion+"/namespaces/:namespace/:resource/:name/"+prom.TPRPrometheusName, http.HandlerFunc(ri.Prometheus.ServeHTTP))
			}
			op.ReverseIndex = ri
		}
	}

	http.Handle("/", router)
	log.Fatalln(http.ListenAndServe(op.Opt.Address, nil))
}

func (op *Operator) RunElasticsearchCleaner() error {
	if op.Config.Elasticsearch == nil {
		return nil
	}

	janitor := es.Janitor{Spec: *op.Config.Elasticsearch}
	err := janitor.Cleanup()
	if err != nil {
		return err
	}
	op.Cron.AddFunc("@every 6h", func() {
		err := janitor.Cleanup()
		if err != nil {
			log.Errorln(err)
		}
	})
	return nil
}

func (op *Operator) RunTrashCanCleaner() error {
	if op.TrashCan == nil {
		return nil
	}

	_, err := op.Cron.AddFunc("@every 6h", func() {
		err := op.TrashCan.Cleanup()
		if err != nil {
			log.Errorln(err)
		}
	})
	return err
}

func (op *Operator) RunSnapshotter() error {
	if op.Config.ClusterSnapshot == nil {
		return nil
	}

	osmconfigPath := filepath.Join(op.Opt.ScratchDir, "osm", "config.yaml")
	err := storage.WriteOSMConfig(op.KubeClient, op.Config.ClusterSnapshot.Storage, op.Opt.OperatorNamespace, osmconfigPath)
	if err != nil {
		return err
	}

	container, err := op.Config.ClusterSnapshot.Storage.Container()
	if err != nil {
		return err
	}

	snapshotter := func() error {
		cfg, err := clientcmd.BuildConfigFromFlags(op.Opt.Master, op.Opt.KubeConfig)
		if err != nil {
			return err
		}

		t := time.Now().UTC()
		snapshotDir := filepath.Join(op.Opt.ScratchDir, "snapshot", t.Format(config.TimestampFormat))
		err = backup.SnapshotCluster(cfg, snapshotDir, op.Config.ClusterSnapshot.Sanitize)
		if err != nil {
			return err
		}

		dest, err := op.Config.ClusterSnapshot.Storage.Location(t)
		if err != nil {
			return err
		}

		sh := shell.NewSession()
		sh.SetDir(op.Opt.ScratchDir)
		sh.ShowCMD = true
		return sh.Command("osm", "push", "--osmconfig", osmconfigPath, "-c", container, snapshotDir, dest).Run()
	}

	err = snapshotter()
	if err != nil {
		return err
	}

	_, err = op.Cron.AddFunc(op.Config.ClusterSnapshot.Schedule, func() {
		err := snapshotter()
		if err != nil {
			log.Errorln(err)
		}
	})
	return err
}

func (op *Operator) RunAndHold() {
	var err error

	err = op.RunElasticsearchCleaner()
	if err != nil {
		log.Fatalln(err)
	}

	err = op.RunTrashCanCleaner()
	if err != nil {
		log.Fatalln(err)
	}

	err = op.RunSnapshotter()
	if err != nil {
		log.Fatalln(err)
	}

	op.RunWatchers()
	op.ListenAndServe()
}
