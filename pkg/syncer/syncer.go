package syncer

import (
	"strings"

	"github.com/appscode/kubed/pkg/config"
	"github.com/appscode/kubed/pkg/util"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ConfigSyncer struct {
	KubeClient         kubernetes.Interface
	ExternalKubeConfig string
}

type syncOpt struct {
	sync       bool
	nsSelector string // should we parse and store as Selector ?
	contexts   []string
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
		if err = s.syncConfigMapIntoNamespace(&configMap, ns); err != nil {
			return err
		}
	}

	secrets, err := s.KubeClient.CoreV1().Secrets(core.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, secret := range secrets.Items {
		if err = s.syncSecretIntoNamespace(&secret, ns); err != nil {
			return err
		}
	}
	return nil
}

func getSyncOption(annotations map[string]string) (opt syncOpt, err error) {
	if util.HasKey(annotations, config.ConfigSyncKey) {
		opt.sync = true
		opt.nsSelector = util.GetString(annotations, config.ConfigSyncKey)
		if opt.nsSelector == "true" {
			opt.nsSelector = ""
		}
	}
	if contexts := util.GetString(annotations, config.ConfigSyncContexts); contexts != "" {
		opt.contexts = strings.Split(contexts, ",")
	}
	return
}
