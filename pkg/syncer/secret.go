package syncer

import (
	"fmt"

	"github.com/appscode/go/log"
	"github.com/appscode/kubed/pkg/util"
	core_util "github.com/appscode/kutil/core/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
)

func (s *ConfigSyncer) SyncSecret(src *core.Secret) error {
	opt := getSyncOption(src.Annotations)

	if opt.sync { // delete that were in old-ns but not in new-ns and upsert to new-ns
		newNs, err := util.NamespaceSetForSelector(s.KubeClient, opt.nsSelector)
		if err != nil {
			return err
		}
		if err := s.syncSecretIntoNamespaces(s.KubeClient, src, newNs, true); err != nil {
			return err
		}
	} else { // no sync, delete that were previously added
		if err := s.syncSecretIntoNamespaces(s.KubeClient, src, sets.NewString(), true); err != nil {
			return err
		}
	}

	return s.syncSecretIntoContexts(src, opt.contexts)
}

// source deleted, delete that were previously added
func (s *ConfigSyncer) SyncDeletedSecret(src *core.Secret) error {
	if err := s.syncSecretIntoNamespaces(s.KubeClient, src, sets.NewString(), true); err != nil {
		return err
	}
	return s.syncSecretIntoContexts(src, sets.NewString())
}

func (s *ConfigSyncer) syncSecretIntoContexts(src *core.Secret, contexts sets.String) error {
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
		if context.Namespace == "" { // use source namespace if not specified via context
			context.Namespace = src.Namespace
		}
		err := s.syncSecretIntoNamespaces(context.Client, src, sets.NewString(context.Namespace), false)
		if err != nil {
			return err
		}
	}

	// delete from other contexts, ignore errors here
	allContexts := sets.StringKeySet(s.Contexts)
	oldContexts := allContexts.Difference(contexts)
	for _, ctx := range oldContexts.List() {
		context, _ := s.Contexts[ctx]
		err := s.syncSecretIntoNamespaces(context.Client, src, sets.NewString(), false)
		if err != nil {
			log.Infoln(err)
		}
	}

	return nil
}

// upsert into newNs set, delete from (oldNs-newNs) set
// use skipSrcNs = true for sync in source cluster
func (s *ConfigSyncer) syncSecretIntoNamespaces(k8sClient kubernetes.Interface, src *core.Secret, newNs sets.String, skipSrcNs bool) error {
	oldNs, err := util.NamespaceSetForSecretSelector(k8sClient, s.SyncerLabelSelector(src.Name, src.Namespace, s.ClusterName))
	if err != nil {
		return err
	}
	oldNs = oldNs.Difference(newNs)
	if skipSrcNs {
		oldNs.Delete(src.Namespace)
		newNs.Delete(src.Namespace)
	}
	for _, ns := range oldNs.List() {
		if err := k8sClient.CoreV1().Secrets(ns).Delete(src.Name, &metav1.DeleteOptions{}); err != nil {
			return err
		}
	}
	for _, ns := range newNs.List() {
		if err = s.upsertSecret(k8sClient, src, ns); err != nil {
			return err
		}
	}
	return nil
}

func (s *ConfigSyncer) syncSecretIntoNewNamespace(src *core.Secret, namespace *core.Namespace) error {
	opt := getSyncOption(src.Annotations)
	if !opt.sync {
		return nil
	}
	if selector, err := labels.Parse(opt.nsSelector); err != nil {
		return err
	} else if selector.Matches(labels.Set(namespace.Labels)) {
		return s.upsertSecret(s.KubeClient, src, namespace.Name)
	}
	return nil
}

func (s *ConfigSyncer) upsertSecret(k8sClient kubernetes.Interface, src *core.Secret, namespace string) error {
	meta := metav1.ObjectMeta{
		Name:      src.Name,
		Namespace: namespace,
	}
	_, _, err := core_util.CreateOrPatchSecret(k8sClient, meta, func(obj *core.Secret) *core.Secret {
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
