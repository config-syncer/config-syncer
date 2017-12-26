package syncer

import (
	"encoding/json"

	"github.com/appscode/go/log"
	"github.com/appscode/kubed/pkg/config"
	"github.com/appscode/kubed/pkg/util"
	core_util "github.com/appscode/kutil/core/v1"
	"github.com/appscode/kutil/tools/clientcmd"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
)

func (s *ConfigSyncer) SyncConfigMap(src *core.ConfigMap) error {
	opt, _ := getSyncOption(src.Annotations)

	oldNs, err := util.ConfigMapNamespaceSet(s.KubeClient, config.ConfigOriginKey+"="+src.Name+"."+src.Namespace)
	if err != nil {
		return err
	}

	if opt.sync {
		newNs, err := util.NamespaceSetForSelector(s.KubeClient, opt.nsSelector)
		if err != nil {
			return err
		}
		oldNs = oldNs.Difference(newNs) // delete that were in old but not in new
		oldNs.Delete(src.Namespace)     // skip source
		if err := util.DeleteConfigMapFromNamespaces(s.KubeClient, src.Name, oldNs.List()); err != nil {
			return err
		}

		// upsert to new
		for _, ns := range newNs.List() {
			if err = s.upsertConfigMap(s.KubeClient, src, ns); err != nil {
				return err
			}
		}
	} else { // sync removed, delete that were in old
		oldNs.Delete(src.Namespace) // skip source
		if err := util.DeleteConfigMapFromNamespaces(s.KubeClient, src.Name, oldNs.List()); err != nil {
			return err
		}
	}

	// sync to contexts
	if allContexts, err := util.ContextNameSet(s.KubeConfig); err != nil {
		log.Errorf("Failed to parse context list. Reason: %s\n", err.Error())
	} else {
		for _, context := range allContexts.List() {
			if err = s.syncConfigMapIntoContext(src, context, opt.contexts); err != nil {
				log.Errorf("Failed to sync configmap %s into context %s. Reason: %s\n", src.Name, context, err.Error())
			}
		}
	}

	return nil
}

func (s *ConfigSyncer) SyncDeletedConfigMap(src *core.ConfigMap) error {
	opt, _ := getSyncOption(src.Annotations)

	// sync to namespaces
	if opt.sync {
		if err := s.KubeClient.CoreV1().ConfigMaps(src.Namespace).DeleteCollection(
			&metav1.DeleteOptions{},
			metav1.ListOptions{
				LabelSelector: config.ConfigOriginKey + "=" + src.Name + "." + src.Namespace,
			},
		); err != nil {
			return err
		}
	}

	// sync to contexts
	if allContexts, err := util.ContextNameSet(s.KubeConfig); err != nil {
		log.Errorf("Failed to parse context list. Reason: %s\n", err.Error())
	} else {
		for _, context := range allContexts.List() {
			if err = s.syncConfigMapIntoContext(src, context, sets.NewString()); err != nil {
				log.Errorf("Failed to sync configmap %s into context %s. Reason: %s\n", src.Name, context, err.Error())
			}
		}
	}

	return nil
}

func (s *ConfigSyncer) syncConfigMapIntoNamespace(src *core.ConfigMap, namespace *core.Namespace) error {
	opt, err := getSyncOption(src.Annotations)
	if err != nil {
		return err
	} else if !opt.sync {
		return nil // nothing to sync
	}

	if selector, err := labels.Parse(opt.nsSelector); err != nil {
		return err
	} else if selector.Matches(labels.Set(namespace.Labels)) {
		return s.upsertConfigMap(s.KubeClient, src, namespace.Name)
	} else {

	}

	return nil
}

func (s *ConfigSyncer) syncConfigMapIntoContext(src *core.ConfigMap, context string, contexts sets.String) error {
	client, err := clientcmd.ClientFromContext(s.KubeConfig, context)
	if err != nil {
		return err
	}

	newNs, err := clientcmd.NamespaceFromContext(s.KubeConfig, context)
	if err != nil {
		return err
	}
	if newNs == "" {
		newNs = src.Namespace
	}

	oldNs, err := util.ConfigMapNamespaceSet(client, config.ConfigOriginKey+"="+src.Name+"."+src.Namespace)
	if err != nil {
		return err
	}

	if contexts.Has(context) {
		// in case kubeconfig changes, delete that were in old but not in new
		oldNs.Delete(newNs)
		util.DeleteConfigMapFromNamespaces(client, src.Name, oldNs.List())
		// add to new
		if err = s.upsertConfigMap(client, src, newNs); err != nil {
			return err
		}
	} else { // no sync for this context, delete that were in old
		for _, ns := range oldNs.List() {
			if err = client.CoreV1().ConfigMaps(ns).Delete(src.Name, &metav1.DeleteOptions{}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *ConfigSyncer) upsertConfigMap(k8sClient kubernetes.Interface, src *core.ConfigMap, namespace string) error {
	if namespace == src.Namespace {
		return nil
	}

	meta := metav1.ObjectMeta{
		Name:      src.Name,
		Namespace: namespace,
	}

	_, _, err := core_util.CreateOrPatchConfigMap(k8sClient, meta, func(obj *core.ConfigMap) *core.ConfigMap {
		obj.Data = src.Data
		obj.Labels = src.Labels

		obj.Labels[config.ConfigOriginKey] = src.Name + "." + src.Namespace

		newAnnotations := map[string]string{}
		if metav1.HasAnnotation(obj.ObjectMeta, config.ConfigSyncKey) {
			newAnnotations[config.ConfigSyncKey] = obj.Annotations[config.ConfigSyncKey]
		}
		if metav1.HasAnnotation(obj.ObjectMeta, config.ConfigSyncContexts) {
			newAnnotations[config.ConfigSyncContexts] = obj.Annotations[config.ConfigSyncContexts]
		}
		for k, v := range src.Annotations {
			if k != config.ConfigSyncKey && k != config.ConfigSyncContexts {
				newAnnotations[k] = v
			}
		}
		obj.Annotations = newAnnotations

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
