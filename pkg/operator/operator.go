package operator

import (
	"encoding/json"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/appscode/envconfig"
	"github.com/appscode/go/log"
	v "github.com/appscode/go/version"
	"github.com/appscode/kubed/pkg/api"
	"github.com/appscode/kubed/pkg/config"
	"github.com/appscode/kubed/pkg/elasticsearch"
	"github.com/appscode/kubed/pkg/eventer"
	"github.com/appscode/kubed/pkg/indexers"
	"github.com/appscode/kubed/pkg/influxdb"
	rbin "github.com/appscode/kubed/pkg/recyclebin"
	"github.com/appscode/kubed/pkg/storage"
	"github.com/appscode/kubed/pkg/syncer"
	"github.com/appscode/kubed/pkg/util"
	"github.com/appscode/kutil/tools/backup"
	"github.com/appscode/pat"
	srch_cs "github.com/appscode/searchlight/client/typed/monitoring/v1alpha1"
	scs "github.com/appscode/stash/client/typed/stash/v1alpha1"
	vcs "github.com/appscode/voyager/client/typed/voyager/v1beta1"
	shell "github.com/codeskyblue/go-sh"
	prom "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"
	kcs "github.com/k8sdb/apimachinery/client/typed/kubedb/v1alpha1"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron"
	ecs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Options struct {
	Master     string
	KubeConfig string

	ConfigPath string
	APIAddress string
	WebAddress string

	EnableConfigSync  bool
	ScratchDir        string
	OperatorNamespace string

	ResyncPeriod time.Duration
}

type Operator struct {
	KubeClient        kubernetes.Interface
	VoyagerClient     vcs.VoyagerV1beta1Interface
	SearchlightClient srch_cs.MonitoringV1alpha1Interface
	StashClient       scs.StashV1alpha1Interface
	PromClient        prom.MonitoringV1Interface
	KubeDBClient      kcs.KubedbV1alpha1Interface
	CRDClient         ecs.ApiextensionsV1beta1Interface

	Opt    Options
	Config config.ClusterConfig

	SearchIndex    *indexers.ResourceIndexer
	ReverseIndex   *indexers.ReverseIndexer
	TrashCan       *rbin.RecycleBin
	Eventer        *eventer.EventForwarder
	Cron           *cron.Cron
	NotifierLoader envconfig.LoaderFunc
	ConfigSyncer   *syncer.ConfigSyncer

	sync.Mutex
}

func (op *Operator) Setup() error {
	cfg, err := config.LoadConfig(op.Opt.ConfigPath)
	if err != nil {
		return err
	}
	if op.Opt.APIAddress != "" {
		cfg.APIServer.Address = op.Opt.APIAddress
	}
	err = cfg.Validate()
	if err != nil {
		return err
	}
	op.Config = *cfg

	op.NotifierLoader, err = op.getLoader()
	if err != nil {
		return err
	}

	if op.Config.RecycleBin != nil {
		if op.Config.RecycleBin.Path == "" {
			op.Config.RecycleBin.Path = filepath.Join(op.Opt.ScratchDir, "transhcan")
		}
		op.TrashCan = &rbin.RecycleBin{
			ClusterName: op.Config.ClusterName,
			Spec:        *op.Config.RecycleBin,
			Loader:      op.NotifierLoader,
		}
	}

	if op.Config.EventForwarder != nil {
		op.Eventer = &eventer.EventForwarder{
			ClusterName: op.Config.ClusterName,
			Receivers:   op.Config.EventForwarder.Receivers,
			Loader:      op.NotifierLoader,
		}
	}

	if op.Config.EnableConfigSyncer {
		op.ConfigSyncer = &syncer.ConfigSyncer{KubeClient: op.KubeClient}
	}

	op.Cron = cron.New()
	op.Cron.Start()

	for _, j := range cfg.Janitors {
		if j.Kind == config.JanitorInfluxDB {
			janitor := influx.Janitor{Spec: *j.InfluxDB, TTL: j.TTL.Duration}
			err = janitor.Cleanup()
			if err != nil {
				return err
			}
		}
	}

	// Enable full text indexing to have search feature
	indexDir := filepath.Join(op.Opt.ScratchDir, "bleve")
	if op.Config.APIServer.EnableSearchIndex {
		si, err := indexers.NewResourceIndexer(indexDir)
		if err != nil {
			return err
		}
		op.SearchIndex = si
	}
	// Enable pod -> service, service -> serviceMonitor indexing
	if op.Config.APIServer.EnableReverseIndex {
		ri, err := indexers.NewReverseIndexer(op.KubeClient, op.PromClient, indexDir)
		if err != nil {
			return err
		}
		op.ReverseIndex = ri
	}

	op.Opt.ResyncPeriod = time.Minute * 2
	return nil
}

func (op *Operator) getLoader() (envconfig.LoaderFunc, error) {
	if op.Config.NotifierSecretName == "" {
		return func(key string) (string, bool) {
			return "", false
		}, nil
	}
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
	go op.WatchAlertmanagerV1()
	go op.WatchCertificateSigningRequests()
	go op.WatchClusterAlerts()
	go op.WatchClusterRoleBindingV1beta1()
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
	go op.WatchPrometheusV1()
	go op.WatchReplicaSets()
	go op.WatchReplicationControllers()
	go op.WatchRestics()
	go op.WatchRoleBindingV1beta1()
	go op.WatchRoleV1alpha1()
	go op.WatchRoleV1beta1()
	go op.WatchSecrets()
	go op.watchService()
	go op.WatchEndpoints()
	go op.WatchServiceMonitorV1()
	go op.WatchStatefulSets()
	go op.WatchStorageClassV1()
	go op.WatchVoyagerCertificates()
	go op.WatchVoyagerIngresses()
}

func (op *Operator) RunAPIServer() {
	// router is default HTTP request multiplexer for kubed. It matches the URL of each
	// incoming request against a list of registered patterns with their associated
	// methods and calls the handler for the pattern that most closely matches the
	// URL.
	//
	// Pattern matching attempts each pattern in the order in which they were
	// registered.
	router := pat.New()
	if op.Config.APIServer.EnableSearchIndex {
		op.SearchIndex.RegisterRouters(router)
	}
	// Enable pod -> service, service -> serviceMonitor indexing
	if op.Config.APIServer.EnableReverseIndex {
		router.Get("/api/v1/namespaces/:namespace/:resource/:name/services", http.HandlerFunc(op.ReverseIndex.Service.ServeHTTP))
		if util.IsPreferredAPIResource(op.KubeClient, prom.Group+"/"+prom.Version, prom.ServiceMonitorsKind) {
			// Add Indexer only if Server support this resource
			router.Get("/apis/"+prom.Group+"/"+prom.Version+"/namespaces/:namespace/:resource/:name/"+prom.ServiceMonitorName, http.HandlerFunc(op.ReverseIndex.ServiceMonitor.ServeHTTP))
		}
		if util.IsPreferredAPIResource(op.KubeClient, prom.Group+"/"+prom.Version, prom.PrometheusesKind) {
			// Add Indexer only if Server support this resource
			router.Get("/apis/"+prom.Group+"/"+prom.Version+"/namespaces/:namespace/:resource/:name/"+prom.PrometheusName, http.HandlerFunc(op.ReverseIndex.Prometheus.ServeHTTP))
		}
	}

	router.Get("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))
	router.Get("/metadata", http.HandlerFunc(op.metadataHandler))
	log.Fatalln(http.ListenAndServe(op.Config.APIServer.Address, router))
}

func (op *Operator) metadataHandler(w http.ResponseWriter, r *http.Request) {
	resp := &api.KubedMetadata{
		OperatorNamespace:   op.Opt.OperatorNamespace,
		SearchEnabled:       op.Config.APIServer.EnableSearchIndex,
		ReverseIndexEnabled: op.Config.APIServer.EnableReverseIndex,
		Version:             &v.Version,
	}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("x-content-type-options", "nosniff")
}

func (op *Operator) RunElasticsearchCleaner() error {
	for _, j := range op.Config.Janitors {
		if j.Kind == config.JanitorElasticsearch {
			var authInfo *config.JanitorAuthInfo

			if j.Elasticsearch.SecretName != "" {
				secret, err := op.KubeClient.CoreV1().Secrets(op.Opt.OperatorNamespace).
					Get(j.Elasticsearch.SecretName, metav1.GetOptions{})
				if err != nil && !kerr.IsNotFound(err) {
					return err
				}
				if secret != nil {
					authInfo = config.LoadJanitorAuthInfo(secret.Data)
				}
			}

			janitor := es.Janitor{Spec: *j.Elasticsearch, AuthInfo: authInfo, TTL: j.TTL.Duration}
			err := janitor.Cleanup()
			if err != nil {
				return err
			}
			op.Cron.AddFunc("@every 1h", func() {
				err := janitor.Cleanup()
				if err != nil {
					log.Errorln(err)
				}
			})
		}
	}
	return nil
}

func (op *Operator) RunTrashCanCleaner() error {
	if op.TrashCan == nil {
		return nil
	}

	return op.Cron.AddFunc("@every 1h", func() {
		err := op.TrashCan.Cleanup()
		if err != nil {
			log.Errorln(err)
		}
	})
}

func (op *Operator) RunSnapshotter() error {
	if op.Config.Snapshotter == nil {
		return nil
	}

	osmconfigPath := filepath.Join(op.Opt.ScratchDir, "osm", "config.yaml")
	err := storage.WriteOSMConfig(op.KubeClient, op.Config.Snapshotter.Backend, op.Opt.OperatorNamespace, osmconfigPath)
	if err != nil {
		return err
	}

	container, err := op.Config.Snapshotter.Backend.Container()
	if err != nil {
		return err
	}

	// test credentials
	sh := shell.NewSession()
	sh.SetDir(op.Opt.ScratchDir)
	sh.ShowCMD = true
	snapshotter := func() error {
		restConfig, err := clientcmd.BuildConfigFromFlags(op.Opt.Master, op.Opt.KubeConfig)
		if err != nil {
			return err
		}

		mgr := backup.NewBackupManager(op.Config.ClusterName, restConfig, op.Config.Snapshotter.Sanitize)
		snapshotFile, err := mgr.BackupToTar(filepath.Join(op.Opt.ScratchDir, "snapshot"))
		if err != nil {
			return err
		}
		defer func() {
			if err := os.Remove(snapshotFile); err != nil {
				log.Errorln(err)
			}
		}()
		dest, err := op.Config.Snapshotter.Location(filepath.Base(snapshotFile))
		if err != nil {
			return err
		}
		return sh.Command("osm", "push", "--osmconfig", osmconfigPath, "-c", container, snapshotFile, dest).Run()
	}
	// start taking first backup
	go func() {
		err := snapshotter()
		if err != nil {
			log.Errorln(err)
		}
	}()
	return op.Cron.AddFunc(op.Config.Snapshotter.Schedule, func() {
		err := snapshotter()
		if err != nil {
			log.Errorln(err)
		}
	})
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
	go op.RunAPIServer()

	m := pat.New()
	m.Get("/metrics", promhttp.Handler())
	http.Handle("/", m)
	log.Infoln("Listening on", op.Opt.WebAddress)
	log.Fatal(http.ListenAndServe(op.Opt.WebAddress, nil))
}
