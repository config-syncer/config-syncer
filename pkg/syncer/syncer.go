package syncer

import (
	"encoding/json"
	"strings"

	"github.com/appscode/kubed/pkg/config"
	"github.com/appscode/kubed/pkg/util"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
)

type ConfigSyncer struct {
	KubeClient  kubernetes.Interface
	ClusterName string
	Contexts    map[string]ClusterContext
}

type ClusterContext struct {
	Client    kubernetes.Interface
	Namespace string
	Address   string
}

type syncOpt struct {
	sync       bool
	nsSelector string // should we parse and store as Selector ?
	contexts   sets.String
}

func getSyncOption(annotations map[string]string) syncOpt {
	opt := syncOpt{}
	if util.HasKey(annotations, config.ConfigSyncKey) {
		opt.sync = true
		opt.nsSelector = util.GetString(annotations, config.ConfigSyncKey)
		if opt.nsSelector == "true" {
			opt.nsSelector = ""
		}
	}
	if contexts := util.GetString(annotations, config.ConfigSyncContexts); contexts != "" {
		opt.contexts = sets.NewString(strings.Split(contexts, ",")...)
	}
	return opt
}

func (s *ConfigSyncer) SyncIntoNamespace(namespace string) error {
	ns, err := s.KubeClient.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
	if err != nil {
		return err
	}

	configMaps, err := s.KubeClient.CoreV1().ConfigMaps(core.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, configMap := range configMaps.Items {
		if err = s.syncConfigMapIntoNewNamespace(&configMap, ns); err != nil {
			return err
		}
	}

	secrets, err := s.KubeClient.CoreV1().Secrets(core.NamespaceAll).List(metav1.ListOptions{})
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

func (s *ConfigSyncer) SyncerLabels(name, namespace, cluster string) labels.Set {
	return labels.Set{
		config.OriginNameLabelKey:      name,
		config.OriginNamespaceLabelKey: namespace,
		config.OriginClusterLabelKey:   cluster,
	}
}

func (s *ConfigSyncer) SyncerLabelSelector(name, namespace, cluster string) string {
	return labels.SelectorFromSet(s.SyncerLabels(name, namespace, cluster)).String()
}

func (s *ConfigSyncer) SyncerAnnotations(oldAnnotations, srcAnnotations map[string]string, srcRef core.ObjectReference) map[string]string {
	newAnnotations := map[string]string{}

	// preserve sync annotations
	if v, ok := oldAnnotations[config.ConfigSyncKey]; ok {
		newAnnotations[config.ConfigSyncKey] = v
	}
	if v, ok := oldAnnotations[config.ConfigSyncContexts]; ok {
		newAnnotations[config.ConfigSyncContexts] = v
	}

	for k, v := range srcAnnotations {
		if k != config.ConfigSyncKey && k != config.ConfigSyncContexts {
			newAnnotations[k] = v
		}
	}

	// set origin reference
	ref, _ := json.Marshal(srcRef)
	newAnnotations[config.ConfigOriginKey] = string(ref)

	return newAnnotations
}
