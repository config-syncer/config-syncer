package framework

import (
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
