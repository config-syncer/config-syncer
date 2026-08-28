/*
Copyright The Config Syncer Authors.

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
	"context"
	"reflect"

	"kubeops.dev/config-syncer/pkg/syncer"

	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	kutil "kmodules.xyz/client-go"
	"kmodules.xyz/client-go/tools/clientcmd"
)

func (fi *Invocation) NewSecret() *core.Secret {
	return &core.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fi.App(),
			Namespace: fi.Namespace(),
			Labels: map[string]string{
				"app": fi.App(),
			},
		},
		StringData: map[string]string{
			"you":  "only",
			"live": "once",
		},
	}
}

func (fi *Invocation) EventuallyNumOfSecrets(namespace string) GomegaAsyncAssertion {
	return fi.EventuallyNumOfSecretsForClient(fi.KubeClient, namespace)
}

func (fi *Invocation) EventuallyNumOfSecretsForContext(kubeConfigPath string, context string) GomegaAsyncAssertion {
	client, err := clientcmd.ClientFromContext(kubeConfigPath, context)
	Expect(err).ShouldNot(HaveOccurred())
	ns, err := clientcmd.NamespaceFromContext(kubeConfigPath, context)
	Expect(err).ShouldNot(HaveOccurred())

	if ns == "" {
		ns = fi.Namespace()
	}

	return fi.EventuallyNumOfSecretsForClient(client, ns)
}

func (fi *Invocation) EventuallyNumOfSecretsForClient(client kubernetes.Interface, namespace string) GomegaAsyncAssertion {
	return Eventually(func() int {
		secrets, err := client.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: labels.Set{
				"app": fi.App(),
			}.String(),
		})
		Expect(err).NotTo(HaveOccurred())
		return len(secrets.Items)
	})
}

func (fi *Invocation) EventuallySecretSynced(source *core.Secret) GomegaAsyncAssertion {
	opt := syncer.GetSyncOptions(source.Annotations)

	return Eventually(func() bool {
		if opt.NamespaceSelector != nil {
			namespaces, err := syncer.NamespacesForSelector(fi.KubeClient, *opt.NamespaceSelector)
			Expect(err).NotTo(HaveOccurred())

			for _, ns := range namespaces.List() {
				if ns == source.Name {
					continue
				}
				_, err := fi.KubeClient.CoreV1().Secrets(ns).Get(context.TODO(), source.Name, metav1.GetOptions{})
				if err != nil {
					return false
				}
			}
			return true
		}
		return false
	})
}

func (fi *Invocation) EventuallySecretNotSynced(source *core.Secret) GomegaAsyncAssertion {
	return Eventually(func() bool {
		namespaces, err := fi.KubeClient.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		Expect(err).NotTo(HaveOccurred())

		for _, ns := range namespaces.Items {
			if ns.Name == source.Name {
				continue
			}
			_, err := fi.KubeClient.CoreV1().Secrets(ns.Namespace).Get(context.TODO(), source.Name, metav1.GetOptions{})
			if err == nil {
				return false
			}
		}
		return true
	})
}

func (fi *Invocation) EventuallySecretSyncedToNamespace(source *core.Secret, namespace string) GomegaAsyncAssertion {
	return Eventually(func() bool {
		_, err := fi.KubeClient.CoreV1().Secrets(namespace).Get(context.TODO(), source.Name, metav1.GetOptions{})
		return err == nil
	})
}

func (fi *Invocation) EventuallySyncedSecretsUpdated(source *core.Secret) GomegaAsyncAssertion {
	opt := syncer.GetSyncOptions(source.Annotations)

	return Eventually(func() bool {
		if opt.NamespaceSelector != nil {
			namespaces, err := syncer.NamespacesForSelector(fi.KubeClient, *opt.NamespaceSelector)
			Expect(err).NotTo(HaveOccurred())

			for _, ns := range namespaces.List() {
				if ns == source.Namespace {
					continue
				}
				secretReplica, err := fi.KubeClient.CoreV1().Secrets(ns).Get(context.TODO(), source.Name, metav1.GetOptions{})
				if err != nil {
					return false
				}
				if !reflect.DeepEqual(source.Data, secretReplica.Data) {
					return false
				}
			}
			return true
		}
		return false
	})
}

func (fi *Invocation) EventuallySyncedSecretsDeleted(source *core.Secret) GomegaAsyncAssertion {
	opt := syncer.GetSyncOptions(source.Annotations)

	return Eventually(func() bool {
		if opt.NamespaceSelector != nil {
			namespaces, err := syncer.NamespacesForSelector(fi.KubeClient, *opt.NamespaceSelector)
			Expect(err).NotTo(HaveOccurred())

			for _, ns := range namespaces.List() {
				if ns == source.Namespace {
					continue
				}
				_, err := fi.KubeClient.CoreV1().Secrets(ns).Get(context.TODO(), source.Name, metav1.GetOptions{})
				if err == nil {
					return false
				}
			}
			return true
		}
		return false
	})
}

func (fi *Invocation) DeleteAllSecrets() {
	secrets, err := fi.KubeClient.CoreV1().Secrets(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labels.Set{
			"app": fi.App(),
		}.String(),
	})
	Expect(err).NotTo(HaveOccurred())

	for _, value := range secrets.Items {
		err := fi.KubeClient.CoreV1().Secrets(value.Namespace).Delete(context.TODO(), value.Name, metav1.DeleteOptions{})
		if kerr.IsNotFound(err) {
			err = nil
		}
		Expect(err).NotTo(HaveOccurred())
	}
}

func (f *Framework) CreateSecret(obj *core.Secret) (*core.Secret, error) {
	return f.KubeClient.CoreV1().Secrets(obj.Namespace).Create(context.TODO(), obj, metav1.CreateOptions{})
}

func (f *Framework) DeleteSecret(meta metav1.ObjectMeta) error {
	return f.KubeClient.CoreV1().Secrets(meta.Namespace).Delete(context.TODO(), meta.Name, *deleteInForeground())
}

func (fi *Invocation) SecretForWebhookNotifier() *core.Secret {
	namespace := OperatorNamespace
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

func (fi *Invocation) WaitUntilSecretCreated(meta metav1.ObjectMeta) error {
	return wait.PollImmediate(kutil.RetryInterval, kutil.ReadinessTimeout, func() (done bool, err error) {
		if _, err := fi.KubeClient.CoreV1().Secrets(meta.Namespace).Get(context.TODO(), meta.Name, metav1.GetOptions{}); err != nil {
			if kerr.IsNotFound(err) {
				return false, nil
			} else {
				return true, err
			}
		}
		return true, nil
	})
}

func (fi *Invocation) WaitUntilSecretDeleted(meta metav1.ObjectMeta) error {
	return wait.PollImmediate(kutil.RetryInterval, kutil.GCTimeout, func() (done bool, err error) {
		if _, err := fi.KubeClient.CoreV1().Secrets(meta.Namespace).Get(context.TODO(), meta.Name, metav1.GetOptions{}); err != nil {
			if kerr.IsNotFound(err) {
				return true, nil
			} else {
				return true, err
			}
		}
		return false, nil
	})
}
