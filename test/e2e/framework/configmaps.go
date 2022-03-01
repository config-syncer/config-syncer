/*
Copyright The Config Syncer Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package framework

import (
	"context"
	"path/filepath"
	"reflect"
	"strings"

	"kubeops.dev/config-syncer/pkg/syncer"

	. "github.com/onsi/gomega"
	"gomodules.xyz/encoding/yaml"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"kmodules.xyz/client-go/tools/clientcmd"
	"kmodules.xyz/client-go/tools/exec"
)

func (fi *Invocation) NewConfigMap() *core.ConfigMap {
	return &core.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fi.App(),
			Namespace: fi.Namespace(),
			Labels: map[string]string{
				"app": fi.App(),
			},
		},
		Data: map[string]string{
			"you":  "only",
			"live": "once",
		},
	}
}

func (fi *Invocation) EventuallyNumOfConfigmaps(namespace string) GomegaAsyncAssertion {
	return fi.EventuallyNumOfConfigmapsForClient(fi.KubeClient, namespace)
}

func (fi *Invocation) EventuallyNumOfConfigmapsForContext(kubeConfigPath string, context string) GomegaAsyncAssertion {
	client, err := clientcmd.ClientFromContext(kubeConfigPath, context)
	Expect(err).ShouldNot(HaveOccurred())
	ns, err := clientcmd.NamespaceFromContext(kubeConfigPath, context)
	Expect(err).ShouldNot(HaveOccurred())
	if ns == "" {
		ns = fi.Namespace()
	}

	return fi.EventuallyNumOfConfigmapsForClient(client, ns)
}

func (fi *Invocation) EventuallyNumOfConfigmapsForClient(client kubernetes.Interface, namespace string) GomegaAsyncAssertion {
	return Eventually(func() int {
		cfgMaps, err := client.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: labels.Set{
				"app": fi.App(),
			}.String(),
		})
		Expect(err).NotTo(HaveOccurred())
		return len(cfgMaps.Items)
	})
}

func (fi *Invocation) EventuallyConfigMapSynced(source *core.ConfigMap) GomegaAsyncAssertion {
	opt := syncer.GetSyncOptions(source.Annotations)

	return Eventually(func() bool {
		if opt.NamespaceSelector != nil {
			namespaces, err := syncer.NamespacesForSelector(fi.KubeClient, *opt.NamespaceSelector)
			Expect(err).NotTo(HaveOccurred())

			for _, ns := range namespaces.List() {
				if ns == source.Name {
					continue
				}
				_, err := fi.KubeClient.CoreV1().ConfigMaps(ns).Get(context.TODO(), source.Name, metav1.GetOptions{})
				if err != nil {
					return false
				}
			}
			return true
		}
		return false
	})
}

func (fi *Invocation) EventuallyConfigMapNotSynced(source *core.ConfigMap) GomegaAsyncAssertion {
	return Eventually(func() bool {
		namespaces, err := fi.KubeClient.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		Expect(err).NotTo(HaveOccurred())

		for _, ns := range namespaces.Items {
			if ns.Name == source.Name {
				continue
			}
			_, err := fi.KubeClient.CoreV1().ConfigMaps(ns.Namespace).Get(context.TODO(), source.Name, metav1.GetOptions{})
			if err == nil {
				return false
			}
		}
		return true
	})
}

func (fi *Invocation) EventuallyConfigMapSyncedToNamespace(source *core.ConfigMap, namespace string) GomegaAsyncAssertion {
	return Eventually(func() bool {
		_, err := fi.KubeClient.CoreV1().ConfigMaps(namespace).Get(context.TODO(), source.Name, metav1.GetOptions{})
		return err == nil
	})
}

func (fi *Invocation) EventuallySyncedConfigMapsUpdated(source *core.ConfigMap) GomegaAsyncAssertion {
	opt := syncer.GetSyncOptions(source.Annotations)

	return Eventually(func() bool {
		if opt.NamespaceSelector != nil {
			namespaces, err := syncer.NamespacesForSelector(fi.KubeClient, *opt.NamespaceSelector)
			Expect(err).NotTo(HaveOccurred())

			for _, ns := range namespaces.List() {
				if ns == source.Namespace {
					continue
				}
				cmReplica, err := fi.KubeClient.CoreV1().ConfigMaps(ns).Get(context.TODO(), source.Name, metav1.GetOptions{})
				if err != nil {
					return false
				}
				if !reflect.DeepEqual(source.Data, cmReplica.Data) ||
					!reflect.DeepEqual(source.BinaryData, cmReplica.BinaryData) {
					return false
				}
			}
			return true
		}
		return false
	})
}

func (fi *Invocation) EventuallySyncedConfigMapsDeleted(source *core.ConfigMap) GomegaAsyncAssertion {
	opt := syncer.GetSyncOptions(source.Annotations)

	return Eventually(func() bool {
		if opt.NamespaceSelector != nil {
			namespaces, err := syncer.NamespacesForSelector(fi.KubeClient, *opt.NamespaceSelector)
			Expect(err).NotTo(HaveOccurred())

			for _, ns := range namespaces.List() {
				if ns == source.Namespace {
					continue
				}
				_, err := fi.KubeClient.CoreV1().ConfigMaps(ns).Get(context.TODO(), source.Name, metav1.GetOptions{})
				if err == nil {
					return false
				}
			}
			return true
		}
		return false
	})
}

func (fi *Invocation) ReadConfigMapFromRecycleBin(recycleBinLocation string, cm *core.ConfigMap) (*core.ConfigMap, error) {
	deletedConfigMap := &core.ConfigMap{}
	dir := filepath.Join(recycleBinLocation, filepath.Dir(cm.SelfLink))

	pod, err := fi.OperatorPod()
	if err != nil {
		return nil, err
	}

	// list the name of recycled configMaps
	output, err := exec.ExecIntoPod(fi.ClientConfig, pod, exec.Command("ls", dir))
	if err != nil {
		return nil, err
	}

	// read the recycled configMaps
	var recycledCMName string
	fileNames := strings.Split(output, "\n")
	for _, fileName := range fileNames {
		if strings.HasPrefix(fileName, cm.Name) && strings.HasSuffix(fileName, ".yaml") {
			recycledCMName = fileName
			break
		}
	}
	output, err = exec.ExecIntoPod(fi.ClientConfig, pod, exec.Command("cat", filepath.Join(dir, recycledCMName)))
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal([]byte(output), &deletedConfigMap)
	if err != nil {
		return nil, err
	}
	return deletedConfigMap, nil
}

func (fi *Invocation) DeleteAllConfigmaps() {
	cfgMaps, err := fi.KubeClient.CoreV1().ConfigMaps(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labels.Set{
			"app": fi.App(),
		}.String(),
	})
	Expect(err).NotTo(HaveOccurred())

	for _, value := range cfgMaps.Items {
		err := fi.DeleteConfigMap(value.ObjectMeta)
		if kerr.IsNotFound(err) {
			Expect(err).NotTo(HaveOccurred())
		}
	}
}

func (fi *Invocation) CreateConfigMap(configMap *core.ConfigMap) (*core.ConfigMap, error) {
	return fi.KubeClient.CoreV1().ConfigMaps(configMap.Namespace).Create(context.TODO(), configMap, metav1.CreateOptions{})
}

func (fi *Invocation) DeleteConfigMap(meta metav1.ObjectMeta) error {
	return fi.KubeClient.CoreV1().ConfigMaps(meta.Namespace).Delete(context.TODO(), meta.Name, metav1.DeleteOptions{})
}
