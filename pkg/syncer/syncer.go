package syncer

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/appscode/go/types"
	api "github.com/appscode/kubed/apis/kubed/v1alpha1"
	"github.com/appscode/kutil/meta"
	clientcmd_util "github.com/appscode/kutil/tools/clientcmd"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"
)

type ConfigSyncer struct {
	kubeClient kubernetes.Interface
	recorder   record.EventRecorder

	clusterName string
	contexts    map[string]clusterContext
	enable      bool
	lock        sync.RWMutex
}

func New(KubeClient kubernetes.Interface, Recorder record.EventRecorder) *ConfigSyncer {
	return &ConfigSyncer{
		kubeClient: KubeClient,
		recorder:   Recorder,
	}
}

func (s *ConfigSyncer) Configure(clusterName string, kubeconfigFile string, enable bool) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.clusterName = clusterName
	s.contexts = map[string]clusterContext{}
	s.enable = enable

	// Parse external kubeconfig file, assume that it doesn't include source cluster
	if kubeconfigFile != "" {
		kConfig, err := clientcmd.LoadFromFile(kubeconfigFile)
		if err != nil {
			return fmt.Errorf("failed to parse context list. Reason: %v", err)
		}

		for contextName := range kConfig.Contexts {
			ctx := clusterContext{}

			cfg, err := clientcmd_util.BuildConfigFromContext(kubeconfigFile, contextName)
			if err != nil {
				continue
			}
			if ctx.Client, err = kubernetes.NewForConfig(cfg); err != nil {
				continue
			}
			if ctx.Namespace, err = clientcmd_util.NamespaceFromContext(kubeconfigFile, contextName); err != nil {
				continue
			}

			u, err := url.Parse(cfg.Host)
			if err != nil {
				continue
			}
			host := u.Hostname()
			port := u.Port()
			if port == "" {
				if u.Scheme == "https" {
					port = "443"
				} else if u.Scheme == "http" {
					port = "80"
				}
			}
			ctx.Address = host + ":" + port
			s.contexts[contextName] = ctx
		}
	}
	return nil
}

type clusterContext struct {
	Client    kubernetes.Interface
	Namespace string
	Address   string
}

type options struct {
	nsSelector *string // if nil, delete from cluster
	contexts   sets.String
}

func getSyncOptions(annotations map[string]string) options {
	opts := options{}
	if v, err := meta.GetStringValue(annotations, api.ConfigSyncKey); err == nil {
		if v == "true" {
			opts.nsSelector = types.StringP(labels.Everything().String())
		} else {
			opts.nsSelector = &v
		}
	}
	if contexts, _ := meta.GetStringValue(annotations, api.ConfigSyncContexts); contexts != "" {
		opts.contexts = sets.NewString(strings.Split(contexts, ",")...)
	}
	return opts
}

func (s *ConfigSyncer) SyncIntoNamespace(namespace string) error {
	ns, err := s.kubeClient.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
	if err != nil {
		return err
	}

	configMaps, err := s.kubeClient.CoreV1().ConfigMaps(core.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, configMap := range configMaps.Items {
		if err = s.syncConfigMapIntoNewNamespace(&configMap, ns); err != nil {
			return err
		}
	}

	secrets, err := s.kubeClient.CoreV1().Secrets(core.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, secret := range secrets.Items {
		if err = s.syncSecretIntoNewNamespace(&secret, ns); err != nil {
			return err
		}
	}
	return nil
}

func (s *ConfigSyncer) syncerLabels(name, namespace, cluster string) labels.Set {
	return labels.Set{
		api.OriginNameLabelKey:      name,
		api.OriginNamespaceLabelKey: namespace,
		api.OriginClusterLabelKey:   cluster,
	}
}

func (s *ConfigSyncer) syncerLabelSelector(name, namespace, cluster string) string {
	return labels.SelectorFromSet(s.syncerLabels(name, namespace, cluster)).String()
}

func (s *ConfigSyncer) syncerAnnotations(oldAnnotations, srcAnnotations map[string]string, srcRef core.ObjectReference) map[string]string {
	newAnnotations := map[string]string{}

	// preserve sync annotations
	if v, ok := oldAnnotations[api.ConfigSyncKey]; ok {
		newAnnotations[api.ConfigSyncKey] = v
	}
	if v, ok := oldAnnotations[api.ConfigSyncContexts]; ok {
		newAnnotations[api.ConfigSyncContexts] = v
	}

	for k, v := range srcAnnotations {
		if k != api.ConfigSyncKey && k != api.ConfigSyncContexts {
			newAnnotations[k] = v
		}
	}

	// set origin reference
	ref, _ := json.Marshal(srcRef)
	newAnnotations[api.ConfigOriginKey] = string(ref)

	return newAnnotations
}

func (s *ConfigSyncer) namespacesForSelector(selector string) (sets.String, error) {
	namespaces, err := s.kubeClient.CoreV1().Namespaces().List(metav1.ListOptions{
		LabelSelector: selector,
	})
	if err != nil {
		return nil, err
	}
	ns := sets.NewString()
	for _, obj := range namespaces.Items {
		ns.Insert(obj.Name)
	}
	return ns, nil
}
