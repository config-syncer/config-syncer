package e2e

import (
	"time"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/kubed/pkg/config"
	"github.com/appscode/kubed/test/framework"
	core_util "github.com/appscode/kutil/core/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
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
				Name: rand.WithUniqSuffix("kubed-test-label"),
				Labels: map[string]string{
					"app": f.App(),
				},
			},
		}
	})

	AfterEach(func() {
		secrets, err := f.KubeClient.CoreV1().Secrets(metav1.NamespaceAll).List(metav1.ListOptions{
			LabelSelector: labels.Set{
				"app": f.App(),
			}.String(),
		})
		Expect(err).NotTo(HaveOccurred())

		for _, value := range secrets.Items {
			err := f.KubeClient.CoreV1().Secrets(value.Namespace).Delete(value.Name, &metav1.DeleteOptions{})
			Expect(err).NotTo(HaveOccurred())
		}

		err = f.KubeClient.CoreV1().Namespaces().Delete(nsWithLabel.Name, &metav1.DeleteOptions{})
		time.Sleep(time.Minute) // time to delete namespace. Is there any WaitUntilNamespaceDeleted?
		// Expect(err).NotTo(HaveOccurred())
	})

	Describe("Secret Syncer Test", func() {
		FIt("Should add secret to all namespaces", func() {
			By("Creating secret")
			c, err := root.KubeClient.CoreV1().Secrets(secret.Namespace).Create(secret)
			Expect(err).NotTo(HaveOccurred())
			f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", 1))

			By("Adding sync annotation")
			c, err = core_util.PatchSecret(f.KubeClient, c, func(obj *core.Secret) *core.Secret {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, config.ConfigSyncKey, "true")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", numOfNamespaces()))

			By("Creating new namespace")
			_, err = root.KubeClient.CoreV1().Namespaces().Create(nsWithLabel)
			Expect(err).ShouldNot(HaveOccurred())
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", numOfNamespaces()))

			By("Removing sync annotation")
			c, err = core_util.PatchSecret(f.KubeClient, c, func(obj *core.Secret) *core.Secret {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, config.ConfigSyncKey, "false")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())
			f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", 1))
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
			c, err = core_util.PatchSecret(f.KubeClient, c, func(obj *core.Secret) *core.Secret {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, config.ConfigSyncKey, "true")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", numOfNamespaces()))

			By("Adding selector annotation")
			c, err = core_util.PatchSecret(f.KubeClient, c, func(obj *core.Secret) *core.Secret {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, config.ConfigSyncNsSelector, "app="+f.App())
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
			c, err = core_util.PatchSecret(f.KubeClient, c, func(obj *core.Secret) *core.Secret {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, config.ConfigSyncNsSelector, "app=do-not-match")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())
			f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", 1))

			By("Removing selector annotation")
			c, err = core_util.PatchSecret(f.KubeClient, c, func(obj *core.Secret) *core.Secret {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, config.ConfigSyncNsSelector, "")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())
			f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", numOfNamespaces()))
		})
	})
})
