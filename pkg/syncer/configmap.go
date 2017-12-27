package syncer

import (
	"github.com/appscode/go/log"
	"github.com/appscode/kubed/pkg/util"
	core_util "github.com/appscode/kutil/core/v1"
	"github.com/appscode/kutil/tools/clientcmd"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	"fmt"
)

func (s *ConfigSyncer) SyncConfigMap(src *core.ConfigMap) error {
	opt := getSyncOption(src.Annotations)

	oldNs, err := util.ConfigMapNamespaceSet(s.KubeClient, s.SyncerLabelSelector(src.Name, src.Namespace, s.ClusterName))
	if err != nil {
		return err
	}
	oldNs.Delete(src.Namespace) // skip source-ns

	if opt.sync {
		// delete that were in old-ns but not in new-ns
		newNs, err := util.NamespaceSetForSelector(s.KubeClient, opt.nsSelector)
		if err != nil {
			return err
		}
		oldNs = oldNs.Difference(newNs)
		if err := util.DeleteConfigMapFromNamespaces(s.KubeClient, src.Name, oldNs); err != nil {
			return err
		}

		// upsert to new-ns
		for _, ns := range newNs.List() {
			if err = s.upsertConfigMap(s.KubeClient, src, ns); err != nil {
				return err
			}
		}
	} else { // no sync, delete that were previously added
		if err := util.DeleteConfigMapFromNamespaces(s.KubeClient, src.Name, oldNs); err != nil {
			return err
		}
	}

	s.syncConfigMapIntoAllContexts(src, opt.contexts)

	return nil
}

func (s *ConfigSyncer) SyncDeletedConfigMap(src *core.ConfigMap) error {
	// source deleted, delete that were previously added
	oldNs, err := util.ConfigMapNamespaceSet(s.KubeClient, s.SyncerLabelSelector(src.Name, src.Namespace, s.ClusterName))
	if err != nil {
		return err
	}
	if err := util.DeleteConfigMapFromNamespaces(s.KubeClient, src.Name, oldNs); err != nil {
		return err
	}

	s.syncConfigMapIntoAllContexts(src, sets.NewString())

	return nil
}

func (s *ConfigSyncer) syncConfigMapIntoNewNamespace(src *core.ConfigMap, namespace *core.Namespace) error {
	opt := getSyncOption(src.Annotations)
	if !opt.sync {
		return nil
	}

	if selector, err := labels.Parse(opt.nsSelector); err != nil {
		return err
	} else if selector.Matches(labels.Set(namespace.Labels)) {
		return s.upsertConfigMap(s.KubeClient, src, namespace.Name)
	} else {

	}

	return nil
}

// upsert into newNs set, delete from (oldNs-newNs) set, skip srcNs
func (s *ConfigSyncer) syncConfigMapIntoNamespaces(k8sClient kubernetes.Interface, src *core.ConfigMap, newNs sets.String, skipSource bool) error {
	oldNs, err := util.ConfigMapNamespaceSet(k8sClient, s.SyncerLabelSelector(src.Name, src.Namespace, s.ClusterName))
	if err != nil {
		return err
	}
	oldNs = oldNs.Difference(newNs)
	if skipSource{
		oldNs.Delete(src.Namespace)
		newNs.Delete(src.Namespace)
	}
	if err := util.DeleteConfigMapFromNamespaces(k8sClient, src.Name, oldNs); err != nil {
		return err
	}
	for _, ns := range newNs.List() {
		if err = s.upsertConfigMap(k8sClient, src, ns); err != nil {
			return err
		}
	}
	return nil
}

func (s *ConfigSyncer) syncConfigMapIntoAllContexts(src *core.ConfigMap, contexts sets.String) error {
	// validate contexts specified via annotation
	taken := map[string]struct{}{}
	for _, ctx := range contexts.List() {
		context, found := s.Contexts[ctx]
		if !found {
			return fmt.Errorf("context %s not found in kubeconfig file", ctx)
		}
		if _, found = taken[context.Address]; found {
			return fmt.Errorf("multiple contexts poniting same cluster")
		}
		taken[context.Address] = struct{}{}
	}

	// sync to contexts specified via annotation, do not ignore errors here
	for _, ctx := range contexts.List() {
		context, _ := s.Contexts[ctx]

		// delete that were in old-ns but not in new-ns
		oldNs, err := util.ConfigMapNamespaceSet(context.Client, s.SyncerLabelSelector(src.Name, src.Namespace, s.ClusterName))
		if err != nil {
			return err
		}
		oldNs.Delete(context.Namespace)
		if err = util.DeleteConfigMapFromNamespaces(context.Client, src.Name, oldNs); err != nil {
			return err
		}

		// upsert to new
		if err = s.upsertConfigMap(context.Client, src, context.Namespace); err != nil {
			return err
		}
	}

	// delete from other contexts, ignore errors here
	allContexts := sets.StringKeySet(s.Contexts)
	oldContexts := allContexts.Difference(contexts)
	for _, ctx := range oldContexts.List() {
		context, _ := s.Contexts[ctx]
		oldNs, err := util.ConfigMapNamespaceSet(context.Client, s.SyncerLabelSelector(src.Name, src.Namespace, s.ClusterName))
		if err != nil {
			log.Infoln(err)
		}
		if err = util.DeleteConfigMapFromNamespaces(context.Client, src.Name, oldNs); err != nil {
			log.Infoln(err)
		}
	}
	return nil
}

func (s *ConfigSyncer) syncConfigMapIntoContext(src *core.ConfigMap, context string, contexts sets.String) error {
	client, err := clientcmd.ClientFromContext(s.KubeConfig, context)
	if err != nil {
		return err
	}

	oldNs, err := util.ConfigMapNamespaceSet(client, s.SyncerLabelSelector(src.Name, src.Namespace, s.ClusterName))
	if err != nil {
		return err
	}

	if contexts.Has(context) {
		// in case kubeconfig changes, delete that were in old but not in new
		newNs, err := clientcmd.NamespaceFromContext(s.KubeConfig, context)
		if err != nil {
			return err
		}
		if newNs == "" {
			newNs = src.Namespace
		}
		oldNs.Delete(newNs)
		if err = util.DeleteConfigMapFromNamespaces(client, src.Name, oldNs); err != nil {
			return err
		}
		// upsert to new
		if err = s.upsertConfigMap(client, src, newNs); err != nil {
			return err
		}
	} else { // no sync, delete that were previously added
		if err = util.DeleteConfigMapFromNamespaces(client, src.Name, oldNs); err != nil {
			return err
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
		obj.Labels = labels.Merge(src.Labels, s.SyncerLabels(src.Name, src.Namespace, s.ClusterName))

		ref := core.ObjectReference{
			APIVersion:      src.APIVersion,
			Kind:            src.Kind,
			Name:            src.Name,
			Namespace:       src.Namespace,
			UID:             src.UID,
			ResourceVersion: src.ResourceVersion,
		}
		obj.Annotations = s.SyncerAnnotations(obj.Annotations, src.Annotations, ref)

		return obj
	})

	return err
}
