package framework

import (
	"github.com/appscode/kubed/pkg/util"
	. "github.com/onsi/gomega"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

func (f *Invocation) EventuallyNumOfConfigmaps(namespace string) GomegaAsyncAssertion {
	return f.EventuallyNumOfConfigmapsForClient(f.KubeClient, namespace)
}

func (f *Invocation) EventuallyNumOfConfigmapsForContext(kubeConfigPath string, context string) GomegaAsyncAssertion {
	client, ns, err := util.ClientAndNamespaceForContext(kubeConfigPath, context)
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
