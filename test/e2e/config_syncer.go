package e2e

import (
	"strings"

	"github.com/appscode/kubed/pkg/config"
	"github.com/appscode/kubed/test/framework"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Config-syncer", func() {
	var (
		f                         *framework.Invocation
		configSelector            *metav1.LabelSelector
		shouldNsAndConfigmapEqual = func() {
			ns, err := f.KubeClient.CoreV1().Namespaces().List(metav1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() int {
				cfgmaps, err := f.KubeClient.CoreV1().ConfigMaps(metav1.NamespaceAll).List(metav1.ListOptions{
					LabelSelector: configSelector.String(),
				})
				Expect(err).NotTo(HaveOccurred())

				return len(cfgmaps.Items)
			}).Should(Equal(len(ns.Items)))
		}
		cfgMap = &core.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "ConfigMap",
			},
			ObjectMeta: metav1.ObjectMeta{},
			Data: map[string]string{
				"you":   "only",
				"leave": "once",
			},
		}
	)

	BeforeEach(func() {
		f = root.Invoke()
		secret := &core.Secret{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
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
		tmp, err := f.KubeClient.CoreV1().ConfigMaps(metav1.NamespaceAll).List(metav1.ListOptions{LabelSelector: metav1.FormatLabelSelector(configSelector)})
		Expect(err).NotTo(HaveOccurred())
		for _, value := range tmp.Items {
			err := f.KubeClient.CoreV1().ConfigMaps(value.Namespace).Delete(value.Name, &metav1.DeleteOptions{})
			Expect(err).NotTo(HaveOccurred())
		}
	})

	Describe("Config-syncer test", func() {
		Context("Config-sync with update config map", func() {
			BeforeEach(func() {
				cfgMap.ObjectMeta.Name = f.App()
				cfgMap.ObjectMeta.Namespace = f.Namespace()

				if cfgMap.ObjectMeta.Labels == nil {
					cfgMap.ObjectMeta.Labels = make(map[string]string)
				}
				cfgMap.ObjectMeta.Labels["app"] = f.App()

				c, err := root.KubeClient.CoreV1().ConfigMaps(root.Config.TestNamespace).Create(cfgMap)
				Expect(err).NotTo(HaveOccurred())

				metav1.SetMetaDataAnnotation(&c.ObjectMeta, config.ConfigSyncKey, "true")

				c, err = root.KubeClient.CoreV1().ConfigMaps(root.Config.TestNamespace).Update(c)
				Expect(err).NotTo(HaveOccurred())

				configSelector = metav1.SetAsLabelSelector(c.Labels)
			})
			It("Check config-syncer works", shouldNsAndConfigmapEqual)
		})

		Context("Config-sync with create config map", func() {
			BeforeEach(func() {
				cfgMap.ObjectMeta.Name = f.App()
				cfgMap.ObjectMeta.Namespace = f.Namespace()
				if cfgMap.ObjectMeta.Labels == nil {
					cfgMap.ObjectMeta.Labels = make(map[string]string)
				}
				cfgMap.ObjectMeta.Labels["app"] = f.App()

				metav1.SetMetaDataAnnotation(&cfgMap.ObjectMeta, config.ConfigSyncKey, "true")
				c, err := root.KubeClient.CoreV1().ConfigMaps(root.Config.TestNamespace).Create(cfgMap)
				Expect(err).NotTo(HaveOccurred())

				configSelector = metav1.SetAsLabelSelector(c.Labels)
			})
			It("checkout", shouldNsAndConfigmapEqual)
		})
	})
})
