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
	"sync"
	"time"

	"github.com/appscode/go/log"
	api "github.com/appscode/kubed/apis/kubed/v1alpha1"
	"github.com/appscode/kubed/pkg/syncer"

	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
	_ "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	core_informers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	_ "kmodules.xyz/client-go/apiextensions/v1beta1"
	"kmodules.xyz/client-go/tools/fsnotify"
)

type Operator struct {
	Config

	ClientConfig *rest.Config

	recorder     record.EventRecorder
	configSyncer *syncer.ConfigSyncer

	KubeClient kubernetes.Interface

	kubeInformerFactory informers.SharedInformerFactory

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
	op.clusterConfig = *cfg

	err = op.configSyncer.Configure(op.clusterConfig.ClusterName, op.clusterConfig.KubeConfigFile, op.clusterConfig.EnableConfigSyncer)
	if err != nil {
		return err
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
	secretInformer.AddEventHandler(op.configSyncer.SecretHandler())

	nsInformer := op.kubeInformerFactory.Core().V1().Namespaces().Informer()
	nsInformer.AddEventHandler(op.configSyncer.NamespaceHandler())
}

func (op *Operator) RunWatchers(stopCh <-chan struct{}) {
	op.kubeInformerFactory.Start(stopCh)

	res := op.kubeInformerFactory.WaitForCacheSync(stopCh)
	for _, v := range res {
		if !v {
			runtime.HandleError(errors.Errorf("timed out waiting for caches to sync"))
			return
		}
	}
}

func (op *Operator) Run(stopCh <-chan struct{}) {
	op.RunWatchers(stopCh)
	go func() {
		utilruntime.Must(op.watcher.Run(stopCh))
	}()

	<-stopCh
	log.Infoln("Stopping kubed controller")
}
