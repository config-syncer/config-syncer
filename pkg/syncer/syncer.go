package syncer

import (
	"encoding/json"

	"github.com/appscode/kubed/pkg/config"
	"github.com/appscode/kubed/pkg/util"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ConfigSyncer struct {
	KubeClient kubernetes.Interface
}

func (s *ConfigSyncer) SyncConfigMap(oldSrc, newSrc *core.ConfigMap) error {
	var oldSynced, newSynced bool
	if oldSrc != nil {
		oldSynced, _ = util.GetBool(oldSrc.Annotations, config.ConfigSyncKey)
	}
	if newSrc != nil {
		if ok, err := util.GetBool(newSrc.Annotations, config.ConfigSyncKey); err != nil {
			return err // Don't remove by mistake
		} else {
			newSynced = ok
		}
	}
	if newSynced {
		namespaces, err := s.KubeClient.CoreV1().Namespaces().List(metav1.ListOptions{})
		if err != nil {
			return err
		}
		for _, ns := range namespaces.Items {
			s.upsertConfigMap(newSrc, ns.Name)
		}
	} else if oldSynced {
		namespaces, err := s.KubeClient.CoreV1().Namespaces().List(metav1.ListOptions{})
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

func (s *ConfigSyncer) upsertConfigMap(src *core.ConfigMap, namespace string) error {
	ok, err := util.GetBool(src.Annotations, config.ConfigSyncKey)
	if err != nil {
		return err
	}
	if !ok {
		return nil // nothing to sync
	}

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
			if k != config.ConfigSyncKey {
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
		if k != config.ConfigSyncKey {
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

func (s *ConfigSyncer) SyncSecret(oldSrc, newSrc *core.Secret) error {
	var oldSynced, newSynced bool
	if oldSrc != nil {
		oldSynced, _ = util.GetBool(oldSrc.Annotations, config.ConfigSyncKey)
	}
	if newSrc != nil {
		if ok, err := util.GetBool(newSrc.Annotations, config.ConfigSyncKey); err != nil {
			return err // Don't remove by mistake
		} else {
			newSynced = ok
		}
	}
	if newSynced {
		namespaces, err := s.KubeClient.CoreV1().Namespaces().List(metav1.ListOptions{})
		if err != nil {
			return err
		}
		for _, ns := range namespaces.Items {
			s.upsertSecret(newSrc, ns.Name)
		}
	} else if oldSynced {
		namespaces, err := s.KubeClient.CoreV1().Namespaces().List(metav1.ListOptions{})
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

func (s *ConfigSyncer) upsertSecret(src *core.Secret, namespace string) error {
	ok, err := util.GetBool(src.Annotations, config.ConfigSyncKey)
	if err != nil {
		return err
	}
	if !ok {
		return nil // nothing to sync
	}

	if namespace == src.Namespace {
		return nil
	}
	nu, err := s.KubeClient.CoreV1().Secrets(namespace).Get(src.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		// create
		n := *src
		n.Namespace = namespace
		n.UID = ""
		n.ResourceVersion = ""
		n.Annotations = map[string]string{}
		for k, v := range src.Annotations {
			if k != config.ConfigSyncKey {
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
		_, err := s.KubeClient.CoreV1().Secrets(namespace).Create(&n)
		return err
	}

	// update
	nu.Data = src.Data
	nu.Labels = src.Labels
	nu.Annotations = map[string]string{}
	for k, v := range src.Annotations {
		if k != config.ConfigSyncKey {
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
	_, err = s.KubeClient.CoreV1().Secrets(namespace).Update(nu)
	return err
}

func (s *ConfigSyncer) SyncIntoNamespace(namespace string) error {
	_, err := s.KubeClient.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
	if err != nil {
		return err
	}

	cfgmaps, err := s.KubeClient.CoreV1().ConfigMaps(core.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, cfgmap := range cfgmaps.Items {
		s.upsertConfigMap(&cfgmap, namespace)
	}

	secrets, err := s.KubeClient.CoreV1().Secrets(core.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, secret := range secrets.Items {
		s.upsertSecret(&secret, namespace)
	}
	return nil
}
