package framework

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/appscode/go/encoding/yaml"
	"github.com/appscode/kubed/pkg/syncer"
	"github.com/appscode/kutil/tools/clientcmd"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

func (f *Invocation) NewConfigMap() *core.ConfigMap {
	return &core.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      f.App(),
			Namespace: f.Namespace(),
			Labels: map[string]string{
				"app": f.App(),
			},
		},
		Data: map[string]string{
			"you":   "only",
			"leave": "once",
		},
	}
}

func (f *Invocation) EventuallyNumOfConfigmaps(namespace string) GomegaAsyncAssertion {
	return f.EventuallyNumOfConfigmapsForClient(f.KubeClient, namespace)
}

func (f *Invocation) EventuallyNumOfConfigmapsForContext(kubeConfigPath string, context string) GomegaAsyncAssertion {
	client, err := clientcmd.ClientFromContext(kubeConfigPath, context)
	Expect(err).ShouldNot(HaveOccurred())
	ns, err := clientcmd.NamespaceFromContext(kubeConfigPath, context)
	Expect(err).ShouldNot(HaveOccurred())
	if ns == "" {
		ns = f.Namespace()
	}

	return f.EventuallyNumOfConfigmapsForClient(client, ns)
}

func (f *Invocation) EventuallyNumOfConfigmapsForClient(client kubernetes.Interface, namespace string) GomegaAsyncAssertion {
	return Eventually(func() int {
		cfgMaps, err := client.CoreV1().ConfigMaps(namespace).List(metav1.ListOptions{
			LabelSelector: labels.Set{
				"app": f.App(),
			}.String(),
		})
		Expect(err).NotTo(HaveOccurred())
		return len(cfgMaps.Items)
	})
}

func (f *Invocation) EventuallyConfigMapSynced(source *core.ConfigMap) GomegaAsyncAssertion {
	opt := syncer.GetSyncOptions(source.Annotations)

	return Eventually(func() bool {

		if opt.NamespaceSelector != nil {
			namespaces, err := syncer.NamespacesForSelector(f.KubeClient, *opt.NamespaceSelector)
			Expect(err).NotTo(HaveOccurred())

			for _, ns := range namespaces.List() {
				if ns == source.Name {
					continue
				}
				_, err := f.KubeClient.CoreV1().ConfigMaps(ns).Get(source.Name, metav1.GetOptions{})
				if err != nil {
					return false
				}
			}
			return true

		} else if opt.Contexts != nil {
			//TODO: Check across context
		}
		return false
	})
}

func (f *Invocation) EventuallyConfigMapNotSynced(source *core.ConfigMap) GomegaAsyncAssertion {

	return Eventually(func() bool {
		namespaces, err := f.KubeClient.CoreV1().Namespaces().List(metav1.ListOptions{})
		Expect(err).NotTo(HaveOccurred())

		for _, ns := range namespaces.Items {
			if ns.Name == source.Name {
				continue
			}
			_, err := f.KubeClient.CoreV1().ConfigMaps(ns.Namespace).Get(source.Name, metav1.GetOptions{})
			if err == nil {
				return false
			}
		}
		return true
	})
}

func (f *Invocation) EventuallyConfigMapSyncedToNamespace(source *core.ConfigMap, namespace string) GomegaAsyncAssertion {
	return Eventually(func() bool {
		_, err := f.KubeClient.CoreV1().ConfigMaps(namespace).Get(source.Name, metav1.GetOptions{})
		if err != nil {
			return false
		}

		return true
	})
}

func (f *Invocation) EventuallySyncedConfigMapsUpdated(source *core.ConfigMap) GomegaAsyncAssertion {
	opt := syncer.GetSyncOptions(source.Annotations)

	return Eventually(func() bool {

		if opt.NamespaceSelector != nil {
			namespaces, err := syncer.NamespacesForSelector(f.KubeClient, *opt.NamespaceSelector)
			Expect(err).NotTo(HaveOccurred())

			for _, ns := range namespaces.List() {
				if ns == source.Namespace {
					continue
				}
				cmReplica, err := f.KubeClient.CoreV1().ConfigMaps(ns).Get(source.Name, metav1.GetOptions{})
				if err != nil {
					return false
				}
				if !reflect.DeepEqual(source.Data, cmReplica.Data) {
					return false
				}
			}
			return true

		} else if opt.Contexts != nil {
			//TODO: Check across context
		}
		return false
	})
}

func (f *Invocation) EventuallySyncedConfigMapsDeleted(source *core.ConfigMap) GomegaAsyncAssertion {
	opt := syncer.GetSyncOptions(source.Annotations)

	return Eventually(func() bool {

		if opt.NamespaceSelector != nil {
			namespaces, err := syncer.NamespacesForSelector(f.KubeClient, *opt.NamespaceSelector)
			Expect(err).NotTo(HaveOccurred())

			for _, ns := range namespaces.List() {
				if ns == source.Namespace {
					continue
				}
				_, err := f.KubeClient.CoreV1().ConfigMaps(ns).Get(source.Name, metav1.GetOptions{})
				if err == nil {
					return false
				}
			}
			return true

		} else if opt.Contexts != nil {
			//TODO: Check across context
		}
		return false
	})
}

func (f *Invocation) ReadConfigMapFromRecycleBin(recycleBinLocation string, cm *core.ConfigMap) (*core.ConfigMap, error) {
	deletedConfigMap := &core.ConfigMap{}
	dir := filepath.Join(recycleBinLocation, filepath.Dir(cm.SelfLink))

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), cm.Name) && strings.HasSuffix(file.Name(), ".yaml") {
			data, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
			if err != nil {
				return nil, err
			}

			err = yaml.Unmarshal(data, &deletedConfigMap)
			if err != nil {
				return nil, err
			}
			return deletedConfigMap, nil
		}
	}
	return deletedConfigMap, fmt.Errorf("configmap not found")
}

func (f *Invocation) DeleteAllConfigmaps() {
	cfgMaps, err := f.KubeClient.CoreV1().ConfigMaps(metav1.NamespaceAll).List(metav1.ListOptions{
		LabelSelector: labels.Set{
			"app": f.App(),
		}.String(),
	})
	Expect(err).NotTo(HaveOccurred())

	for _, value := range cfgMaps.Items {
		err := f.DeleteConfigMap(value.ObjectMeta)
		if kerr.IsNotFound(err) {
			Expect(err).NotTo(HaveOccurred())
		}
	}
}

func (f *Invocation) CreateConfigMap(configMap *core.ConfigMap) (*core.ConfigMap, error) {
	return f.KubeClient.CoreV1().ConfigMaps(configMap.Namespace).Create(configMap)
}
func (f *Invocation) DeleteConfigMap(meta metav1.ObjectMeta) error {
	return f.KubeClient.CoreV1().ConfigMaps(meta.Namespace).Delete(meta.Name, &metav1.DeleteOptions{})
}
