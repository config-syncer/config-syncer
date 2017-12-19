package framework

import (
	. "github.com/onsi/gomega"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (f *Invocation) EventuallyNumOfConfigmaps(namespace string) GomegaAsyncAssertion {
	return Eventually(func() int {
		cfgMaps, err := f.KubeClient.CoreV1().ConfigMaps(namespace).List(metav1.ListOptions{
			LabelSelector: labels.Set{
				"app": f.App(),
			}.String(),
		})
		Expect(err).NotTo(HaveOccurred())
		return len(cfgMaps.Items)
	})
}

func (f *Invocation) DeleteAllConfigmaps() {
	cfgMaps, err := f.KubeClient.CoreV1().ConfigMaps(metav1.NamespaceAll).List(metav1.ListOptions{
		LabelSelector: labels.Set{
			"app": f.App(),
		}.String(),
	})
	Expect(err).NotTo(HaveOccurred())

	for _, value := range cfgMaps.Items {
		err := f.KubeClient.CoreV1().ConfigMaps(value.Namespace).Delete(value.Name, &metav1.DeleteOptions{})
		if kerr.IsNotFound(err) {
			err = nil
		}
		Expect(err).NotTo(HaveOccurred())
	}
}
