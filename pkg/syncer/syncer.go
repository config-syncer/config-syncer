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
		s.syncConfigMapIntoNamespace(&configMap, ns) // ignore error ?
	}

	secrets, err := s.KubeClient.CoreV1().Secrets(core.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, secret := range secrets.Items {
		s.syncSecretIntoNamespace(&secret, ns)
	}
	return nil
}

func getSyncOption(annotations map[string]string) (opt syncOpt, err error) {
	if opt.sync, err = util.GetBool(annotations, config.ConfigSyncKey); err != nil {
		return
	}
	opt.nsSelector = util.GetString(annotations, config.ConfigSyncNsSelector)
	contexts := util.GetString(annotations, config.ConfigSyncContexts)
	opt.contexts = strings.Split(contexts, ",")
	return
}
