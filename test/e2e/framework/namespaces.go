/*
Copyright The Kubed Authors.

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
	"github.com/appscode/go/crypto/rand"

	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kmodules.xyz/client-go/tools/clientcmd"
)

func (f *Framework) Namespace() string {
	return f.namespace
}

func (f *Framework) EnsureNamespace() error {
	_, err := f.KubeClient.CoreV1().Namespaces().Get(f.namespace, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err = f.KubeClient.CoreV1().Namespaces().Create(&core.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: f.namespace,
			},
		})
	}
	return err
}

func (f *Framework) NewNamespace(name string) *core.Namespace {
	return &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: rand.WithUniqSuffix(name),
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

func (f *Framework) CreateNamespace(ns *core.Namespace) error {
	_, err := f.KubeClient.CoreV1().Namespaces().Create(ns)
	return err
}
func (f *Framework) DeleteNamespace(name string) error {
	return f.KubeClient.CoreV1().Namespaces().Delete(name, &metav1.DeleteOptions{})
}

func (f *Framework) EventuallyNamespaceDeleted(ns string) GomegaAsyncAssertion {
	return Eventually(func() bool {
		_, err := f.KubeClient.CoreV1().Namespaces().Get(ns, metav1.GetOptions{})
		return kerr.IsNotFound(err)
	})
}

func (f *Invocation) EnsureNamespaceForContext(kubeConfigPath string, context string) {
	client, err := clientcmd.ClientFromContext(kubeConfigPath, context)
	Expect(err).ShouldNot(HaveOccurred())
	ns, err := clientcmd.NamespaceFromContext(kubeConfigPath, context)
	Expect(err).ShouldNot(HaveOccurred())

	if ns == "" {
		ns = f.Namespace()
	}

	_, err = client.CoreV1().Namespaces().Get(ns, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err = client.CoreV1().Namespaces().Create(&core.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: ns,
			},
		})
		Expect(err).ShouldNot(HaveOccurred())

		Eventually(func() bool {
			_, err := client.CoreV1().Namespaces().Get(ns, metav1.GetOptions{})
			return kerr.IsNotFound(err)
		}).Should(BeFalse())
	}
}

func (f *Invocation) DeleteNamespaceForContext(kubeConfigPath string, context string) {
	client, err := clientcmd.ClientFromContext(kubeConfigPath, context)
	Expect(err).ShouldNot(HaveOccurred())
	ns, err := clientcmd.NamespaceFromContext(kubeConfigPath, context)
	Expect(err).ShouldNot(HaveOccurred())

	if ns == "" {
		ns = f.Namespace()
	}

	err = client.CoreV1().Namespaces().Delete(ns, &metav1.DeleteOptions{})
	if errors.IsNotFound(err) {
		err = nil
	}
	Expect(err).ShouldNot(HaveOccurred())

	Eventually(func() bool {
		_, err := client.CoreV1().Namespaces().Get(ns, metav1.GetOptions{})
		return kerr.IsNotFound(err)
	}).Should(BeTrue())
}
