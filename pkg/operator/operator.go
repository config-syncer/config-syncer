/*
Copyright The Kubed Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package operator

import (
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/appscode/go/log"
	api "github.com/appscode/kubed/apis/kubed/v1alpha1"
	"github.com/appscode/kubed/pkg/syncer"

	shell "github.com/codeskyblue/go-sh"
	"github.com/pkg/errors"
	"gomodules.xyz/envconfig"
	core "k8s.io/api/core/v1"
	_ "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	core_informers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	_ "kmodules.xyz/client-go/apiextensions/v1beta1"
	"kmodules.xyz/client-go/tools/backup"
	"kmodules.xyz/client-go/tools/fsnotify"
	storage "kmodules.xyz/objectstore-api/osm"
)

type Operator struct {
	Config

	ClientConfig *rest.Config

	notifierCred envconfig.LoaderFunc
	recorder     record.EventRecorder
	configSyncer *syncer.ConfigSyncer

	KubeClient kubernetes.Interface

	kubeInformerFactory informers.SharedInformerFactory

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
	err := storage.WriteOSMConfig(op.KubeClient, op.OperatorNamespace, op.clusterConfig.Snapshotter.Backend, osmconfigPath)
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
