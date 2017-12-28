package syncer

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/appscode/go/types"
	"github.com/appscode/kubed/pkg/config"
	"github.com/appscode/kutil/meta"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/reference"
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

type options struct {
	nsSelector *string // if nil, delete from cluster
	contexts   sets.String
}

func getSyncOptions(annotations map[string]string) options {
	opts := options{}
	if meta.HasKey(annotations, config.ConfigSyncKey) {
		opts.nsSelector = types.StringP(meta.GetString(annotations, config.ConfigSyncKey))
		if *opts.nsSelector == "true" {
			opts.nsSelector = types.StringP(labels.Everything().String())
		}
	}
	if contexts := meta.GetString(annotations, config.ConfigSyncContexts); contexts != "" {
		opts.contexts = sets.NewString(strings.Split(contexts, ",")...)
	}
	return opts
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

func (s *ConfigSyncer) syncerLabels(name, namespace, cluster string) labels.Set {
	return labels.Set{
		config.OriginNameLabelKey:      name,
		config.OriginNamespaceLabelKey: namespace,
		config.OriginClusterLabelKey:   cluster,
	}
}

func (s *ConfigSyncer) syncerLabelSelector(name, namespace, cluster string) string {
	return labels.SelectorFromSet(s.syncerLabels(name, namespace, cluster)).String()
}

func (s *ConfigSyncer) syncerAnnotations(oldAnnotations, srcAnnotations map[string]string, srcRef core.ObjectReference) map[string]string {
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

func (s *ConfigSyncer) namespacesForSelector(selector string) (sets.String, error) {
	namespaces, err := s.KubeClient.CoreV1().Namespaces().List(metav1.ListOptions{
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

func (s *ConfigSyncer) createEvent(component string, obj runtime.Object, eventType, reason, message string) (*core.Event, error) {
	ref, err := reference.GetReference(scheme.Scheme, obj)
	if err != nil {
		return nil, err
	}

	t := metav1.Time{Time: time.Now()}

	return s.KubeClient.CoreV1().Events(ref.Namespace).Create(&core.Event{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%v.%x", ref.Name, t.UnixNano()),
			Namespace: ref.Namespace,
		},
		InvolvedObject: *ref,
		Reason:         reason,
		Message:        message,
		FirstTimestamp: t,
		LastTimestamp:  t,
		Count:          1,
		Type:           eventType,
		Source:         core.EventSource{Component: component},
	})
}
