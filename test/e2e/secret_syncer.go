package e2e

import (
	"os"

	api "github.com/appscode/kubed/apis/kubed/v1alpha1"
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
		f             *framework.Invocation
		secret        *core.Secret
		nsWithLabel   *core.Namespace
		stopCh        chan struct{}
		clusterConfig api.ClusterConfig
	)

	BeforeEach(func() {
		f = root.Invoke()
		secret = f.NewSecret()
		nsWithLabel = f.NewNamespaceWithLabel()
	})

	JustBeforeEach(func() {
		By("Starting Operator")
		stopCh = make(chan struct{})
		err := f.RunOperator(stopCh, clusterConfig)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		close(stopCh)
		f.DeleteAllSecrets()

		err := f.KubeClient.CoreV1().Namespaces().Delete(nsWithLabel.Name, &metav1.DeleteOptions{})
		if kerr.IsNotFound(err) {
			err = nil
		}
		Expect(err).NotTo(HaveOccurred())
		f.EventuallyNamespaceDeleted(nsWithLabel.Name).Should(BeTrue())

		err = framework.ResetTestConfigFile()
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Secret Syncer Test", func() {

		BeforeEach(func() {
			clusterConfig = framework.ConfigMapSyncClusterConfig()
		})

		It("Should add secret to all namespaces", func() {

			By("Creating secret")
			s, err := root.KubeClient.CoreV1().Secrets(secret.Namespace).Create(secret)
			Expect(err).NotTo(HaveOccurred())

			By("Checking secret has not synced yet")
			f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", 1))

			By("Adding sync annotation")
			s, _, err = core_util.PatchSecret(f.KubeClient, s, func(obj *core.Secret) *core.Secret {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncKey, "true")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())

			By("Checking secret has synced to all namespaces")
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", f.NumberOfNameSpace()))

			By("Creating new namespace")
			_, err = root.KubeClient.CoreV1().Namespaces().Create(nsWithLabel)
			Expect(err).ShouldNot(HaveOccurred())

			By("Checking secret has been added to new namespace")
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", f.NumberOfNameSpace()))
			_, err = f.KubeClient.CoreV1().Secrets(secret.Namespace).Get(secret.Name, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())

			By("Removing sync annotation")
			s, err = f.KubeClient.CoreV1().Secrets(secret.Namespace).Get(secret.Name, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			s, _, err = core_util.PatchSecret(f.KubeClient, s, func(obj *core.Secret) *core.Secret {
				obj.Annotations = meta.RemoveKey(obj.Annotations, api.ConfigSyncKey)
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())

			By("Checking secret has removed from other namespaces")
			f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", 1))
		})
	})

	Describe("Secret Syncer Backward Compatibility Test", func() {

		BeforeEach(func() {
			clusterConfig = framework.ConfigMapSyncClusterConfig()
		})

		It("Should add secret to all namespaces", func() {

			By("Creating secret")
			s, err := root.KubeClient.CoreV1().Secrets(secret.Namespace).Create(secret)
			Expect(err).NotTo(HaveOccurred())

			By("Checking secret has not synced yet")
			f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", 1))

			By("Adding sync annotation")
			s, _, err = core_util.PatchSecret(f.KubeClient, s, func(obj *core.Secret) *core.Secret {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncKey, "true")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())

			By("Checking secret has synced to all namespaces")
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", f.NumberOfNameSpace()))
		})
	})

	Describe("Secret Syncer With Namespace Selector", func() {

		BeforeEach(func() {
			clusterConfig = framework.ConfigMapSyncClusterConfig()
		})

		It("Should add secret to selected namespaces", func() {

			By("Creating secret")
			s, err := root.KubeClient.CoreV1().Secrets(secret.Namespace).Create(secret)
			Expect(err).NotTo(HaveOccurred())

			By("Checking secret has not synced yet")
			f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", 1))

			By("Adding sync annotation")
			s, _, err = core_util.PatchSecret(f.KubeClient, s, func(obj *core.Secret) *core.Secret {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncKey, "true")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())

			By("Checking secret has synced to all namespaces")
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", f.NumberOfNameSpace()))

			By("Adding selector annotation")
			s, _, err = core_util.PatchSecret(f.KubeClient, s, func(obj *core.Secret) *core.Secret {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncKey, "app="+f.App())
				return obj
			})
			Expect(err).NotTo(HaveOccurred())

			By("Checking secret has removed from no-matching namespaces")
			f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", 1))

			By("Creating new namespace with matching label")
			_, err = root.KubeClient.CoreV1().Namespaces().Create(nsWithLabel)
			Expect(err).ShouldNot(HaveOccurred())

			By("Checking secret has added to new namespace")
			f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfSecrets(nsWithLabel.Name).Should(BeNumerically("==", 1))
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", 2))

			By("Changing selector annotation")
			s, _, err = core_util.PatchSecret(f.KubeClient, s, func(obj *core.Secret) *core.Secret {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncKey, "app=do-not-match")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())

			By("Checking secret has removed from no-matching namespaces")
			f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", 1))

			By("Removing selector annotation")
			s, _, err = core_util.PatchSecret(f.KubeClient, s, func(obj *core.Secret) *core.Secret {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncKey, "")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())

			By("Checking secret has synced to all namespaces")
			f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", f.NumberOfNameSpace()))
		})
	})

	Describe("Secret Syncer Source Deleted", func() {

		BeforeEach(func() {
			clusterConfig = framework.ConfigMapSyncClusterConfig()
		})

		It("Should delete synced secrets from namespaces", func() {

			By("Creating secret")
			s, err := root.KubeClient.CoreV1().Secrets(secret.Namespace).Create(secret)
			Expect(err).NotTo(HaveOccurred())

			By("Checkin secret has not synced yet")
			f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", 1))

			By("Adding sync annotation")
			s, _, err = core_util.PatchSecret(f.KubeClient, s, func(obj *core.Secret) *core.Secret {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncKey, "")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())

			By("Checking secret has synced to all namespaces")
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", f.NumberOfNameSpace()))

			By("Creating new namespace")
			_, err = root.KubeClient.CoreV1().Namespaces().Create(nsWithLabel)
			Expect(err).ShouldNot(HaveOccurred())

			By("Checking secret has synced to new namespace")
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", f.NumberOfNameSpace()))

			By("Deleting source secret")
			err = f.KubeClient.CoreV1().Secrets(secret.Namespace).Delete(secret.Name, &metav1.DeleteOptions{})
			Expect(err).ShouldNot(HaveOccurred())

			By("Checking secret has removed from all namespaces")
			f.EventuallyNumOfSecrets(metav1.NamespaceAll).Should(BeNumerically("==", 0))
		})
	})

	Describe("Secret Context Syncer Test", func() {
		var (
			kubeConfigPath = "/home/dipta/all/kubed-test/kubeconfig"
			context        = "gke_tigerworks-kube_us-central1-f_kite"
		)

		BeforeEach(func() {
			clusterConfig = framework.ConfigMapSyncClusterConfig()
			clusterConfig.ClusterName = "minikube"
			clusterConfig.KubeConfigFile = kubeConfigPath

			if _, err := os.Stat(kubeConfigPath); err != nil {
				Skip(`"config" file not found on` + kubeConfigPath)
			}

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
