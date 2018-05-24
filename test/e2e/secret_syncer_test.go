package e2e_test

import (
	"os"
	"time"

	api "github.com/appscode/kubed/apis/kubed/v1alpha1"
	"github.com/appscode/kubed/test/e2e/framework"
	core_util "github.com/appscode/kutil/core/v1"
	"github.com/appscode/kutil/meta"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Secret-Syncer", func() {
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
		if f.SelfHostedOperator {
			By("Restarting kubed operator")
			err:=f.RestartKubedOperator(&clusterConfig)
			Expect(err).NotTo(HaveOccurred())
		} else {
			By("Starting Kubed")
			stopCh = make(chan struct{})
			err := f.RunKubed(stopCh, clusterConfig)
			Expect(err).NotTo(HaveOccurred())

			By("Waiting for API server to be ready")
			root.EventuallyAPIServerReady().Should(Succeed())
			time.Sleep(time.Second * 5)
		}
	})

	AfterEach(func() {
		if !f.SelfHostedOperator {
			close(stopCh)
		}
		f.DeleteAllSecrets()

		err := f.KubeClient.CoreV1().Namespaces().Delete(nsWithLabel.Name, &metav1.DeleteOptions{})
		if kerr.IsNotFound(err) {
			err = nil
		}
		Expect(err).NotTo(HaveOccurred())
		f.EventuallyNamespaceDeleted(nsWithLabel.Name).Should(BeTrue())
	})

	var (
		shouldSyncSecretToAllNamespaces = func() {
			By("Creating secret")
			sourceSecret, err := f.KubeClient.CoreV1().Secrets(secret.Namespace).Create(secret)
			Expect(err).NotTo(HaveOccurred())

			By("Checking secret has not synced yet")
			f.EventuallySecretNotSynced(sourceSecret).Should(BeTrue())

			By("Adding sync annotation")
			sourceSecret, _, err = core_util.PatchSecret(f.KubeClient, sourceSecret, func(obj *core.Secret) *core.Secret {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncKey, "")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())

			By("Checking secret has synced to all namespaces")
			f.EventuallySecretSynced(sourceSecret).Should(BeTrue())
		}
	)

	Describe("Across Namespaces", func() {

		BeforeEach(func() {
			clusterConfig = framework.ConfigSyncClusterConfig()
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
				source, err := f.KubeClient.CoreV1().Secrets(secret.Namespace).Get(secret.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())

				_, _, err = core_util.PatchSecret(f.KubeClient, source, func(obj *core.Secret) *core.Secret {
					obj.Annotations = meta.RemoveKey(obj.Annotations, api.ConfigSyncKey)
					return obj
				})
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking synced secrets has been deleted")
				f.EventuallySyncedSecretsDeleted(source)
			})
		})

		Context("Source Update", func() {

			It("should update synced secrets", func() {
				shouldSyncSecretToAllNamespaces()

				By("Updating source secret")
				source, err := f.KubeClient.CoreV1().Secrets(secret.Namespace).Get(secret.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())

				source, _, err = core_util.PatchSecret(f.KubeClient, source, func(obj *core.Secret) *core.Secret {
					obj.Data["data"] = []byte("test")
					return obj
				})
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
				source, _, err = core_util.PatchSecret(f.KubeClient, source, func(obj *core.Secret) *core.Secret {
					metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncKey, "true")
					return obj
				})
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking secret has synced to all namespaces")
				f.EventuallySecretSynced(source).Should(BeTrue())
			})
		})

		Context("Namespace Selector", func() {

			It("should add secret to selected namespaces", func() {

				shouldSyncSecretToAllNamespaces()

				By("Adding selector annotation")
				source, err := f.KubeClient.CoreV1().Secrets(secret.Namespace).Get(secret.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())

				source, _, err = core_util.PatchSecret(f.KubeClient, source, func(obj *core.Secret) *core.Secret {
					metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncKey, "app="+f.App())
					return obj
				})
				Expect(err).NotTo(HaveOccurred())

				By("Checking secret has not synced to other namespaces")
				f.EventuallySecretNotSynced(source).Should(BeTrue())

				By("Creating new namespace with label")
				err = f.CreateNamespace(nsWithLabel)
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking secret synced to new namespace")
				f.EventuallySecretSyncedToNamespace(source, nsWithLabel.Name)

				By("Changing selector annotation")
				_, _, err = core_util.PatchSecret(f.KubeClient, source, func(obj *core.Secret) *core.Secret {
					metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncKey, "app=do-not-match")
					return obj
				})
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking synced secrets has been deleted")
				f.EventuallySyncedSecretsDeleted(source)

				By("Removing selector annotation")
				source, err = f.KubeClient.CoreV1().Secrets(source.Namespace).Get(source.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())

				source, _, err = core_util.PatchSecret(f.KubeClient, source, func(obj *core.Secret) *core.Secret {
					metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncKey, "")
					return obj
				})
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking secret synced to all namespaces")
				f.EventuallySecretSynced(source).Should(BeTrue())
			})
		})

		Context("Source Deleted", func() {

			It("should delete synced secrets", func() {
				shouldSyncSecretToAllNamespaces()

				By("Deleting source secret")
				source, err := f.KubeClient.CoreV1().Secrets(secret.Namespace).Get(secret.Name, metav1.GetOptions{})
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
				source, err := f.KubeClient.CoreV1().Secrets(secret.Namespace).Get(secret.Name, metav1.GetOptions{})
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
				kubeConfigPath = "/home/dipta/all/kubed-test/kubeconfig"
				context        = "gke_tigerworks-kube_us-central1-f_kite"
			)

			BeforeEach(func() {
				clusterConfig = framework.ConfigSyncClusterConfig()
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
				secret, err := f.KubeClient.CoreV1().Secrets(secret.Namespace).Create(secret)
				Expect(err).NotTo(HaveOccurred())

				By("Adding sync annotation")
				secret, _, err = core_util.PatchSecret(f.KubeClient, secret, func(obj *core.Secret) *core.Secret {
					metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncContexts, context)
					return obj
				})
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking secret added to contexts")
				f.EventuallyNumOfConfigmapsForContext(kubeConfigPath, context).Should(BeNumerically("==", 1))

				By("Removing sync annotation")
				secret, _, err = core_util.PatchSecret(f.KubeClient, secret, func(obj *core.Secret) *core.Secret {
					obj.Annotations = meta.RemoveKey(obj.Annotations, api.ConfigSyncContexts)
					return obj
				})
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking secret removed from contexts")
				f.EventuallyNumOfConfigmapsForContext(kubeConfigPath, context).Should(BeNumerically("==", 0))
			})
		})
	})
})
