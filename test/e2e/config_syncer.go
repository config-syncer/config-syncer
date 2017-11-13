package e2e

import (
	"strings"

	"github.com/appscode/kubed/pkg/config"
	"github.com/appscode/kubed/test/framework"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"github.com/appscode/go/crypto/rand"
)

var _ = Describe("Config-syncer", func() {
	var (
		f              *framework.Invocation
		cfgMap         *core.ConfigMap
		namespaceCount = func() int {
			ns, err := f.KubeClient.CoreV1().Namespaces().List(metav1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())
			return len(ns.Items)
		}
		configmapCounts = func(nsName string) GomegaAsyncAssertion {
			return Eventually(func() int {
				cfgmaps, err := f.KubeClient.CoreV1().ConfigMaps(nsName).List(metav1.ListOptions{
					LabelSelector: labels.Set{
						"app": f.App(),
					}.String(),
				})
				Expect(err).NotTo(HaveOccurred())
				return len(cfgmaps.Items)
			})
		}

		shouldNsAndConfigmapEqual = func() {
			ns, err := f.KubeClient.CoreV1().Namespaces().List(metav1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() int {
				cfgmaps, err := f.KubeClient.CoreV1().ConfigMaps(metav1.NamespaceAll).List(metav1.ListOptions{
					LabelSelector: labels.Set{
						"app": f.App(),
					}.String(),
				})
				Expect(err).NotTo(HaveOccurred())
				return len(cfgmaps.Items)
			}).Should(Equal(len(ns.Items)))
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

		secret := &core.Secret{
			TypeMeta: metav1.TypeMeta{
				APIVersion: core.SchemeGroupVersion.String(),
				Kind:       "Secret",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "kubed-config",
				Namespace: metav1.NamespaceSystem,
				Labels: map[string]string{
					"app": "kubed",
				},
			},
			StringData: map[string]string{
				"config.yaml": strings.TrimSpace(`
enableConfigSyncer: true
`),
			},
		}

		_, err := f.KubeClient.CoreV1().Secrets(metav1.NamespaceSystem).Update(secret)
		Expect(err).NotTo(HaveOccurred())
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
			BeforeEach(func() {
				c, err := root.KubeClient.CoreV1().ConfigMaps(root.Config.TestNamespace).Create(cfgMap)
				Expect(err).NotTo(HaveOccurred())

				// TODO: check that it only exists in one namespace
				cfgmaps, err := root.KubeClient.CoreV1().ConfigMaps(metav1.NamespaceAll).List(metav1.ListOptions{
					LabelSelector: labels.Set{
						"app": f.App(),
					}.String(),
				})
				Expect(len(cfgmaps.Items)).Should(BeNumerically("==", 1))

				metav1.SetMetaDataAnnotation(&c.ObjectMeta, config.ConfigSyncKey, "true")
				c, err = root.KubeClient.CoreV1().ConfigMaps(root.Config.TestNamespace).Update(c)
				Expect(err).NotTo(HaveOccurred())
			})

			// Now create a NS, check that configmap is copied into that ns
			// Remove annotations, confirm only the original source is present
			It("Check config-syncer works", func() {
				nsName := rand.WithUniqSuffix("test-ns")
				namespace := &core.Namespace{
					TypeMeta: metav1.TypeMeta{
						APIVersion: core.SchemeGroupVersion.String(),
						Kind:       "Namespace",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: nsName,
					},
				}
				_, err := root.KubeClient.CoreV1().Namespaces().Create(namespace)
				Expect(err).ShouldNot(HaveOccurred())
				configmapCounts(nsName).Should(BeNumerically("==", 1))

				nsCount := namespaceCount()
				configmapCounts("").Should(BeNumerically("==", nsCount))

				metav1.SetMetaDataAnnotation(&cfgMap.ObjectMeta, config.ConfigOriginKey, "false")
				_, err = root.KubeClient.CoreV1().ConfigMaps(root.Config.TestNamespace).Update(cfgMap)
				configmapCounts("").Should(BeNumerically("==", 1))
			})
		})
	})
})
