package util

import (
	"strconv"
	"time"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func IsPreferredAPIResource(kubeClient kubernetes.Interface, groupVersion, kind string) bool {
	if resourceList, err := kubeClient.Discovery().ServerPreferredResources(); err == nil {
		for _, resources := range resourceList {
			if resources.GroupVersion != groupVersion {
				continue
			}
			for _, resource := range resources.APIResources {
				if resources.GroupVersion == groupVersion && resource.Kind == kind {
					return true
				}
			}
		}
	}
	return false
}

const (
	MaxSyncInterval = 5 * time.Minute
)

func IsRecent(t metav1.Time) bool {
	return time.Now().Sub(t.Time) < MaxSyncInterval
}

func ObfuscateSecret(in core.Secret) *core.Secret {
	data := make(map[string][]byte)
	for k := range in.Data {
		data[k] = []byte("-")
	}
	in.Data = data
	return &in
}

func GetBool(m map[string]string, key string) (bool, error) {
	if m == nil {
		return false, nil
	}
	v, ok := m[key]
	if !ok || v == "" {
		return false, nil
	}
	return strconv.ParseBool(v)
}

func GetString(m map[string]string, key string) string {
	if m == nil {
		return ""
	}
	return m[key]
}

func HasKey(m map[string]string, key string) bool {
	if m == nil {
		return false
	}
	_, ok := m[key]
	return ok
}

func RemoveKey(m map[string]string, key string) map[string]string {
	if m == nil {
		return nil
	}
	delete(m, key)
	return m
}

func ContextNameSet(kubeConfigPath string) (sets.String, error) {
	kConfig, err := clientcmd.LoadFromFile(kubeConfigPath)
	if err != nil {
		return nil, err
	}
	return sets.StringKeySet(kConfig.Contexts), nil
}

func ConfigMapNamespaceSet(k8sClient kubernetes.Interface, selector string) (sets.String, error) {
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

func NamespaceSetForSelector(k8sClient kubernetes.Interface, selector string) (sets.String, error) {
	namespaces, err := k8sClient.CoreV1().Namespaces().List(metav1.ListOptions{
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

func DeleteConfigMapFromNamespaces(k8sClient kubernetes.Interface, name string, namespaces []string) error {
	for _, ns := range namespaces {
		if err := k8sClient.CoreV1().ConfigMaps(ns).Delete(name, &metav1.DeleteOptions{}); err != nil {
			return err
		}
	}
	return nil
}
