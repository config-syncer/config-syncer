package syncer

import (
	"encoding/json"

	"github.com/appscode/kubed/pkg/config"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (s *ConfigSyncer) SyncConfigMap(oldSrc, newSrc *core.ConfigMap) error {
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

	if newOpt.sync {
		namespaces, err := s.KubeClient.CoreV1().Namespaces().List(metav1.ListOptions{
			LabelSelector: newOpt.nsSelector,
		})
		if err != nil {
			return err
		}
		for _, ns := range namespaces.Items {
			s.upsertConfigMap(newSrc, ns.Name)
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
					s.KubeClient.CoreV1().ConfigMaps(oldNs.Name).Delete(newSrc.Name, &metav1.DeleteOptions{})
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
			s.KubeClient.CoreV1().ConfigMaps(ns.Name).Delete(oldSrc.Name, &metav1.DeleteOptions{})
		}
	}
	return nil
}

func (s *ConfigSyncer) SyncConfigMapIntoNamespace(src *core.ConfigMap, namespace *core.Namespace) error {
	opt, err := getSyncOption(src.Annotations)
	if err != nil {
		return err
	} else if !opt.sync {
		return nil // nothing to sync
	}

	if selector, err := labels.Parse(opt.nsSelector); err != nil {
		return err
	} else if selector.Matches(labels.Set(namespace.Labels)) {
		return s.upsertConfigMap(src, namespace.Name)
	}

	return nil
}

func (s *ConfigSyncer) upsertConfigMap(src *core.ConfigMap, namespace string) error {
	if namespace == src.Namespace {
		return nil
	}
	nu, err := s.KubeClient.CoreV1().ConfigMaps(namespace).Get(src.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		// create
		n := *src
		n.Namespace = namespace
		n.UID = ""
		n.ResourceVersion = ""
		n.Annotations = map[string]string{}
		for k, v := range src.Annotations {
			if k != config.ConfigSyncKey && k != config.ConfigSyncNsSelector {
				n.Annotations[k] = v
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
		n.Annotations[config.ConfigOriginKey] = string(ref)

		_, err := s.KubeClient.CoreV1().ConfigMaps(namespace).Create(&n)
		return err
	}
	// update
	nu.Data = src.Data
	nu.Labels = src.Labels
	nu.Annotations = map[string]string{}
	for k, v := range src.Annotations {
		if k != config.ConfigSyncKey && k != config.ConfigSyncNsSelector {
			nu.Annotations[k] = v
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
	nu.Annotations[config.ConfigOriginKey] = string(ref)
	_, err = s.KubeClient.CoreV1().ConfigMaps(namespace).Update(nu)
	return err
}
