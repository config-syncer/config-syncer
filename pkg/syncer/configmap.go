package syncer

import (
	"fmt"

	"github.com/appscode/go/log"
	"github.com/appscode/kubed/pkg/api"
	"github.com/appscode/kubed/pkg/eventer"
	core_util "github.com/appscode/kutil/core/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
)

func (s *ConfigSyncer) SyncConfigMap(src *core.ConfigMap) error {
	opts := getSyncOptions(src.Annotations)

	if opts.nsSelector != nil { // delete that were in old-ns but not in new-ns and upsert to new-ns
		newNs, err := s.namespacesForSelector(*opts.nsSelector)
		if err != nil {
			return err
		}
		if err := s.syncConfigMapIntoNamespaces(s.kubeClient, src, newNs, true, ""); err != nil {
			return err
		}
	} else { // no sync, delete that were previously added
		if err := s.syncConfigMapIntoNamespaces(s.kubeClient, src, sets.NewString(), true, ""); err != nil {
			return err
		}
	}

	return s.syncConfigMapIntoContexts(src, opts.contexts)
}

// source deleted, delete that were previously added
func (s *ConfigSyncer) SyncDeletedConfigMap(src *core.ConfigMap) error {
	if err := s.syncConfigMapIntoNamespaces(s.kubeClient, src, sets.NewString(), true, ""); err != nil {
		return err
	}
	return s.syncConfigMapIntoContexts(src, sets.NewString())
}

func (s *ConfigSyncer) syncConfigMapIntoContexts(src *core.ConfigMap, contexts sets.String) error {
	// validate contexts specified via annotation
	taken := map[string]struct{}{}
	for _, ctx := range contexts.List() {
		context, found := s.contexts[ctx]
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
		context, _ := s.contexts[ctx]
		if context.Namespace == "" { // use source namespace if not specified via context
			context.Namespace = src.Namespace
		}
		err := s.syncConfigMapIntoNamespaces(context.Client, src, sets.NewString(context.Namespace), false, ctx)
		if err != nil {
			return err
		}
	}

	// delete from other contexts, ignore errors here
	for ctxName, ctx := range s.contexts {
		if _, found := taken[ctx.Address]; !found {
			err := s.syncConfigMapIntoNamespaces(ctx.Client, src, sets.NewString(), false, ctxName)
			if err != nil {
				log.Infoln(err)
			}
			taken[ctx.Address] = struct{}{} // to avoid deleting form same cluster twice
		}
	}

	return nil
}

// upsert into newNs set, delete from (oldNs-newNs) set
// use skipSrcNs = true for sync in source cluster
func (s *ConfigSyncer) syncConfigMapIntoNamespaces(k8sClient kubernetes.Interface, src *core.ConfigMap, newNs sets.String, skipSrcNs bool, context string) error {
	oldNs, err := namespaceSetForConfigMapSelector(k8sClient, s.syncerLabelSelector(src.Name, src.Namespace, s.clusterName))
	if err != nil {
		return err
	}
	oldNs = oldNs.Difference(newNs)
	if skipSrcNs {
		oldNs.Delete(src.Namespace)
		newNs.Delete(src.Namespace)
	}
	for _, ns := range oldNs.List() {
		if err := k8sClient.CoreV1().ConfigMaps(ns).Delete(src.Name, &metav1.DeleteOptions{}); err != nil {
			return err
		}
	}
	for _, ns := range newNs.List() {
		if err = s.upsertConfigMap(k8sClient, src, ns, context); err != nil {
			return err
		}
	}
	return nil
}

func (s *ConfigSyncer) syncConfigMapIntoNewNamespace(src *core.ConfigMap, namespace *core.Namespace) error {
	opts := getSyncOptions(src.Annotations)
	if opts.nsSelector == nil {
		return nil
	}
	if selector, err := labels.Parse(*opts.nsSelector); err != nil {
		return err
	} else if selector.Matches(labels.Set(namespace.Labels)) {
		return s.upsertConfigMap(s.kubeClient, src, namespace.Name, "")
	}
	return nil
}

func (s *ConfigSyncer) upsertConfigMap(k8sClient kubernetes.Interface, src *core.ConfigMap, namespace, context string) error {
	meta := metav1.ObjectMeta{
		Name:      src.Name,
		Namespace: namespace,
	}
	_, _, err := core_util.CreateOrPatchConfigMap(k8sClient, meta, func(obj *core.ConfigMap) *core.ConfigMap {
		// check origin cluster, if not match overwrite and create an event
		if v, ok := obj.Labels[api.OriginClusterLabelKey]; ok && v != s.clusterName {
			s.recorder.Eventf(
				src,
				core.EventTypeWarning,
				eventer.EventReasonOriginConflict,
				"Origin cluster changed from %s in context %s", v, context,
			)
		}

		obj.Data = src.Data
		obj.Labels = labels.Merge(src.Labels, s.syncerLabels(src.Name, src.Namespace, s.clusterName))

		ref := core.ObjectReference{
			APIVersion:      src.APIVersion,
			Kind:            src.Kind,
			Name:            src.Name,
			Namespace:       src.Namespace,
			UID:             src.UID,
			ResourceVersion: src.ResourceVersion,
		}
		obj.Annotations = s.syncerAnnotations(obj.Annotations, src.Annotations, ref)

		return obj
	})

	return err
}

func namespaceSetForConfigMapSelector(k8sClient kubernetes.Interface, selector string) (sets.String, error) {
	cfgMaps, err := k8sClient.CoreV1().ConfigMaps(metav1.NamespaceAll).List(metav1.ListOptions{
		LabelSelector: selector,
	})
	if err != nil {
		return nil, err
	}
	ns := sets.NewString()
	for _, obj := range cfgMaps.Items {
		ns.Insert(obj.Namespace)
	}
	return ns, nil
}
