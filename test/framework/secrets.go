package framework

import (
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (f *Invocation) EventuallyNumOfSecrets(namespace string) GomegaAsyncAssertion {
	return Eventually(func() int {
		secrets, err := f.KubeClient.CoreV1().Secrets(namespace).List(metav1.ListOptions{
			LabelSelector: labels.Set{
				"app": f.App(),
			}.String(),
		})
		Expect(err).NotTo(HaveOccurred())
		return len(secrets.Items)
	})
}
