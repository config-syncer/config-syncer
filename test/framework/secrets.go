package framework

import (
	"github.com/appscode/kubed/pkg/util"
	. "github.com/onsi/gomega"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

func (f *Invocation) EventuallyNumOfSecrets(namespace string) GomegaAsyncAssertion {
	return f.EventuallyNumOfSecretsForClient(f.KubeClient, namespace)
}

func (f *Invocation) EventuallyNumOfSecretsForContext(kubeConfigPath string, context string) GomegaAsyncAssertion {
	client, ns, err := util.ClientAndNamespaceForContext(kubeConfigPath, context)
	Expect(err).ShouldNot(HaveOccurred())

	if ns == "" {
		ns = f.Namespace()
	}

	return f.EventuallyNumOfSecretsForClient(client, ns)
}

func (f *Invocation) EventuallyNumOfSecretsForClient(client kubernetes.Interface, namespace string) GomegaAsyncAssertion {
	return Eventually(func() int {
		secrets, err := client.CoreV1().Secrets(namespace).List(metav1.ListOptions{
			LabelSelector: labels.Set{
				"app": f.App(),
			}.String(),
		})
		Expect(err).NotTo(HaveOccurred())
		return len(secrets.Items)
	})
}

func (f *Invocation) DeleteAllSecrets() {
	secrets, err := f.KubeClient.CoreV1().Secrets(metav1.NamespaceAll).List(metav1.ListOptions{
		LabelSelector: labels.Set{
			"app": f.App(),
		}.String(),
	})
	Expect(err).NotTo(HaveOccurred())

	for _, value := range secrets.Items {
		err := f.KubeClient.CoreV1().Secrets(value.Namespace).Delete(value.Name, &metav1.DeleteOptions{})
		if kerr.IsNotFound(err) {
			err = nil
		}
		Expect(err).NotTo(HaveOccurred())
	}
}
