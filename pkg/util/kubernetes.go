package util

import (
	"fmt"
	"strconv"
	"time"

	"net/url"

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

func NamespaceSetForConfigMapSelector(k8sClient kubernetes.Interface, selector string) (sets.String, error) {
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

func NamespaceSetForSecretSelector(k8sClient kubernetes.Interface, selector string) (sets.String, error) {
	secret, err := k8sClient.CoreV1().Secrets(metav1.NamespaceAll).List(metav1.ListOptions{
		LabelSelector: selector,
	})
	if err != nil {
		return nil, err
	}
	ns := sets.NewString()
	for _, obj := range secret.Items {
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

func AddressFromContext(kubeConfigPath, contextName string) (string, error) {
	kConfig, err := clientcmd.LoadFromFile(kubeConfigPath)
	if err != nil {
		return "", err
	}
	ctx, found := kConfig.Contexts[contextName]
	if !found {
		return "", fmt.Errorf("context %s not found in kubeconfig file %s", contextName, kubeConfigPath)
	}
	cluster, found := kConfig.Clusters[ctx.Cluster]
	if !found {
		return "", fmt.Errorf("cluster %s not found in kubeconfig file %s", ctx.Cluster, kubeConfigPath)
	}
	serverUrl, err := url.Parse(cluster.Server)
	if err != nil {
		return "", err
	}
	if serverUrl.Port() == "" {
		if serverUrl.Scheme == "https" {
			return serverUrl.Host + ":443", nil
		} else if serverUrl.Scheme == "http" {
			return serverUrl.Host + ":80", nil
		} else {
			return "", fmt.Errorf("port/scheme not found for context %s", contextName)
		}
	}
	return serverUrl.Host, nil
}
