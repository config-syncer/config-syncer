package framework

import (
	"reflect"

	"github.com/appscode/go/crypto/rand"
	api "github.com/appscode/kubed/apis/kubed/v1alpha1"
	"github.com/appscode/kubed/pkg/syncer"
	"github.com/appscode/kutil/tools/clientcmd"
	"github.com/ghodss/yaml"
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

func (f *Invocation) EventuallySecretSynced(source *core.Secret) GomegaAsyncAssertion {
	opt := syncer.GetSyncOptions(source.Annotations)

	return Eventually(func() bool {

		if opt.NamespaceSelector != nil {
			namespaces, err := syncer.NamespacesForSelector(f.KubeClient, *opt.NamespaceSelector)
			Expect(err).NotTo(HaveOccurred())

			for _, ns := range namespaces.List() {
				if ns == source.Name {
					continue
				}
				_, err := f.KubeClient.CoreV1().Secrets(ns).Get(source.Name, metav1.GetOptions{})
				if err != nil {
					return false
				}
			}
			return true

		} else if opt.Contexts != nil {
			//TODO: Check across context
		}
		return false
	})
}

func (f *Invocation) EventuallySecretNotSynced(source *core.Secret) GomegaAsyncAssertion {

	return Eventually(func() bool {
		namespaces, err := f.KubeClient.CoreV1().Namespaces().List(metav1.ListOptions{})
		Expect(err).NotTo(HaveOccurred())

		for _, ns := range namespaces.Items {
			if ns.Name == source.Name {
				continue
			}
			_, err := f.KubeClient.CoreV1().Secrets(ns.Namespace).Get(source.Name, metav1.GetOptions{})
			if err == nil {
				return false
			}
		}
		return true
	})
}

func (f *Invocation) EventuallySecretSyncedToNamespace(source *core.Secret, namespace string) GomegaAsyncAssertion {
	return Eventually(func() bool {
		_, err := f.KubeClient.CoreV1().Secrets(namespace).Get(source.Name, metav1.GetOptions{})
		if err != nil {
			return false
		}

		return true
	})
}

func (f *Invocation) EventuallySyncedSecretsUpdated(source *core.Secret) GomegaAsyncAssertion {
	opt := syncer.GetSyncOptions(source.Annotations)

	return Eventually(func() bool {

		if opt.NamespaceSelector != nil {
			namespaces, err := syncer.NamespacesForSelector(f.KubeClient, *opt.NamespaceSelector)
			Expect(err).NotTo(HaveOccurred())

			for _, ns := range namespaces.List() {
				if ns == source.Namespace {
					continue
				}
				secretReplica, err := f.KubeClient.CoreV1().Secrets(ns).Get(source.Name, metav1.GetOptions{})
				if err != nil {
					return false
				}
				if !reflect.DeepEqual(source.Data, secretReplica.Data) {
					return false
				}
			}
			return true

		} else if opt.Contexts != nil {
			//TODO: Check across context
		}
		return false
	})
}

func (f *Invocation) EventuallySyncedSecretsDeleted(source *core.Secret) GomegaAsyncAssertion {
	opt := syncer.GetSyncOptions(source.Annotations)

	return Eventually(func() bool {

		if opt.NamespaceSelector != nil {
			namespaces, err := syncer.NamespacesForSelector(f.KubeClient, *opt.NamespaceSelector)
			Expect(err).NotTo(HaveOccurred())

			for _, ns := range namespaces.List() {
				if ns == source.Namespace {
					continue
				}
				_, err := f.KubeClient.CoreV1().Secrets(ns).Get(source.Name, metav1.GetOptions{})
				if err == nil {
					return false
				}
			}
			return true

		} else if opt.Contexts != nil {
			//TODO: Check across context
		}
		return false
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

func (f *Framework) CreateSecret(obj *core.Secret) (*core.Secret, error) {
	return f.KubeClient.CoreV1().Secrets(obj.Namespace).Create(obj)
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
		secret.Data[api.CA_CERT_DATA] = f.CertStore.CACertBytes()
	}

	return secret
}

func (fi *Invocation) SecretForWebhookNotifier() *core.Secret {
	namespace := fi.namespace
	if fi.SelfHostedOperator {
		namespace = OperatorNamespace
	}
	return &core.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: metav1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "notifier-config",
			Namespace: namespace,
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

func (f *Framework) KubeConfigSecret(config *api.ClusterConfig, meta metav1.ObjectMeta) (*core.Secret, error) {
	data, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	return &core.Secret{
		ObjectMeta: meta,
		Data: map[string][]byte{
			"config.yaml": data,
		},
	}, nil
}
