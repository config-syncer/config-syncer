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

package e2e_test

import (
	"context"
	"os"

	"kubeops.dev/config-syncer/pkg/operator"
	"kubeops.dev/config-syncer/pkg/syncer"
	"kubeops.dev/config-syncer/test/e2e/framework"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	core_util "kmodules.xyz/client-go/core/v1"
	"kmodules.xyz/client-go/meta"
)

var _ = Describe("Secret-Syncer", func() {
	var (
		f           *framework.Invocation
		secret      *core.Secret
		nsWithLabel *core.Namespace
		config      operator.Config
	)

	BeforeEach(func() {
		f = root.Invoke()
		secret = f.NewSecret()
		nsWithLabel = f.NewNamespaceWithLabel()
	})

	AfterEach(func() {
		f.DeleteAllSecrets()

		err := f.KubeClient.CoreV1().Namespaces().Delete(context.TODO(), nsWithLabel.Name, metav1.DeleteOptions{})
		if kerr.IsNotFound(err) {
			err = nil
		}
		Expect(err).NotTo(HaveOccurred())
		f.EventuallyNamespaceDeleted(nsWithLabel.Name).Should(BeTrue())
	})

	shouldSyncSecretToAllNamespaces := func() {
		By("Creating secret")
		sourceSecret, err := f.KubeClient.CoreV1().Secrets(secret.Namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
		Expect(err).NotTo(HaveOccurred())

		By("Checking secret has not synced yet")
		f.EventuallySecretNotSynced(sourceSecret).Should(BeTrue())

		By("Adding sync annotation")
		sourceSecret, _, err = core_util.PatchSecret(context.TODO(), f.KubeClient, sourceSecret, func(obj *core.Secret) *core.Secret {
			metav1.SetMetaDataAnnotation(&obj.ObjectMeta, syncer.ConfigSyncKey, "")
			return obj
		}, metav1.PatchOptions{})
		Expect(err).ShouldNot(HaveOccurred())

		By("Checking secret has synced to all namespaces")
		f.EventuallySecretSynced(sourceSecret).Should(BeTrue())
	}

	Describe("Across Namespaces", func() {
		BeforeEach(func() {
			config = operator.Config{}
		})

		Context("All Namespaces", func() {
			It("should sync secret to all namespaces", shouldSyncSecretToAllNamespaces)
		})

		Context("New Namespace", func() {
			It("should synced secret to new namespace", func() {
				shouldSyncSecretToAllNamespaces()

				By("Creating new namespace")
				err := f.CreateNamespace(nsWithLabel)
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking new namespace has the secret")
				f.EventuallySecretSyncedToNamespace(secret, nsWithLabel.Name).Should(BeTrue())
			})
		})

		Context("Remove Sync Annotation", func() {
			It("should delete synced secrets", func() {
				shouldSyncSecretToAllNamespaces()

				By("Removing sync annotation")
				source, err := f.KubeClient.CoreV1().Secrets(secret.Namespace).Get(context.TODO(), secret.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())

				_, _, err = core_util.PatchSecret(context.TODO(), f.KubeClient, source, func(obj *core.Secret) *core.Secret {
					obj.Annotations = meta.RemoveKey(obj.Annotations, syncer.ConfigSyncKey)
					return obj
				}, metav1.PatchOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking synced secrets has been deleted")
				f.EventuallySyncedSecretsDeleted(source)
			})
		})

		Context("Source Update", func() {
			It("should update synced secrets", func() {
				shouldSyncSecretToAllNamespaces()

				By("Updating source secret")
				source, err := f.KubeClient.CoreV1().Secrets(secret.Namespace).Get(context.TODO(), secret.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())

				source, _, err = core_util.PatchSecret(context.TODO(), f.KubeClient, source, func(obj *core.Secret) *core.Secret {
					if obj.Data == nil {
						obj.Data = map[string][]byte{}
					}
					obj.Data["data"] = []byte("test")
					return obj
				}, metav1.PatchOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking synced secrets has been updated")
				f.EventuallySyncedSecretsUpdated(source).Should(BeTrue())
			})
		})

		Context("Backward Compatibility", func() {
			It("should sync secret to all namespaces", func() {
				By("Creating secret")
				source, err := f.CreateSecret(secret)
				Expect(err).NotTo(HaveOccurred())

				By("Checking secret has not synced yet")
				f.EventuallySecretNotSynced(source).Should(BeTrue())

				By("Adding sync=true annotation")
				source, _, err = core_util.PatchSecret(context.TODO(), f.KubeClient, source, func(obj *core.Secret) *core.Secret {
					metav1.SetMetaDataAnnotation(&obj.ObjectMeta, syncer.ConfigSyncKey, "true")
					return obj
				}, metav1.PatchOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking secret has synced to all namespaces")
				f.EventuallySecretSynced(source).Should(BeTrue())
			})
		})

		Context("Namespace Selector", func() {
			It("should add secret to selected namespaces", func() {
				shouldSyncSecretToAllNamespaces()

				By("Adding selector annotation")
				source, err := f.KubeClient.CoreV1().Secrets(secret.Namespace).Get(context.TODO(), secret.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())

				source, _, err = core_util.PatchSecret(context.TODO(), f.KubeClient, source, func(obj *core.Secret) *core.Secret {
					metav1.SetMetaDataAnnotation(&obj.ObjectMeta, syncer.ConfigSyncKey, "app="+f.App())
					return obj
				}, metav1.PatchOptions{})
				Expect(err).NotTo(HaveOccurred())

				By("Checking secret has not synced to other namespaces")
				f.EventuallySecretNotSynced(source).Should(BeTrue())

				By("Creating new namespace with label")
				err = f.CreateNamespace(nsWithLabel)
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking secret synced to new namespace")
				f.EventuallySecretSyncedToNamespace(source, nsWithLabel.Name)

				By("Changing selector annotation")
				_, _, err = core_util.PatchSecret(context.TODO(), f.KubeClient, source, func(obj *core.Secret) *core.Secret {
					metav1.SetMetaDataAnnotation(&obj.ObjectMeta, syncer.ConfigSyncKey, "app=do-not-match")
					return obj
				}, metav1.PatchOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking synced secrets has been deleted")
				f.EventuallySyncedSecretsDeleted(source)

				By("Removing selector annotation")
				source, err = f.KubeClient.CoreV1().Secrets(source.Namespace).Get(context.TODO(), source.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())

				source, _, err = core_util.PatchSecret(context.TODO(), f.KubeClient, source, func(obj *core.Secret) *core.Secret {
					metav1.SetMetaDataAnnotation(&obj.ObjectMeta, syncer.ConfigSyncKey, "")
					return obj
				}, metav1.PatchOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking secret synced to all namespaces")
				f.EventuallySecretSynced(source).Should(BeTrue())
			})
		})

		Context("Source Deleted", func() {
			It("should delete synced secrets", func() {
				shouldSyncSecretToAllNamespaces()

				By("Deleting source secret")
				source, err := f.KubeClient.CoreV1().Secrets(secret.Namespace).Get(context.TODO(), secret.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())

				err = f.DeleteSecret(source.ObjectMeta)
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking synced secrets has been deleted")
				f.EventuallySyncedSecretsDeleted(source).Should(BeTrue())
			})
		})

		Context("Source Namespace Deleted", func() {
			var sourceNamespace *core.Namespace

			BeforeEach(func() {
				sourceNamespace = f.NewNamespace("source")
				secret.Namespace = sourceNamespace.Name
			})

			It("should delete synced secrets", func() {
				By("Creating source namespace")
				err := f.CreateNamespace(sourceNamespace)
				Expect(err).NotTo(HaveOccurred())

				shouldSyncSecretToAllNamespaces()

				By("Deleting source namespace")
				source, err := f.KubeClient.CoreV1().Secrets(secret.Namespace).Get(context.TODO(), secret.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())

				err = f.DeleteNamespace(sourceNamespace.Name)
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking synced secrets has been deleted")
				f.EventuallySyncedSecretsDeleted(source).Should(BeTrue())
			})
		})
	})

	Describe("Across Cluster", func() {
		Context("Secret Context Syncer Test", func() {
			var (
				kubeConfigPath = "/home/dipta/all/config-syncer-test/kubeconfig"
				ctx            = "gke_tigerworks-kube_us-central1-f_kite"
			)

			BeforeEach(func() {
				config = operator.Config{}
				config.ClusterName = "minikube"
				config.KubeConfigFile = kubeConfigPath

				if _, err := os.Stat(kubeConfigPath); err != nil {
					Skip(`"config" file not found on` + kubeConfigPath)
				}

				By("Creating namespace for context")
				f.EnsureNamespaceForContext(kubeConfigPath, ctx)
			})

			AfterEach(func() {
				By("Deleting namespaces for contexts")
				f.DeleteNamespaceForContext(kubeConfigPath, ctx)
			})

			XIt("Should add secret to contexts", func() {
				By("Creating source ns in remote cluster")
				f.EnsureNamespaceForContext(kubeConfigPath, ctx)

				By("Creating secret")
				secret, err := f.KubeClient.CoreV1().Secrets(secret.Namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
				Expect(err).NotTo(HaveOccurred())

				By("Adding sync annotation")
				secret, _, err = core_util.PatchSecret(context.TODO(), f.KubeClient, secret, func(obj *core.Secret) *core.Secret {
					metav1.SetMetaDataAnnotation(&obj.ObjectMeta, syncer.ConfigSyncContexts, ctx)
					return obj
				}, metav1.PatchOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking secret added to contexts")
				f.EventuallyNumOfConfigmapsForContext(kubeConfigPath, ctx).Should(BeNumerically("==", 1))

				By("Removing sync annotation")
				_, _, err = core_util.PatchSecret(context.TODO(), f.KubeClient, secret, func(obj *core.Secret) *core.Secret {
					obj.Annotations = meta.RemoveKey(obj.Annotations, syncer.ConfigSyncContexts)
					return obj
				}, metav1.PatchOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking secret removed from contexts")
				f.EventuallyNumOfConfigmapsForContext(kubeConfigPath, ctx).Should(BeNumerically("==", 0))
			})
		})
	})
})
