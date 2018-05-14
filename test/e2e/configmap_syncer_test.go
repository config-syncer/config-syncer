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

var _ = Describe("Config-syncer", func() {
	var (
		f             *framework.Invocation
		cfgMap        *core.ConfigMap
		nsWithLabel   *core.Namespace
		stopCh        chan struct{}
		clusterConfig api.ClusterConfig
	)

	BeforeEach(func() {
		f = root.Invoke()
		cfgMap = f.NewConfigMap()
		nsWithLabel = f.NewNamespaceWithLabel()
	})

	JustBeforeEach(func() {
		By("Starting Kubed")
		stopCh = make(chan struct{})
		err := f.RunKubed(stopCh, clusterConfig)
		Expect(err).NotTo(HaveOccurred())

		By("Waiting for API server to be ready")
		root.EventuallyAPIServerReady().Should(Succeed())
		time.Sleep(time.Second * 5)
	})

	AfterEach(func() {
		close(stopCh)
		f.DeleteAllConfigmaps()

		err := f.KubeClient.CoreV1().Namespaces().Delete(nsWithLabel.Name, &metav1.DeleteOptions{})
		if kerr.IsNotFound(err) {
			err = nil
		}
		Expect(err).NotTo(HaveOccurred())
		f.EventuallyNamespaceDeleted(nsWithLabel.Name).Should(BeTrue())

		err = framework.ResetTestConfigFile()
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("ConfigMap Syncer Test", func() {
		BeforeEach(func() {
			clusterConfig = framework.ConfigMapSyncClusterConfig()
		})

		It("Should add configmap to all namespaces", func() {

			By("Creating configmap")
			c, err := f.KubeClient.CoreV1().ConfigMaps(cfgMap.Namespace).Create(cfgMap)
			Expect(err).NotTo(HaveOccurred())

			By("Checking configmap has not sync. yet")
			f.EventuallyNumOfConfigmaps(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", 1))

			By("Adding sync annotation")
			c, _, err = core_util.PatchConfigMap(f.KubeClient, c, func(obj *core.ConfigMap) *core.ConfigMap {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncKey, "")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())

			By("Checking configmap has synced to all namespaces")
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", f.NumberOfNameSpace()))

			By("Creating new namespace")
			_, err = f.KubeClient.CoreV1().Namespaces().Create(nsWithLabel)
			Expect(err).ShouldNot(HaveOccurred())

			By("Checking new namespace has the configmap")
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", f.NumberOfNameSpace()))
			_, err = f.KubeClient.CoreV1().ConfigMaps(nsWithLabel.Name).Get(cfgMap.Name, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())

			By("Removing sync annotation")
			c, err = f.KubeClient.CoreV1().ConfigMaps(cfgMap.Namespace).Get(cfgMap.Name, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			c, _, err = core_util.PatchConfigMap(f.KubeClient, c, func(obj *core.ConfigMap) *core.ConfigMap {
				obj.Annotations = meta.RemoveKey(obj.Annotations, api.ConfigSyncKey)
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())

			By("Checking synced configmaps are deleted")
			f.EventuallyNumOfConfigmaps(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", 1))
		})
	})

	Describe("ConfigMap Syncer Backward Compatibility Test", func() {

		BeforeEach(func() {
			clusterConfig = framework.ConfigMapSyncClusterConfig()
		})

		It("Should add configmap to all namespaces", func() {

			By("Creating configmap")
			c, err := f.KubeClient.CoreV1().ConfigMaps(cfgMap.Namespace).Create(cfgMap)
			Expect(err).NotTo(HaveOccurred())

			By("Checking configmap has not synced yet")
			f.EventuallyNumOfConfigmaps(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", 1))

			By("Adding sync=true annotation")
			c, _, err = core_util.PatchConfigMap(f.KubeClient, c, func(obj *core.ConfigMap) *core.ConfigMap {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncKey, "true")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())

			By("Checking configmap has synced to all namespaces")
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", f.NumberOfNameSpace()))
		})
	})

	Describe("ConfigMap Syncer With Namespace Selector", func() {

		BeforeEach(func() {
			clusterConfig = framework.ConfigMapSyncClusterConfig()
		})

		It("Should add configmap to selected namespaces", func() {

			By("Creating configmap")
			c, err := f.KubeClient.CoreV1().ConfigMaps(cfgMap.Namespace).Create(cfgMap)
			Expect(err).NotTo(HaveOccurred())

			By("Checking configmap has not synced yet")
			f.EventuallyNumOfConfigmaps(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", 1))

			By("Adding sync annotation")
			c, _, err = core_util.PatchConfigMap(f.KubeClient, c, func(obj *core.ConfigMap) *core.ConfigMap {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncKey, "")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())

			By("Checking configmap has synced to all namespaces")
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", f.NumberOfNameSpace()))

			By("Adding selector annotation")
			c, _, err = core_util.PatchConfigMap(f.KubeClient, c, func(obj *core.ConfigMap) *core.ConfigMap {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncKey, "app="+f.App())
				return obj
			})
			Expect(err).NotTo(HaveOccurred())

			By("Checking configmap has not synced to other namespaces")
			f.EventuallyNumOfConfigmaps(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", 1))

			By("Creating new namespace with label")
			_, err = f.KubeClient.CoreV1().Namespaces().Create(nsWithLabel)
			Expect(err).ShouldNot(HaveOccurred())

			By("Checking configmap is added to only new namespace")
			f.EventuallyNumOfConfigmaps(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfConfigmaps(nsWithLabel.Name).Should(BeNumerically("==", 1))
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", 2))

			By("Changing selector annotation")
			c, _, err = core_util.PatchConfigMap(f.KubeClient, c, func(obj *core.ConfigMap) *core.ConfigMap {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncKey, "app=do-not-match")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())

			By("Checking config map not synced to other namespaces")
			f.EventuallyNumOfConfigmaps(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", 1))

			By("Removing selector annotation")
			c, _, err = core_util.PatchConfigMap(f.KubeClient, c, func(obj *core.ConfigMap) *core.ConfigMap {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncKey, "")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())

			By("Checking configmap synced to all namespaces")
			f.EventuallyNumOfConfigmaps(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", f.NumberOfNameSpace()))
		})
	})

	Describe("ConfigMap Syncer Source Deleted", func() {

		BeforeEach(func() {
			clusterConfig = framework.ConfigMapSyncClusterConfig()
		})

		It("Should delete synced configmaps from namespaces", func() {

			By("Creating configmap")
			c, err := f.KubeClient.CoreV1().ConfigMaps(cfgMap.Namespace).Create(cfgMap)
			Expect(err).NotTo(HaveOccurred())

			By("Checking configmap has not synced yet")
			f.EventuallyNumOfConfigmaps(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", 1))

			By("Adding sync annotation")
			c, _, err = core_util.PatchConfigMap(f.KubeClient, c, func(obj *core.ConfigMap) *core.ConfigMap {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncKey, "")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())

			By("Checking configmap synced to all namespaces")
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", f.NumberOfNameSpace()))

			By("Creating new namespace")
			_, err = f.KubeClient.CoreV1().Namespaces().Create(nsWithLabel)
			Expect(err).ShouldNot(HaveOccurred())

			By("Checking configmap added to new namespaces")
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", f.NumberOfNameSpace()))

			By("Deleting source configmap")
			err = f.KubeClient.CoreV1().ConfigMaps(cfgMap.Namespace).Delete(cfgMap.Name, &metav1.DeleteOptions{})
			Expect(err).ShouldNot(HaveOccurred())

			By("Checking synced configmaps has deleted")
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", 0))
		})
	})

	Describe("ConfigMap Context Syncer Test", func() {
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

		It("Should add configmap to contexts", func() {
			By("Creating source ns in remote cluster")
			f.EnsureNamespaceForContext(kubeConfigPath, context)

			By("Creating configmap")
			cfgMap, err := f.KubeClient.CoreV1().ConfigMaps(cfgMap.Namespace).Create(cfgMap)
			Expect(err).NotTo(HaveOccurred())

			By("Adding sync annotation")
			cfgMap, _, err = core_util.PatchConfigMap(f.KubeClient, cfgMap, func(obj *core.ConfigMap) *core.ConfigMap {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, api.ConfigSyncContexts, context)
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())

			By("Checking configmap added to contexts")
			f.EventuallyNumOfConfigmapsForContext(kubeConfigPath, context).Should(BeNumerically("==", 1))

			By("Removing sync annotation")
			cfgMap, _, err = core_util.PatchConfigMap(f.KubeClient, cfgMap, func(obj *core.ConfigMap) *core.ConfigMap {
				obj.Annotations = meta.RemoveKey(obj.Annotations, api.ConfigSyncContexts)
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())

			By("Checking configmap removed from contexts")
			f.EventuallyNumOfConfigmapsForContext(kubeConfigPath, context).Should(BeNumerically("==", 0))
		})
	})
})
