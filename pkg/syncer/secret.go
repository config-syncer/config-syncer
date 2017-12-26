package syncer

import (
	"encoding/json"

	"github.com/appscode/go/log"
	"github.com/appscode/kubed/pkg/config"
	"github.com/appscode/kubed/pkg/util"
	core_util "github.com/appscode/kutil/core/v1"
	"github.com/appscode/kutil/tools/clientcmd"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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

	if allContexts, err := util.ContextNameSet(s.KubeConfig); err != nil {
		log.Errorf("Failed to parse context list. Reason: %s\n", err.Error())
	} else {
		for _, context := range allContexts.List() {
			if newSrc != nil {
				if err = s.syncSecretIntoContext(newSrc, context, newOpt.contexts); err != nil {
					log.Errorf("Failed to sync secret %s into context %s. Reason: %s\n", newSrc.Name, context, err.Error())
				}
			} else {
				if err = s.syncSecretIntoContext(oldSrc, context, newOpt.contexts); err != nil {
					log.Errorf("Failed to sync secret %s into context %s. Reason: %s\n", oldSrc.Name, context, err.Error())
				}
			}
		}
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

func (s *ConfigSyncer) syncSecretIntoContext(src *core.Secret, context string, newContexts sets.String) error {
	client, err := clientcmd.ClientFromContext(s.KubeConfig, context)
	if err != nil {
		return err
	}

	ns, err := clientcmd.NamespaceFromContext(s.KubeConfig, context)
	if err != nil {
		return err
	}
	if ns == "" {
		ns = src.Namespace
	}

	if newContexts.Has(context) {
		if err = s.upsertSecret(client, src, ns); err != nil {
			return err
		}
	} else {
		if secret, err := client.CoreV1().Secrets(ns).Get(src.Name, metav1.GetOptions{}); err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
		} else if metav1.HasAnnotation(secret.ObjectMeta, config.ConfigOriginKey) { // delete only if it was added by kubed
			if err = client.CoreV1().Secrets(ns).Delete(src.Name, &metav1.DeleteOptions{}); err != nil {
				return err
			}
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
