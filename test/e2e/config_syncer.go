package e2e

import (
	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/kubed/pkg/config"
	"github.com/appscode/kubed/test/framework"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

var _ = Describe("Config-syncer", func() {
	var (
		f               *framework.Invocation
		cfgMap          *core.ConfigMap
		numOfNamespaces = func() int {
			ns, err := f.KubeClient.CoreV1().Namespaces().List(metav1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())
			return len(ns.Items)
		}
	)

	BeforeEach(func() {
		f = root.Invoke()

		cfgMap = &core.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				APIVersion: core.SchemeGroupVersion.String(),
				Kind:       "ConfigMap",
			},
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
	})

	AfterEach(func() {
		cfgmaps, err := f.KubeClient.CoreV1().ConfigMaps(metav1.NamespaceAll).List(metav1.ListOptions{
			LabelSelector: labels.Set{
				"app": f.App(),
			}.String(),
		})
		Expect(err).NotTo(HaveOccurred())
		for _, value := range cfgmaps.Items {
			err := f.KubeClient.CoreV1().ConfigMaps(value.Namespace).Delete(value.Name, &metav1.DeleteOptions{})
			Expect(err).NotTo(HaveOccurred())
		}
	})

	Describe("Config-syncer test", func() {
		Context("Config-sync with update config map", func() {
			var (
				namespace string
			)

			BeforeEach(func() {
				namespace = rand.WithUniqSuffix("test-ns")
			})

			It("Check config-syncer works", func() {
				c, err := root.KubeClient.CoreV1().ConfigMaps(root.Config.TestNamespace).Create(cfgMap)
				Expect(err).NotTo(HaveOccurred())
				f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", 1))

				metav1.SetMetaDataAnnotation(&c.ObjectMeta, config.ConfigSyncKey, "true")
				c, err = root.KubeClient.CoreV1().ConfigMaps(root.Config.TestNamespace).Update(c)
				Expect(err).NotTo(HaveOccurred())

				nsObj := &core.Namespace{
					TypeMeta: metav1.TypeMeta{
						APIVersion: core.SchemeGroupVersion.String(),
						Kind:       "Namespace",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: namespace,
					},
				}
				_, err = root.KubeClient.CoreV1().Namespaces().Create(nsObj)
				Expect(err).ShouldNot(HaveOccurred())
				f.EventuallyNumOfConfigmaps(namespace).Should(BeNumerically("==", 1))

				f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", numOfNamespaces()))

				metav1.SetMetaDataAnnotation(&cfgMap.ObjectMeta, config.ConfigOriginKey, "false")
				_, err = root.KubeClient.CoreV1().ConfigMaps(root.Config.TestNamespace).Update(cfgMap)
				f.EventuallyNumOfConfigmaps(metav1.NamespaceAll).Should(BeNumerically("==", 1))
			})

			AfterEach(func() {
				err := f.KubeClient.CoreV1().Namespaces().Delete(namespace, &metav1.DeleteOptions{})
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
