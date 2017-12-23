package syncer

import (
	"encoding/json"

	"github.com/appscode/kubed/pkg/config"
	core_util "github.com/appscode/kutil/core/v1"
	"github.com/appscode/kutil/tools/clientcmd"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
)

func (s *ConfigSyncer) SyncSecret(oldSrc, newSrc *core.Secret) error {
	var (
		oldOpt, newOpt syncOpt
		err            error
	)

	if oldSrc != nil {
		oldOpt, _ = getSyncOption(oldSrc.Annotations)
	}
	if newSrc != nil {
		newOpt, err = getSyncOption(newSrc.Annotations)
		if err != nil {
			return err // Don't remove by mistake
		}
	}

	if err = s.syncSecretIntoContexts(newSrc, oldOpt.contexts, newOpt.contexts); err != nil {
		return err
	}

	if newOpt.sync {
		namespaces, err := s.KubeClient.CoreV1().Namespaces().List(metav1.ListOptions{
			LabelSelector: newOpt.nsSelector,
		})
		if err != nil {
			return err
		}
		for _, ns := range namespaces.Items {
			if err = s.upsertSecret(s.KubeClient, newSrc, ns.Name); err != nil {
				return err
			}
		}

		// if selector changed, delete that were in old but not in new (n^2)
		if oldOpt.sync && newOpt.nsSelector != oldOpt.nsSelector {
			oldNamespaces, err := s.KubeClient.CoreV1().Namespaces().List(metav1.ListOptions{
				LabelSelector: oldOpt.nsSelector,
			})
			if err != nil {
				return err
			}
			for _, oldNs := range oldNamespaces.Items {
				if oldNs.Name == newSrc.Namespace {
					continue
				}
				remove := true
				for _, ns := range namespaces.Items {
					if oldNs.Name == ns.Name {
						remove = false
						break
					}
				}
				if remove {
					s.KubeClient.CoreV1().Secrets(oldNs.Name).Delete(newSrc.Name, &metav1.DeleteOptions{})
				}
			}
		}
	} else if oldOpt.sync {
		namespaces, err := s.KubeClient.CoreV1().Namespaces().List(metav1.ListOptions{
			LabelSelector: oldOpt.nsSelector,
		})
		if err != nil {
			return err
		}
		for _, ns := range namespaces.Items {
			if ns.Name == oldSrc.Namespace {
				continue
			}
			s.KubeClient.CoreV1().Secrets(ns.Name).Delete(oldSrc.Name, &metav1.DeleteOptions{})
		}
	}
	return nil
}

func (s *ConfigSyncer) syncSecretIntoNamespace(src *core.Secret, namespace *core.Namespace) error {
	opt, err := getSyncOption(src.Annotations)
	if err != nil {
		return err
	} else if !opt.sync {
		return nil // nothing to sync
	}

	if selector, err := labels.Parse(opt.nsSelector); err != nil {
		return err
	} else if selector.Matches(labels.Set(namespace.Labels)) {
		return s.upsertSecret(s.KubeClient, src, namespace.Name)
	}

	return nil
}

func (s *ConfigSyncer) syncSecretIntoContexts(src *core.Secret, oldContexts, newContexts sets.String) error {
	for _, oldContext := range oldContexts.Difference(newContexts).List() {
		client, err := clientcmd.ClientFromContext(s.KubeConfig, oldContext)
		if err != nil {
			return err
		}
		ns, err := clientcmd.NamespaceFromContext(s.KubeConfig, oldContext)
		if err != nil {
			return err
		}
		if ns == "" {
			ns = src.Namespace
		}
		if err = client.CoreV1().Secrets(ns).Delete(src.Name, &metav1.DeleteOptions{}); err != nil {
			return err
		}
	}

	for _, newContext := range newContexts.List() {
		client, err := clientcmd.ClientFromContext(s.KubeConfig, newContext)
		if err != nil {
			return err
		}
		ns, err := clientcmd.NamespaceFromContext(s.KubeConfig, newContext)
		if err != nil {
			return err
		}
		if ns == "" {
			ns = src.Namespace
		}
		if err = s.upsertSecret(client, src, ns); err != nil {
			return err
		}
	}

	return nil
}

func (s *ConfigSyncer) upsertSecret(k8sClient kubernetes.Interface, src *core.Secret, namespace string) error {
	if namespace == src.Namespace {
		return nil
	}

	meta := metav1.ObjectMeta{
		Name:      src.Name,
		Namespace: namespace,
	}

	_, _, err := core_util.CreateOrPatchSecret(k8sClient, meta, func(obj *core.Secret) *core.Secret {
		obj.Data = src.Data
		obj.Labels = src.Labels

		obj.Annotations = map[string]string{}
		for k, v := range src.Annotations {
			if k != config.ConfigSyncKey && k != config.ConfigSyncContexts {
				obj.Annotations[k] = v
			}
		}

		ref, _ := json.Marshal(core.ObjectReference{
			APIVersion:      src.APIVersion,
			Kind:            src.Kind,
			Name:            src.Name,
			Namespace:       src.Namespace,
			UID:             src.UID,
			ResourceVersion: src.ResourceVersion,
		})
		obj.Annotations[config.ConfigOriginKey] = string(ref)

		return obj
	})

	return err
}
