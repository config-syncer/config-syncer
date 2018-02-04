package e2e

import (
	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/kubed/pkg/api"
	"github.com/appscode/kubed/test/framework"
	core_util "github.com/appscode/kutil/core/v1"
	"github.com/appscode/kutil/meta"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Secret-syncer", func() {
	var (
		f               *framework.Invocation
		secret          *core.Secret
		nsWithLabel     *core.Namespace
		numOfNamespaces = func() int {
			ns, err := f.KubeClient.CoreV1().Namespaces().List(metav1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())
			return len(ns.Items)
		}
	)

	BeforeEach(func() {
		f = root.Invoke()

		secret = &core.Secret{
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

		nsWithLabel = &core.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: rand.WithUniqSuffix("kubed-e2e-labeled"),
				Labels: map[string]string{
					"app": f.App(),
				},
			},
		}
	})

	AfterEach(func() {
		f.DeleteAllSecrets()

		err := f.KubeClient.CoreV1().Namespaces().Delete(nsWithLabel.Name, &metav1.DeleteOptions{})
		if kerr.IsNotFound(err) {
			err = nil
		}
		Expect(err).NotTo(HaveOccurred())
		f.EventuallyNamespaceDeleted(nsWithLabel.Name).Should(BeTrue())
	})

	Describe("Secret Syncer Test", func() {
		It("Should add secret to all namespaces", func() {
			By("Creating secret")
			c, err := root.KubeClient.CoreV1().Secrets(secret.Namespace).Create(secret)
			Expect(err).NotTo(HaveOccurred())
			f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", 1))

			By("Adding sync annotation")
			c, _, err = core_util.PatchSecret(f.KubeClient, c, func(obj *core.Secret) *core.Secret {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncKey, "true")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", numOfNamespaces()))

			By("Creating new namespace")
			_, err = root.KubeClient.CoreV1().Namespaces().Create(nsWithLabel)
			Expect(err).ShouldNot(HaveOccurred())
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", numOfNamespaces()))

			By("Removing sync annotation")
			c, _, err = core_util.PatchSecret(f.KubeClient, c, func(obj *core.Secret) *core.Secret {
				obj.Annotations = meta.RemoveKey(obj.Annotations, api.ConfigSyncKey)
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())
			f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", 1))
		})
	})

	Describe("Secret Syncer Test Backward Compatibility Tes", func() {
		It("Should add secret to all namespaces", func() {
			By("Creating secret")
			c, err := root.KubeClient.CoreV1().Secrets(secret.Namespace).Create(secret)
			Expect(err).NotTo(HaveOccurred())
			f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", 1))

			By("Adding sync annotation")
			c, _, err = core_util.PatchSecret(f.KubeClient, c, func(obj *core.Secret) *core.Secret {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncKey, "true")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", numOfNamespaces()))
		})
	})

	Describe("Secret Syncer With Namespace Selector", func() {
		It("Should add secret to selected namespaces", func() {
			By("Creating secret")
			c, err := root.KubeClient.CoreV1().Secrets(secret.Namespace).Create(secret)
			Expect(err).NotTo(HaveOccurred())
			f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", 1))

			By("Adding sync annotation")
			c, _, err = core_util.PatchSecret(f.KubeClient, c, func(obj *core.Secret) *core.Secret {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncKey, "true")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", numOfNamespaces()))

			By("Adding selector annotation")
			c, _, err = core_util.PatchSecret(f.KubeClient, c, func(obj *core.Secret) *core.Secret {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncKey, "app="+f.App())
				return obj
			})
			Expect(err).NotTo(HaveOccurred())
			f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", 1))

			By("Creating new namespace with label")
			_, err = root.KubeClient.CoreV1().Namespaces().Create(nsWithLabel)
			Expect(err).ShouldNot(HaveOccurred())
			f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfSecrets(nsWithLabel.Name).Should(BeNumerically("==", 1))
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", 2))

			By("Changing selector annotation")
			c, _, err = core_util.PatchSecret(f.KubeClient, c, func(obj *core.Secret) *core.Secret {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncKey, "app=do-not-match")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())
			f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", 1))

			By("Removing selector annotation")
			c, _, err = core_util.PatchSecret(f.KubeClient, c, func(obj *core.Secret) *core.Secret {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncKey, "")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())
			f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", numOfNamespaces()))
		})
	})

	Describe("Secret Syncer Source Deleted", func() {
		It("Should delete synced secrets from namespaces", func() {
			By("Creating secret")
			c, err := root.KubeClient.CoreV1().Secrets(secret.Namespace).Create(secret)
			Expect(err).NotTo(HaveOccurred())
			f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", 1))

			By("Adding sync annotation")
			c, _, err = core_util.PatchSecret(f.KubeClient, c, func(obj *core.Secret) *core.Secret {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncKey, "")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", numOfNamespaces()))

			By("Creating new namespace")
			_, err = root.KubeClient.CoreV1().Namespaces().Create(nsWithLabel)
			Expect(err).ShouldNot(HaveOccurred())
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", numOfNamespaces()))

			By("Deleting source secret")
			err = f.KubeClient.CoreV1().Secrets(secret.Namespace).Delete(secret.Name, &metav1.DeleteOptions{})
			Expect(err).ShouldNot(HaveOccurred())
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", 0))
		})
	})

	Describe("Secret Context Syncer Test", func() {
		var (
			kubeConfigPath = "/home/dipta/all/kubed-test/kubeconfig"
			context        = "gke_tigerworks-kube_us-central1-f_kite"
		)

		BeforeEach(func() {
			By("Creating namespace for context")
			f.EnsureNamespaceForContext(kubeConfigPath, context)
		})

		AfterEach(func() {
			By("Deleting namespaces for contexts")
			f.DeleteNamespaceForContext(kubeConfigPath, context)
		})

		It("Should add secret to contexts", func() {
			By("Creating source ns in remote cluster")
			f.EnsureNamespaceForContext(kubeConfigPath, context)

			By("Creating secret")
			secret, err := root.KubeClient.CoreV1().Secrets(secret.Namespace).Create(secret)
			Expect(err).NotTo(HaveOccurred())

			By("Adding sync annotation")
			secret, _, err = core_util.PatchSecret(f.KubeClient, secret, func(obj *core.Secret) *core.Secret {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncContexts, context)
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())

			By("Checking secret added to contexts")
			f.EventuallyNumOfSecretsForContext(kubeConfigPath, context).Should(BeNumerically("==", 1))

			By("Removing sync annotation")
			secret, _, err = core_util.PatchSecret(f.KubeClient, secret, func(obj *core.Secret) *core.Secret {
				obj.Annotations = meta.RemoveKey(obj.Annotations, api.ConfigSyncContexts)
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())

			By("Checking secret removed from contexts")
			f.EventuallyNumOfSecretsForContext(kubeConfigPath, context).Should(BeNumerically("==", 0))
		})
	})
})
