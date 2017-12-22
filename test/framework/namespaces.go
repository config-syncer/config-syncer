package framework

import (
	"github.com/appscode/kubed/pkg/util"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (f *Framework) DeleteNamespace() error {
	return f.KubeClient.CoreV1().Namespaces().Delete(f.namespace, &metav1.DeleteOptions{})
}

func (f *Framework) EventuallyNamespaceDeleted(ns string) GomegaAsyncAssertion {
	return Eventually(func() bool {
		_, err := f.KubeClient.CoreV1().Namespaces().Get(ns, metav1.GetOptions{})
		return kerr.IsNotFound(err)
	})
}

func (f *Invocation) EnsureNamespaceForContext(kubeConfigPath string, context string) {
	client, ns, err := util.ClientAndNamespaceForContext(kubeConfigPath, context)
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
	client, ns, err := util.ClientAndNamespaceForContext(kubeConfigPath, context)
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
