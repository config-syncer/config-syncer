package operator

import (
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/appscode/go/log"
	api "github.com/appscode/kubed/apis/kubed/v1alpha1"
	es "github.com/appscode/kubed/pkg/elasticsearch"
	"github.com/appscode/kubed/pkg/eventer"
	influx "github.com/appscode/kubed/pkg/influxdb"
	rbin "github.com/appscode/kubed/pkg/recyclebin"
	indexers "github.com/appscode/kubed/pkg/registry/resource"
	"github.com/appscode/kubed/pkg/storage"
	"github.com/appscode/kubed/pkg/syncer"
	searchlight_api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	srch_cs "github.com/appscode/searchlight/client/clientset/versioned"
	searchlightinformers "github.com/appscode/searchlight/client/informers/externalversions"
	voyager_api "github.com/appscode/voyager/apis/voyager/v1beta1"
	vcs "github.com/appscode/voyager/client/clientset/versioned"
	voyagerinformers "github.com/appscode/voyager/client/informers/externalversions"
	shell "github.com/codeskyblue/go-sh"
	promapi "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	prominformers "github.com/coreos/prometheus-operator/pkg/client/informers/externalversions"
	pcm "github.com/coreos/prometheus-operator/pkg/client/versioned"
	kubedb_api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	kcs "github.com/kubedb/apimachinery/client/clientset/versioned"
	kubedbinformers "github.com/kubedb/apimachinery/client/informers/externalversions"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"gomodules.xyz/envconfig"
	apps "k8s.io/api/apps/v1"
	batch "k8s.io/api/batch/v1"
	certificates "k8s.io/api/certificates/v1beta1"
	core "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	rbac "k8s.io/api/rbac/v1"
	storage_v1 "k8s.io/api/storage/v1"
	_ "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	core_informers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	_ "kmodules.xyz/client-go/apiextensions/v1beta1"
	"kmodules.xyz/client-go/discovery"
	"kmodules.xyz/client-go/tools/backup"
	"kmodules.xyz/client-go/tools/fsnotify"
	"kmodules.xyz/client-go/tools/queue"
	stash_api "stash.appscode.dev/stash/apis/stash/v1alpha1"
	scs "stash.appscode.dev/stash/client/clientset/versioned"
	stashinformers "stash.appscode.dev/stash/client/informers/externalversions"
)

type Operator struct {
	Config

	ClientConfig *rest.Config

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
	PromClient        pcm.Interface

	kubeInformerFactory        informers.SharedInformerFactory
	voyagerInformerFactory     voyagerinformers.SharedInformerFactory
	stashInformerFactory       stashinformers.SharedInformerFactory
	searchlightInformerFactory searchlightinformers.SharedInformerFactory
	kubedbInformerFactory      kubedbinformers.SharedInformerFactory
	promInformerFactory        prominformers.SharedInformerFactory

	Indexer *indexers.ResourceIndexer

	watcher *fsnotify.Watcher

	clusterConfig api.ClusterConfig
	lock          sync.RWMutex
}

func (op *Operator) Configure() error {
	log.Infoln("configuring kubed ...")

	op.lock.Lock()
	defer op.lock.Unlock()

	var err error

	cfg, err := api.LoadConfig(op.ConfigPath)
	if err != nil {
		return err
	}
	err = cfg.Validate()
	if err != nil {
		return err
	}
	op.clusterConfig = *cfg

	if op.clusterConfig.RecycleBin != nil && op.clusterConfig.RecycleBin.Path == "" {
		op.clusterConfig.RecycleBin.Path = filepath.Join(op.ScratchDir, "trashcan")
	}

	op.notifierCred, err = op.getLoader()
	if err != nil {
		return err
	}

	err = op.trashCan.Configure(op.clusterConfig.ClusterName, op.clusterConfig.RecycleBin)
	if err != nil {
		return err
	}

	err = op.eventProcessor.Configure(op.clusterConfig.ClusterName, op.clusterConfig.EventForwarder, op.notifierCred)
	if err != nil {
		return err
	}

	err = op.configSyncer.Configure(op.clusterConfig.ClusterName, op.clusterConfig.KubeConfigFile, op.clusterConfig.EnableConfigSyncer)
	if err != nil {
		return err
	}

	for _, j := range op.clusterConfig.Janitors {
		if j.Kind == api.JanitorInfluxDB {
			janitor := influx.Janitor{Spec: *j.InfluxDB, TTL: j.TTL.Duration}
			err = janitor.Cleanup()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (op *Operator) setupWorkloadInformers() {
	deploymentInformer := op.kubeInformerFactory.Apps().V1().Deployments().Informer()
	op.addEventHandlers(deploymentInformer, apps.SchemeGroupVersion.WithKind("Deployment"))

	rcInformer := op.kubeInformerFactory.Core().V1().ReplicationControllers().Informer()
	op.addEventHandlers(rcInformer, core.SchemeGroupVersion.WithKind("ReplicationController"))

	rsInformer := op.kubeInformerFactory.Apps().V1().ReplicaSets().Informer()
	op.addEventHandlers(rsInformer, apps.SchemeGroupVersion.WithKind("ReplicaSet"))

	daemonSetInformer := op.kubeInformerFactory.Apps().V1().DaemonSets().Informer()
	op.addEventHandlers(daemonSetInformer, apps.SchemeGroupVersion.WithKind("DaemonSet"))

	statefulSetInformer := op.kubeInformerFactory.Apps().V1().StatefulSets().Informer()
	op.addEventHandlers(statefulSetInformer, apps.SchemeGroupVersion.WithKind("StatefulSet"))

	jobInformer := op.kubeInformerFactory.Batch().V1().Jobs().Informer()
	op.addEventHandlers(jobInformer, batch.SchemeGroupVersion.WithKind("Job"))

	op.kubeInformerFactory.Core().V1().Pods().Informer()
}

func (op *Operator) setupNetworkInformers() {
	svcInformer := op.kubeInformerFactory.Core().V1().Services().Informer()
	op.addEventHandlers(svcInformer, core.SchemeGroupVersion.WithKind("Service"))

	ingressInformer := op.kubeInformerFactory.Extensions().V1beta1().Ingresses().Informer()
	op.addEventHandlers(ingressInformer, extensions.SchemeGroupVersion.WithKind("Ingress"))

	netPolicyInformer := op.kubeInformerFactory.Networking().V1().NetworkPolicies().Informer()
	op.addEventHandlers(netPolicyInformer, core.SchemeGroupVersion.WithKind("NetworkPolicy"))
}

func (op *Operator) setupConfigInformers() {
	configMapInformer := op.kubeInformerFactory.InformerFor(&core.ConfigMap{}, func(client kubernetes.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
		return core_informers.NewFilteredConfigMapInformer(
			client,
			op.clusterConfig.ConfigSourceNamespace,
			resyncPeriod,
			cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
			func(options *metav1.ListOptions) {},
		)
	})
	op.addEventHandlers(configMapInformer, core.SchemeGroupVersion.WithKind("ConfigMap"))
	configMapInformer.AddEventHandler(op.configSyncer.ConfigMapHandler())

	secretInformer := op.kubeInformerFactory.InformerFor(&core.Secret{}, func(client kubernetes.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
		return core_informers.NewFilteredSecretInformer(
			client,
			op.clusterConfig.ConfigSourceNamespace,
			resyncPeriod,
			cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
			func(options *metav1.ListOptions) {},
		)
	})
	op.addEventHandlers(secretInformer, core.SchemeGroupVersion.WithKind("Secret"))
	secretInformer.AddEventHandler(op.configSyncer.SecretHandler())

	nsInformer := op.kubeInformerFactory.Core().V1().Namespaces().Informer()
	op.addEventHandlers(nsInformer, core.SchemeGroupVersion.WithKind("Namespace"))
	nsInformer.AddEventHandler(op.configSyncer.NamespaceHandler())
}

func (op *Operator) setupRBACInformers() {
	clusterRoleInformer := op.kubeInformerFactory.Rbac().V1().ClusterRoles().Informer()
	op.addEventHandlers(clusterRoleInformer, rbac.SchemeGroupVersion.WithKind("ClusterRole"))

	clusterRoleBindingInformer := op.kubeInformerFactory.Rbac().V1().ClusterRoleBindings().Informer()
	op.addEventHandlers(clusterRoleBindingInformer, rbac.SchemeGroupVersion.WithKind("ClusterRoleBinding"))

	roleInformer := op.kubeInformerFactory.Rbac().V1().Roles().Informer()
	op.addEventHandlers(roleInformer, rbac.SchemeGroupVersion.WithKind("Role"))

	roleBindingInformer := op.kubeInformerFactory.Rbac().V1().RoleBindings().Informer()
	op.addEventHandlers(roleBindingInformer, rbac.SchemeGroupVersion.WithKind("RoleBinding"))
}

func (op *Operator) setupCoreInformers() {
	nodeInformer := op.kubeInformerFactory.Core().V1().Nodes().Informer()
	op.addEventHandlers(nodeInformer, core.SchemeGroupVersion.WithKind("Node"))

	limitRangeInformer := op.kubeInformerFactory.Core().V1().LimitRanges().Informer()
	op.addEventHandlers(limitRangeInformer, core.SchemeGroupVersion.WithKind("LimitRange"))

	saInformer := op.kubeInformerFactory.Core().V1().ServiceAccounts().Informer()
	op.addEventHandlers(saInformer, core.SchemeGroupVersion.WithKind("ServiceAccount"))
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
	op.addEventHandlers(csrInformer, certificates.SchemeGroupVersion.WithKind("CertificateSigningRequest"))
}

func (op *Operator) setupStorageInformers() {
	pvInformer := op.kubeInformerFactory.Core().V1().PersistentVolumes().Informer()
	op.addEventHandlers(pvInformer, core.SchemeGroupVersion.WithKind("PersistentVolume"))

	pvcInformer := op.kubeInformerFactory.Core().V1().PersistentVolumeClaims().Informer()
	op.addEventHandlers(pvcInformer, core.SchemeGroupVersion.WithKind("PersistentVolumeClaim"))

	storageClassInformer := op.kubeInformerFactory.Storage().V1().StorageClasses().Informer()
	op.addEventHandlers(storageClassInformer, storage_v1.SchemeGroupVersion.WithKind("StorageClass"))
}

func (op *Operator) setupVoyagerInformers() {
	if discovery.IsPreferredAPIResource(op.KubeClient.Discovery(), voyager_api.SchemeGroupVersion.String(), voyager_api.ResourceKindIngress) {
		voyagerIngressInformer := op.voyagerInformerFactory.Voyager().V1beta1().Ingresses().Informer()
		op.addEventHandlers(voyagerIngressInformer, voyager_api.SchemeGroupVersion.WithKind(voyager_api.ResourceKindIngress))

		voyagerCertificateInformer := op.voyagerInformerFactory.Voyager().V1beta1().Certificates().Informer()
		op.addEventHandlers(voyagerCertificateInformer, voyager_api.SchemeGroupVersion.WithKind(voyager_api.ResourceKindCertificate))
	}
}

func (op *Operator) setupStashInformers() {
	if discovery.IsPreferredAPIResource(op.KubeClient.Discovery(), stash_api.SchemeGroupVersion.String(), stash_api.ResourceKindRestic) {
		resticsInformer := op.stashInformerFactory.Stash().V1alpha1().Restics().Informer()
		op.addEventHandlers(resticsInformer, stash_api.SchemeGroupVersion.WithKind(stash_api.ResourceKindRestic))

		recoveryInformer := op.stashInformerFactory.Stash().V1alpha1().Recoveries().Informer()
		op.addEventHandlers(recoveryInformer, stash_api.SchemeGroupVersion.WithKind(stash_api.ResourceKindRecovery))
	}
}

func (op *Operator) setupSearchlightInformers() {
	if discovery.IsPreferredAPIResource(op.KubeClient.Discovery(), searchlight_api.SchemeGroupVersion.String(), searchlight_api.ResourceKindClusterAlert) {
		clusterAlertInformer := op.searchlightInformerFactory.Monitoring().V1alpha1().ClusterAlerts().Informer()
		op.addEventHandlers(clusterAlertInformer, searchlight_api.SchemeGroupVersion.WithKind(searchlight_api.ResourceKindClusterAlert))

		nodeAlertInformer := op.searchlightInformerFactory.Monitoring().V1alpha1().NodeAlerts().Informer()
		op.addEventHandlers(nodeAlertInformer, searchlight_api.SchemeGroupVersion.WithKind(searchlight_api.ResourceKindNodeAlert))

		podAlertInformer := op.searchlightInformerFactory.Monitoring().V1alpha1().PodAlerts().Informer()
		op.addEventHandlers(podAlertInformer, searchlight_api.SchemeGroupVersion.WithKind(searchlight_api.ResourceKindPodAlert))
	}
}

func (op *Operator) setupKubeDBInformers() {
	if discovery.IsPreferredAPIResource(op.KubeClient.Discovery(), kubedb_api.SchemeGroupVersion.String(), kubedb_api.ResourceKindPostgres) {
		pgInformer := op.kubedbInformerFactory.Kubedb().V1alpha1().Postgreses().Informer()
		op.addEventHandlers(pgInformer, kubedb_api.SchemeGroupVersion.WithKind(kubedb_api.ResourceKindPostgres))

		esInformer := op.kubedbInformerFactory.Kubedb().V1alpha1().Elasticsearches().Informer()
		op.addEventHandlers(esInformer, kubedb_api.SchemeGroupVersion.WithKind(kubedb_api.ResourceKindElasticsearch))

		myInformer := op.kubedbInformerFactory.Kubedb().V1alpha1().MySQLs().Informer()
		op.addEventHandlers(myInformer, kubedb_api.SchemeGroupVersion.WithKind(kubedb_api.ResourceKindMySQL))

		mgInformer := op.kubedbInformerFactory.Kubedb().V1alpha1().MongoDBs().Informer()
		op.addEventHandlers(mgInformer, kubedb_api.SchemeGroupVersion.WithKind(kubedb_api.ResourceKindMongoDB))

		rdInformer := op.kubedbInformerFactory.Kubedb().V1alpha1().Redises().Informer()
		op.addEventHandlers(rdInformer, kubedb_api.SchemeGroupVersion.WithKind(kubedb_api.ResourceKindRedis))

		mcInformer := op.kubedbInformerFactory.Kubedb().V1alpha1().Memcacheds().Informer()
		op.addEventHandlers(mcInformer, kubedb_api.SchemeGroupVersion.WithKind(kubedb_api.ResourceKindMemcached))

		dbSnapshotInformer := op.kubedbInformerFactory.Kubedb().V1alpha1().Snapshots().Informer()
		op.addEventHandlers(dbSnapshotInformer, kubedb_api.SchemeGroupVersion.WithKind(kubedb_api.ResourceKindSnapshot))

		dormantDatabaseInformer := op.kubedbInformerFactory.Kubedb().V1alpha1().DormantDatabases().Informer()
		op.addEventHandlers(dormantDatabaseInformer, kubedb_api.SchemeGroupVersion.WithKind(kubedb_api.ResourceKindDormantDatabase))
	}
}

func (op *Operator) setupPrometheusInformers() {
	if discovery.IsPreferredAPIResource(op.KubeClient.Discovery(), promapi.SchemeGroupVersion.String(), promapi.PrometheusesKind) {
		promInf := op.promInformerFactory.Monitoring().V1().Prometheuses().Informer()
		op.addEventHandlers(promInf, promapi.SchemeGroupVersion.WithKind(promapi.PrometheusesKind))

		ruleInf := op.promInformerFactory.Monitoring().V1().PrometheusRules().Informer()
		op.addEventHandlers(ruleInf, promapi.SchemeGroupVersion.WithKind(promapi.PrometheusRuleKind))

		smonInf := op.promInformerFactory.Monitoring().V1().ServiceMonitors().Informer()
		op.addEventHandlers(smonInf, promapi.SchemeGroupVersion.WithKind(promapi.ServiceMonitorsKind))

		amgrInf := op.promInformerFactory.Monitoring().V1().Alertmanagers().Informer()
		op.addEventHandlers(amgrInf, promapi.SchemeGroupVersion.WithKind(promapi.AlertmanagersKind))
	}
}

func (op *Operator) addEventHandlers(informer cache.SharedIndexInformer, gvk schema.GroupVersionKind) {
	informer.AddEventHandler(queue.NewVersionedHandler(op.trashCan, gvk))
	informer.AddEventHandler(queue.NewVersionedHandler(op.eventProcessor, gvk))
	informer.AddEventHandler(queue.NewVersionedHandler(op.Indexer, gvk))
}

func (op *Operator) getLoader() (envconfig.LoaderFunc, error) {
	if op.clusterConfig.NotifierSecretName == "" {
		return func(key string) (string, bool) {
			return "", false
		}, nil
	}
	cfg, err := op.KubeClient.CoreV1().
		Secrets(op.OperatorNamespace).
		Get(op.clusterConfig.NotifierSecretName, metav1.GetOptions{})
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
	op.promInformerFactory.Start(stopCh)

	var res map[reflect.Type]bool

	res = op.kubeInformerFactory.WaitForCacheSync(stopCh)
	for _, v := range res {
		if !v {
			runtime.HandleError(errors.Errorf("timed out waiting for caches to sync"))
			return
		}
	}

	res = op.voyagerInformerFactory.WaitForCacheSync(stopCh)
	for _, v := range res {
		if !v {
			runtime.HandleError(errors.Errorf("timed out waiting for caches to sync"))
			return
		}
	}

	res = op.stashInformerFactory.WaitForCacheSync(stopCh)
	for _, v := range res {
		if !v {
			runtime.HandleError(errors.Errorf("timed out waiting for caches to sync"))
			return
		}
	}

	res = op.searchlightInformerFactory.WaitForCacheSync(stopCh)
	for _, v := range res {
		if !v {
			runtime.HandleError(errors.Errorf("timed out waiting for caches to sync"))
			return
		}
	}

	res = op.kubedbInformerFactory.WaitForCacheSync(stopCh)
	for _, v := range res {
		if !v {
			runtime.HandleError(errors.Errorf("timed out waiting for caches to sync"))
			return
		}
	}

	res = op.promInformerFactory.WaitForCacheSync(stopCh)
	for _, v := range res {
		if !v {
			runtime.HandleError(errors.Errorf("timed out waiting for caches to sync"))
			return
		}
	}
}

func (op *Operator) RunElasticsearchCleaner() error {
	for _, j := range op.clusterConfig.Janitors {
		if j.Kind == api.JanitorElasticsearch {
			var authInfo *api.JanitorAuthInfo

			if j.Elasticsearch.SecretName != "" {
				secret, err := op.KubeClient.CoreV1().Secrets(op.OperatorNamespace).
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

	schedule := "@every 1h"
	if op.Test {
		schedule = "@every 1m"
	}

	_, err := op.cron.AddFunc(schedule, func() {
		err := op.trashCan.Cleanup()
		if err != nil {
			log.Errorln(err)
		}
	})
	return err
}

func (op *Operator) RunSnapshotter() error {
	if op.clusterConfig.Snapshotter == nil {
		return nil
	}

	osmconfigPath := filepath.Join(op.ScratchDir, "osm", "config.yaml")
	err := storage.WriteOSMConfig(op.KubeClient, op.clusterConfig.Snapshotter.Backend, op.OperatorNamespace, osmconfigPath)
	if err != nil {
		return err
	}

	container, err := api.Container(op.clusterConfig.Snapshotter.Backend)
	if err != nil {
		return err
	}

	// test credentials
	sh := shell.NewSession()
	sh.SetDir(op.ScratchDir)
	sh.ShowCMD = true
	snapshotter := func() error {
		mgr := backup.NewBackupManager(op.clusterConfig.ClusterName, op.ClientConfig, op.clusterConfig.Snapshotter.Sanitize)
		snapshotFile, err := mgr.BackupToTar(filepath.Join(op.ScratchDir, "snapshot"))
		if err != nil {
			return err
		}
		defer func() {
			if err := os.Remove(snapshotFile); err != nil {
				log.Errorln(err)
			}
		}()
		dest, err := op.clusterConfig.Snapshotter.Location(filepath.Base(snapshotFile))
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

	if !op.Test { // don't run cronjob for test. it cause problem for consecutive tests.
		_, err := op.cron.AddFunc(op.clusterConfig.Snapshotter.Schedule, func() {
			err := snapshotter()
			if err != nil {
				log.Errorln(err)
			}
		})
		return err
	}
	return nil
}

func (op *Operator) Run(stopCh <-chan struct{}) {
	if err := op.RunElasticsearchCleaner(); err != nil {
		log.Fatalln(err.Error())
	}

	if err := op.RunTrashCanCleaner(); err != nil {
		log.Fatalln(err.Error())
	}

	if err := op.RunSnapshotter(); err != nil {
		log.Fatalln(err.Error())
	}

	op.RunWatchers(stopCh)
	go op.watcher.Run(stopCh)

	<-stopCh
	log.Infoln("Stopping kubed controller")
}
