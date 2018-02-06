package operator

import (
	"encoding/json"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/appscode/envconfig"
	"github.com/appscode/go/log"
	prom_util "github.com/appscode/kube-mon/prometheus/v1"
	"github.com/appscode/kubed/pkg/api"
	"github.com/appscode/kubed/pkg/elasticsearch"
	"github.com/appscode/kubed/pkg/eventer"
	"github.com/appscode/kubed/pkg/indexers"
	"github.com/appscode/kubed/pkg/influxdb"
	rbin "github.com/appscode/kubed/pkg/recyclebin"
	"github.com/appscode/kubed/pkg/storage"
	"github.com/appscode/kubed/pkg/syncer"
	"github.com/appscode/kutil/discovery"
	"github.com/appscode/kutil/tools/backup"
	"github.com/appscode/kutil/tools/fsnotify"
	"github.com/appscode/pat"
	searchlight_api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	srch_cs "github.com/appscode/searchlight/client"
	searchlightinformers "github.com/appscode/searchlight/informers/externalversions"
	stash_api "github.com/appscode/stash/apis/stash/v1alpha1"
	scs "github.com/appscode/stash/client"
	stashinformers "github.com/appscode/stash/informers/externalversions"
	voyager_api "github.com/appscode/voyager/apis/voyager/v1beta1"
	vcs "github.com/appscode/voyager/client"
	voyagerinformers "github.com/appscode/voyager/informers/externalversions"
	shell "github.com/codeskyblue/go-sh"
	prom "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"
	kubedb_api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	kcs "github.com/kubedb/apimachinery/client"
	kubedbinformers "github.com/kubedb/apimachinery/informers/externalversions"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	core_informers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"
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
	QPS          float32
	Burst        int

	PrometheusCrdGroup string
	PrometheusCrdKinds prom.CrdKinds
}

type Operator struct {
	options Options

	notifierCred   envconfig.LoaderFunc
	recorder       record.EventRecorder
	trashCan       *rbin.RecycleBin
	eventProcessor *eventer.EventForwarder
	configSyncer   *syncer.ConfigSyncer

	cron *cron.Cron

	KubeClient        kubernetes.Interface
	VoyagerClient     vcs.Interface
	SearchlightClient srch_cs.Interface
	StashClient       scs.Interface
	KubeDBClient      kcs.Interface
	PromClient        prom.MonitoringV1Interface

	kubeInformerFactory        informers.SharedInformerFactory
	voyagerInformerFactory     voyagerinformers.SharedInformerFactory
	stashInformerFactory       stashinformers.SharedInformerFactory
	searchlightInformerFactory searchlightinformers.SharedInformerFactory
	kubedbInformerFactory      kubedbinformers.SharedInformerFactory
	promInf                    cache.SharedIndexInformer
	smonInf                    cache.SharedIndexInformer
	amgrInf                    cache.SharedIndexInformer

	searchIndexer *indexers.ResourceIndexer

	watcher *fsnotify.Watcher

	config api.ClusterConfig
	lock   sync.RWMutex
}

func New(config *rest.Config, opt Options) (*Operator, error) {
	var err error
	op := &Operator{options: opt}

	if op.KubeClient, err = kubernetes.NewForConfig(config); err != nil {
		return nil, err
	}
	if op.VoyagerClient, err = vcs.NewForConfig(config); err != nil {
		return nil, err
	}
	if op.SearchlightClient, err = srch_cs.NewForConfig(config); err != nil {
		return nil, err
	}
	if op.StashClient, err = scs.NewForConfig(config); err != nil {
		return nil, err
	}
	if op.KubeDBClient, err = kcs.NewForConfig(config); err != nil {
		return nil, err
	}

	op.PromClient, err = prom.NewForConfig(&opt.PrometheusCrdKinds, opt.PrometheusCrdGroup, config)
	if err != nil {
		return nil, err
	}

	op.recorder = eventer.NewEventRecorder(op.KubeClient, "kubed")
	op.trashCan = &rbin.RecycleBin{}
	op.eventProcessor = &eventer.EventForwarder{Client: op.KubeClient.Discovery()}
	op.configSyncer = syncer.New(op.KubeClient, op.recorder)

	op.cron = cron.New()
	op.cron.Start()

	// Enable full text indexing to have search feature
	indexDir := filepath.Join(op.options.ScratchDir, "indices")
	op.searchIndexer, err = indexers.NewResourceIndexer(indexDir)
	if err != nil {
		return nil, err
	}

	op.watcher = &fsnotify.Watcher{
		WatchDir: filepath.Dir(opt.ConfigPath),
		Reload:   op.Configure,
	}

	// ---------------------------
	op.kubeInformerFactory = informers.NewSharedInformerFactory(op.KubeClient, op.options.ResyncPeriod)
	op.voyagerInformerFactory = voyagerinformers.NewSharedInformerFactory(op.VoyagerClient, op.options.ResyncPeriod)
	op.stashInformerFactory = stashinformers.NewSharedInformerFactory(op.StashClient, op.options.ResyncPeriod)
	op.searchlightInformerFactory = searchlightinformers.NewSharedInformerFactory(op.SearchlightClient, op.options.ResyncPeriod)
	op.kubedbInformerFactory = kubedbinformers.NewSharedInformerFactory(op.KubeDBClient, op.options.ResyncPeriod)
	// ---------------------------
	op.setupWorkloadInformers()
	op.setupNetworkInformers()
	op.setupConfigInformers()
	op.setupRBACInformers()
	op.setupNodeInformers()
	op.setupEventInformers()
	op.setupCertificateInformers()
	// ---------------------------
	op.setupVoyagerInformers()
	op.setupStashInformers()
	op.setupSearchlightInformers()
	op.setupKubeDBInformers()
	op.setupPrometheusInformers()
	// ---------------------------

	return op, nil
}

func (op *Operator) Configure() error {
	op.lock.Lock()
	defer op.lock.Unlock()

	var err error

	cfg, err := api.LoadConfig(op.options.ConfigPath)
	if err != nil {
		return err
	}
	if op.options.APIAddress != "" {
		cfg.APIServer.Address = op.options.APIAddress
	}
	err = cfg.Validate()
	if err != nil {
		return err
	}
	op.config = *cfg

	if op.config.RecycleBin != nil && op.config.RecycleBin.Path == "" {
		op.config.RecycleBin.Path = filepath.Join(op.options.ScratchDir, "transhcan")
	}

	op.notifierCred, err = op.getLoader()
	if err != nil {
		return err
	}

	err = op.trashCan.Configure(op.config.ClusterName, op.config.RecycleBin, op.notifierCred)
	if err != nil {
		return err
	}

	err = op.eventProcessor.Configure(op.config.ClusterName, op.config.EventForwarder, op.notifierCred)
	if err != nil {
		return err
	}

	err = op.configSyncer.Configure(op.config.ClusterName, op.config.KubeConfigFile, op.config.EnableConfigSyncer)
	if err != nil {
		return err
	}

	for _, j := range op.config.Janitors {
		if j.Kind == api.JanitorInfluxDB {
			janitor := influx.Janitor{Spec: *j.InfluxDB, TTL: j.TTL.Duration}
			err = janitor.Cleanup()
			if err != nil {
				return err
			}
		}
	}

	return op.searchIndexer.Configure(op.config.APIServer.EnableSearchIndex)
}

func (op *Operator) setupWorkloadInformers() {
	deploymentInformer := op.kubeInformerFactory.Apps().V1beta1().Deployments().Informer()
	op.addEventHandlers(deploymentInformer)

	rcInformer := op.kubeInformerFactory.Core().V1().ReplicationControllers().Informer()
	op.addEventHandlers(rcInformer)

	rsInformer := op.kubeInformerFactory.Extensions().V1beta1().ReplicaSets().Informer()
	op.addEventHandlers(rsInformer)

	daemonSetInformer := op.kubeInformerFactory.Extensions().V1beta1().DaemonSets().Informer()
	op.addEventHandlers(daemonSetInformer)

	jobInformer := op.kubeInformerFactory.Batch().V1().Jobs().Informer()
	op.addEventHandlers(jobInformer)

	op.kubeInformerFactory.Core().V1().Pods().Informer()
}

func (op *Operator) setupNetworkInformers() {
	svcInformer := op.kubeInformerFactory.Core().V1().Services().Informer()
	op.addEventHandlers(svcInformer)

	ingressInformer := op.kubeInformerFactory.Extensions().V1beta1().Ingresses().Informer()
	op.addEventHandlers(ingressInformer)
}

func (op *Operator) setupConfigInformers() {
	configMapInformer := op.kubeInformerFactory.Core().V1().ConfigMaps().Informer()
	op.addEventHandlers(configMapInformer)
	configMapInformer.AddEventHandler(op.configSyncer.ConfigMapHandler())

	secretInformer := op.kubeInformerFactory.Core().V1().Secrets().Informer()
	op.addEventHandlers(secretInformer)
	secretInformer.AddEventHandler(op.configSyncer.SecretHandler())

	nsInformer := op.kubeInformerFactory.Core().V1().Namespaces().Informer()
	nsInformer.AddEventHandler(op.configSyncer.NamespaceHandler())
}

func (op *Operator) setupRBACInformers() {
	clusterRoleInformer := op.kubeInformerFactory.Rbac().V1beta1().ClusterRoles().Informer()
	op.addEventHandlers(clusterRoleInformer)

	clusterRoleBindingInformer := op.kubeInformerFactory.Rbac().V1beta1().ClusterRoleBindings().Informer()
	op.addEventHandlers(clusterRoleBindingInformer)

	roleInformer := op.kubeInformerFactory.Rbac().V1beta1().Roles().Informer()
	op.addEventHandlers(roleInformer)

	roleBindingInformer := op.kubeInformerFactory.Rbac().V1beta1().RoleBindings().Informer()
	op.addEventHandlers(roleBindingInformer)
}

func (op *Operator) setupNodeInformers() {
	nodeInformer := op.kubeInformerFactory.Core().V1().Nodes().Informer()
	op.addEventHandlers(nodeInformer)
}

func (op *Operator) setupEventInformers() {
	eventInformer := op.kubeInformerFactory.InformerFor(&core.Event{}, func(client kubernetes.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
		return core_informers.NewFilteredEventInformer(
			client,
			core.NamespaceAll,
			resyncPeriod,
			cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
			func(options *metav1.ListOptions) {
				options.FieldSelector = fields.OneTermEqualSelector("type", core.EventTypeWarning).String()
			},
		)
	})
	eventInformer.AddEventHandler(op.eventProcessor)
}

func (op *Operator) setupCertificateInformers() {
	csrInformer := op.kubeInformerFactory.Certificates().V1beta1().CertificateSigningRequests().Informer()
	op.addEventHandlers(csrInformer)
}

func (op *Operator) setupStorageInformers() {
	pvInformer := op.kubeInformerFactory.Core().V1().PersistentVolumes().Informer()
	op.addEventHandlers(pvInformer)

	pvcInformer := op.kubeInformerFactory.Core().V1().PersistentVolumeClaims().Informer()
	op.addEventHandlers(pvcInformer)

	storageClassInformer := op.kubeInformerFactory.Storage().V1().StorageClasses().Informer()
	op.addEventHandlers(storageClassInformer)
}

func (op *Operator) setupVoyagerInformers() {
	if discovery.IsPreferredAPIResource(op.KubeClient.Discovery(), voyager_api.SchemeGroupVersion.String(), voyager_api.ResourceKindIngress) {
		voyagerIngressInformer := op.voyagerInformerFactory.Voyager().V1beta1().Ingresses().Informer()
		op.addEventHandlers(voyagerIngressInformer)

		voyagerCertificateInformer := op.voyagerInformerFactory.Voyager().V1beta1().Certificates().Informer()
		op.addEventHandlers(voyagerCertificateInformer)
	}
}

func (op *Operator) setupStashInformers() {
	if discovery.IsPreferredAPIResource(op.KubeClient.Discovery(), stash_api.SchemeGroupVersion.String(), stash_api.ResourceKindRestic) {
		resticsInformer := op.stashInformerFactory.Stash().V1alpha1().Restics().Informer()
		op.addEventHandlers(resticsInformer)

		recoveryInformer := op.stashInformerFactory.Stash().V1alpha1().Recoveries().Informer()
		op.addEventHandlers(recoveryInformer)
	}
}

func (op *Operator) setupSearchlightInformers() {
	if discovery.IsPreferredAPIResource(op.KubeClient.Discovery(), searchlight_api.SchemeGroupVersion.String(), searchlight_api.ResourceKindClusterAlert) {
		clusterAlertInformer := op.searchlightInformerFactory.Monitoring().V1alpha1().ClusterAlerts().Informer()
		op.addEventHandlers(clusterAlertInformer)

		nodeAlertInformer := op.searchlightInformerFactory.Monitoring().V1alpha1().NodeAlerts().Informer()
		op.addEventHandlers(nodeAlertInformer)

		podAlertInformer := op.searchlightInformerFactory.Monitoring().V1alpha1().PodAlerts().Informer()
		op.addEventHandlers(podAlertInformer)
	}
}

func (op *Operator) setupKubeDBInformers() {
	if discovery.IsPreferredAPIResource(op.KubeClient.Discovery(), kubedb_api.SchemeGroupVersion.String(), kubedb_api.ResourceKindPostgres) {
		pgInformer := op.kubedbInformerFactory.Kubedb().V1alpha1().Postgreses().Informer()
		op.addEventHandlers(pgInformer)

		esInformer := op.kubedbInformerFactory.Kubedb().V1alpha1().Postgreses().Informer()
		op.addEventHandlers(esInformer)

		myInformer := op.kubedbInformerFactory.Kubedb().V1alpha1().MySQLs().Informer()
		op.addEventHandlers(myInformer)

		mgInformer := op.kubedbInformerFactory.Kubedb().V1alpha1().MongoDBs().Informer()
		op.addEventHandlers(mgInformer)

		rdInformer := op.kubedbInformerFactory.Kubedb().V1alpha1().Redises().Informer()
		op.addEventHandlers(rdInformer)

		mcInformer := op.kubedbInformerFactory.Kubedb().V1alpha1().Memcacheds().Informer()
		op.addEventHandlers(mcInformer)

		dbSnapshotInformer := op.kubedbInformerFactory.Kubedb().V1alpha1().Snapshots().Informer()
		op.addEventHandlers(dbSnapshotInformer)

		dormantDatabaseInformer := op.kubedbInformerFactory.Kubedb().V1alpha1().DormantDatabases().Informer()
		op.addEventHandlers(dormantDatabaseInformer)
	}
}

func (op *Operator) setupPrometheusInformers() {
	if discovery.IsPreferredAPIResource(op.KubeClient.Discovery(), prom_util.SchemeGroupVersion.String(), prom.PrometheusesKind) {
		op.promInf = cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc:  op.PromClient.Prometheuses(core.NamespaceAll).List,
				WatchFunc: op.PromClient.Prometheuses(core.NamespaceAll).Watch,
			},
			&prom.Prometheus{}, op.options.ResyncPeriod, cache.Indexers{},
		)
		op.addEventHandlers(op.promInf)

		op.smonInf = cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc:  op.PromClient.ServiceMonitors(core.NamespaceAll).List,
				WatchFunc: op.PromClient.ServiceMonitors(core.NamespaceAll).Watch,
			},
			&prom.ServiceMonitor{}, op.options.ResyncPeriod, cache.Indexers{},
		)
		op.addEventHandlers(op.smonInf)

		op.amgrInf = cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc:  op.PromClient.Alertmanagers(core.NamespaceAll).List,
				WatchFunc: op.PromClient.Alertmanagers(core.NamespaceAll).Watch,
			},
			&prom.Alertmanager{}, op.options.ResyncPeriod, cache.Indexers{},
		)
		op.addEventHandlers(op.amgrInf)
	}
}

func (op *Operator) addEventHandlers(informer cache.SharedIndexInformer) {
	informer.AddEventHandler(op.trashCan)
	informer.AddEventHandler(op.eventProcessor)
	informer.AddEventHandler(op.searchIndexer)
}

func (op *Operator) getLoader() (envconfig.LoaderFunc, error) {
	if op.config.NotifierSecretName == "" {
		return func(key string) (string, bool) {
			return "", false
		}, nil
	}
	cfg, err := op.KubeClient.CoreV1().
		Secrets(op.options.OperatorNamespace).
		Get(op.config.NotifierSecretName, metav1.GetOptions{})
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

func (op *Operator) RunWatchers(stopCh <-chan struct{}) {
	op.kubeInformerFactory.Start(stopCh)
	op.voyagerInformerFactory.Start(stopCh)
	op.stashInformerFactory.Start(stopCh)
	op.searchlightInformerFactory.Start(stopCh)
	op.kubedbInformerFactory.Start(stopCh)
	if op.promInf != nil {
		go op.promInf.Run(stopCh)
		go op.smonInf.Run(stopCh)
		go op.amgrInf.Run(stopCh)
	}

	var res map[reflect.Type]bool

	res = op.kubeInformerFactory.WaitForCacheSync(stopCh)
	for _, v := range res {
		if !v {
			runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
			return
		}
	}

	res = op.voyagerInformerFactory.WaitForCacheSync(stopCh)
	for _, v := range res {
		if !v {
			runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
			return
		}
	}

	res = op.stashInformerFactory.WaitForCacheSync(stopCh)
	for _, v := range res {
		if !v {
			runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
			return
		}
	}

	res = op.searchlightInformerFactory.WaitForCacheSync(stopCh)
	for _, v := range res {
		if !v {
			runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
			return
		}
	}

	res = op.kubedbInformerFactory.WaitForCacheSync(stopCh)
	for _, v := range res {
		if !v {
			runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
			return
		}
	}

	if op.promInf != nil {
		if !cache.WaitForCacheSync(stopCh, op.promInf.HasSynced) {
			runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
			return
		}
		if !cache.WaitForCacheSync(stopCh, op.smonInf.HasSynced) {
			runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
			return
		}
		if !cache.WaitForCacheSync(stopCh, op.amgrInf.HasSynced) {
			runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
			return
		}
	}
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
	if op.config.APIServer.EnableSearchIndex {
		op.searchIndexer.RegisterRouters(router)
	}

	router.Get("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))
	router.Get("/metadata", http.HandlerFunc(op.metadataHandler))
	log.Fatalln(http.ListenAndServe(op.config.APIServer.Address, router))
}

func (op *Operator) metadataHandler(w http.ResponseWriter, r *http.Request) {
	resp := &api.KubedMetadata{
		OperatorNamespace: op.options.OperatorNamespace,
		SearchEnabled:     op.config.APIServer.EnableSearchIndex,
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
	for _, j := range op.config.Janitors {
		if j.Kind == api.JanitorElasticsearch {
			var authInfo *api.JanitorAuthInfo

			if j.Elasticsearch.SecretName != "" {
				secret, err := op.KubeClient.CoreV1().Secrets(op.options.OperatorNamespace).
					Get(j.Elasticsearch.SecretName, metav1.GetOptions{})
				if err != nil && !kerr.IsNotFound(err) {
					return err
				}
				if secret != nil {
					authInfo = api.LoadJanitorAuthInfo(secret.Data)
				}
			}

			janitor := es.Janitor{Spec: *j.Elasticsearch, AuthInfo: authInfo, TTL: j.TTL.Duration}
			err := janitor.Cleanup()
			if err != nil {
				return err
			}
			op.cron.AddFunc("@every 1h", func() {
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
	if op.trashCan == nil {
		return nil
	}

	return op.cron.AddFunc("@every 1h", func() {
		err := op.trashCan.Cleanup()
		if err != nil {
			log.Errorln(err)
		}
	})
}

func (op *Operator) RunSnapshotter() error {
	if op.config.Snapshotter == nil {
		return nil
	}

	osmconfigPath := filepath.Join(op.options.ScratchDir, "osm", "config.yaml")
	err := storage.WriteOSMConfig(op.KubeClient, op.config.Snapshotter.Backend, op.options.OperatorNamespace, osmconfigPath)
	if err != nil {
		return err
	}

	container, err := op.config.Snapshotter.Backend.Container()
	if err != nil {
		return err
	}

	// test credentials
	sh := shell.NewSession()
	sh.SetDir(op.options.ScratchDir)
	sh.ShowCMD = true
	snapshotter := func() error {
		restConfig, err := clientcmd.BuildConfigFromFlags(op.options.Master, op.options.KubeConfig)
		if err != nil {
			return err
		}

		mgr := backup.NewBackupManager(op.config.ClusterName, restConfig, op.config.Snapshotter.Sanitize)
		snapshotFile, err := mgr.BackupToTar(filepath.Join(op.options.ScratchDir, "snapshot"))
		if err != nil {
			return err
		}
		defer func() {
			if err := os.Remove(snapshotFile); err != nil {
				log.Errorln(err)
			}
		}()
		dest, err := op.config.Snapshotter.Location(filepath.Base(snapshotFile))
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
	return op.cron.AddFunc(op.config.Snapshotter.Schedule, func() {
		err := snapshotter()
		if err != nil {
			log.Errorln(err)
		}
	})
}

func (op *Operator) RunAndHold(stopCh <-chan struct{}) {
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

	op.RunWatchers(stopCh)
	go op.RunAPIServer()

	go op.watcher.Run(stopCh)

	m := pat.New()
	m.Get("/metrics", promhttp.Handler())
	http.Handle("/", m)
	log.Infoln("Listening on", op.options.WebAddress)
	log.Fatal(http.ListenAndServe(op.options.WebAddress, nil))
}
