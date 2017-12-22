package e2e

import (
	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/kubed/pkg/config"
	"github.com/appscode/kubed/pkg/util"
	"github.com/appscode/kubed/test/framework"
	core_util "github.com/appscode/kutil/core/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Config-syncer", func() {
	var (
		f               *framework.Invocation
		cfgMap          *core.ConfigMap
		nsWithLabel     *core.Namespace
		numOfNamespaces = func() int {
			ns, err := f.KubeClient.CoreV1().Namespaces().List(metav1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())
			return len(ns.Items)
		}
	)

	BeforeEach(func() {
		f = root.Invoke()

		cfgMap = &core.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      f.App(),
				Namespace: f.Namespace(),
				Labels: map[string]string{
					"app": f.App(),
				},
			},
			Data: map[string]string{
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
		f.DeleteAllConfigmaps()

		err := f.KubeClient.CoreV1().Namespaces().Delete(nsWithLabel.Name, &metav1.DeleteOptions{})
		if kerr.IsNotFound(err) {
			err = nil
		}
		Expect(err).NotTo(HaveOccurred())
		f.EventuallyNamespaceDeleted(nsWithLabel.Name).Should(BeTrue())
	})

	Describe("ConfigMap Syncer Test", func() {
		It("Should add configmap to all namespaces", func() {
			By("Creating configmap")
			c, err := root.KubeClient.CoreV1().ConfigMaps(cfgMap.Namespace).Create(cfgMap)
			Expect(err).NotTo(HaveOccurred())
			f.EventuallyNumOfConfigmaps(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", 1))

			By("Adding sync annotation")
			c, err = core_util.PatchConfigMap(f.KubeClient, c, func(obj *core.ConfigMap) *core.ConfigMap {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, config.ConfigSyncKey, "")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", numOfNamespaces()))

			By("Creating new namespace")
			_, err = root.KubeClient.CoreV1().Namespaces().Create(nsWithLabel)
			Expect(err).ShouldNot(HaveOccurred())
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", numOfNamespaces()))

			By("Removing sync annotation")
			c, err = core_util.PatchConfigMap(f.KubeClient, c, func(obj *core.ConfigMap) *core.ConfigMap {
				obj.Annotations = util.RemoveKey(obj.Annotations, config.ConfigSyncKey)
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())
			f.EventuallyNumOfConfigmaps(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", 1))
		})
	})

	Describe("ConfigMap Syncer Backward Compatibility Test", func() {
		It("Should add configmap to all namespaces", func() {
			By("Creating configmap")
			c, err := root.KubeClient.CoreV1().ConfigMaps(cfgMap.Namespace).Create(cfgMap)
			Expect(err).NotTo(HaveOccurred())
			f.EventuallyNumOfConfigmaps(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", 1))

			By("Adding sync=true annotation")
			c, err = core_util.PatchConfigMap(f.KubeClient, c, func(obj *core.ConfigMap) *core.ConfigMap {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, config.ConfigSyncKey, "true")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", numOfNamespaces()))
		})
	})

	Describe("ConfigMap Syncer With Namespace Selector", func() {
		It("Should add configmap to selected namespaces", func() {
			By("Creating configmap")
			c, err := root.KubeClient.CoreV1().ConfigMaps(cfgMap.Namespace).Create(cfgMap)
			Expect(err).NotTo(HaveOccurred())
			f.EventuallyNumOfConfigmaps(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", 1))

			By("Adding sync annotation")
			c, err = core_util.PatchConfigMap(f.KubeClient, c, func(obj *core.ConfigMap) *core.ConfigMap {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, config.ConfigSyncKey, "")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", numOfNamespaces()))

			By("Adding selector annotation")
			c, err = core_util.PatchConfigMap(f.KubeClient, c, func(obj *core.ConfigMap) *core.ConfigMap {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, config.ConfigSyncKey, "app="+f.App())
				return obj
			})
			Expect(err).NotTo(HaveOccurred())
			f.EventuallyNumOfConfigmaps(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", 1))

			By("Creating new namespace with label")
			_, err = root.KubeClient.CoreV1().Namespaces().Create(nsWithLabel)
			Expect(err).ShouldNot(HaveOccurred())
			f.EventuallyNumOfConfigmaps(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfConfigmaps(nsWithLabel.Name).Should(BeNumerically("==", 1))
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", 2))

			By("Changing selector annotation")
			c, err = core_util.PatchConfigMap(f.KubeClient, c, func(obj *core.ConfigMap) *core.ConfigMap {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, config.ConfigSyncKey, "app=do-not-match")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())
			f.EventuallyNumOfConfigmaps(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", 1))

			By("Removing selector annotation")
			c, err = core_util.PatchConfigMap(f.KubeClient, c, func(obj *core.ConfigMap) *core.ConfigMap {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, config.ConfigSyncKey, "")
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())
			f.EventuallyNumOfConfigmaps(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", numOfNamespaces()))
		})
	})

	Describe("ConfigMap Context Syncer Test", func() {
		It("Should add configmap to contexts", func() {
			kubeConfigPath := "/home/dipta/.kube/config"
			context := "kubed-test"

			By("Creating client for context")
			client, ns, err := util.ClientAndNamespaceForContext(kubeConfigPath, context)
			Expect(err).ShouldNot(HaveOccurred())

			if ns == "" {
				ns = f.Namespace()
			}
			By("Using external context " + context + " with namespace " + ns)

			By("Creating configmap")
			c, err := root.KubeClient.CoreV1().ConfigMaps(cfgMap.Namespace).Create(cfgMap)
			Expect(err).NotTo(HaveOccurred())
			f.EventuallyNumOfConfigmaps(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", 1))

			By("Adding sync annotation")
			c, err = core_util.PatchConfigMap(f.KubeClient, c, func(obj *core.ConfigMap) *core.ConfigMap {
				metav1.SetMetaDataAnnotation(&obj.ObjectMeta, config.ConfigSyncContexts, context)
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())
			f.EventuallyNumOfConfigmaps(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfConfigmapsForClient(client, ns).Should(BeNumerically("==", 1))

			By("Removing sync annotation")
			c, err = core_util.PatchConfigMap(f.KubeClient, c, func(obj *core.ConfigMap) *core.ConfigMap {
				obj.Annotations = util.RemoveKey(obj.Annotations, config.ConfigSyncContexts)
				return obj
			})
			Expect(err).ShouldNot(HaveOccurred())
			f.EventuallyNumOfConfigmaps(f.Namespace()).Should(BeNumerically("==", 1))
			f.EventuallyNumOfConfigmapsForClient(client, ns).Should(BeNumerically("==", 0))
		})
	})
})
