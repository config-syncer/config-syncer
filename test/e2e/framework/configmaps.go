package framework

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/go/encoding/yaml"
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

func (f *Invocation) NewNamespaceWithLabel() *core.Namespace {
	return &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: rand.WithUniqSuffix("kubed-e2e-labeled"),
			Labels: map[string]string{
				"app": f.App(),
			},
		},
	}
}

func (f *Invocation) NumberOfNameSpace() int {
	ns, err := f.KubeClient.CoreV1().Namespaces().List(metav1.ListOptions{})
	Expect(err).NotTo(HaveOccurred())
	return len(ns.Items)
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
