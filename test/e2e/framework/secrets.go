package framework

import (
	"github.com/appscode/go/crypto/rand"
	api "github.com/appscode/kubed/apis/kubed/v1alpha1"
	"github.com/appscode/kutil/tools/clientcmd"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

func (f *Invocation) NewSecret() *core.Secret {
	return &core.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      f.App(),
			Namespace: f.Namespace(),
			Labels: map[string]string{
				"app": f.App(),
			},
		},
		StringData: map[string]string{
			"you":   "only",
			"leave": "once",
		},
	}
}

func (f *Invocation) EventuallyNumOfSecrets(namespace string) GomegaAsyncAssertion {
	return f.EventuallyNumOfSecretsForClient(f.KubeClient, namespace)
}

func (f *Invocation) EventuallyNumOfSecretsForContext(kubeConfigPath string, context string) GomegaAsyncAssertion {
	client, err := clientcmd.ClientFromContext(kubeConfigPath, context)
	Expect(err).ShouldNot(HaveOccurred())
	ns, err := clientcmd.NamespaceFromContext(kubeConfigPath, context)
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

func (f *Framework) CreateSecret(obj core.Secret) error {
	_, err := f.KubeClient.CoreV1().Secrets(obj.Namespace).Create(&obj)
	return err
}

func (f *Framework) DeleteSecret(meta metav1.ObjectMeta) error {
	return f.KubeClient.CoreV1().Secrets(meta.Namespace).Delete(meta.Name, deleteInForeground())
}

func (f *Invocation) SecretForMinioBackend(includeCert bool) core.Secret {
	secret := core.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rand.WithUniqSuffix(f.app + "-minio"),
			Namespace: f.namespace,
		},
		Data: map[string][]byte{
			api.AWS_ACCESS_KEY_ID:     []byte(MINIO_ACCESS_KEY_ID),
			api.AWS_SECRET_ACCESS_KEY: []byte(MINIO_SECRET_ACCESS_KEY),
		},
	}

	if includeCert {
		secret.Data[api.CA_CERT_DATA] = f.CertStore.CACert()
	}

	return secret
}

func (fi *Invocation) SecretForWebhookNotifier() *core.Secret {
	return &core.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: metav1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "notifier-config",
			Namespace: fi.namespace,
		},
		Data: map[string][]byte{
			"WEBHOOK_URL": []byte("http://localhost:8181"),
		},
	}
}

func (f *Invocation) WaitUntilSecretCreated(meta metav1.ObjectMeta) error {
	return wait.PollImmediate(interval, timeout, func() (done bool, err error) {
		if _, err := f.KubeClient.CoreV1().Secrets(meta.Namespace).Get(meta.Name, metav1.GetOptions{}); err != nil {
			if kerr.IsNotFound(err) {
				return false, nil
			} else {
				return true, err
			}
		}
		return true, nil
	})
}

func (f *Invocation) WaitUntilSecretDeleted(meta metav1.ObjectMeta) error {
	return wait.PollImmediate(interval, timeout, func() (done bool, err error) {
		if _, err := f.KubeClient.CoreV1().Secrets(meta.Namespace).Get(meta.Name, metav1.GetOptions{}); err != nil {
			if kerr.IsNotFound(err) {
				return true, nil
			} else {
				return true, err
			}
		}
		return false, nil
	})
}
