package e2e

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"io/ioutil"
	"github.com/appscode/kubed/test/framework"
)

var _ = Describe("Config-syncer", func() {
	var (
		f                         *framework.Invocation
		configSelector            *metav1.LabelSelector
		shouldNsAndConfigmapEqual = func() {
			ns, err := f.KubeClient.CoreV1().Namespaces().List(metav1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() int {
				tmp, err := f.KubeClient.CoreV1().ConfigMaps(metav1.NamespaceAll).List(metav1.ListOptions{LabelSelector: metav1.FormatLabelSelector(configSelector)})
				Expect(err).NotTo(HaveOccurred())
				return len(tmp.Items)
			}, "10m", "5s").Should(Equal(len(ns.Items)))
		}
	)

	BeforeEach(func() {
		f = root.Invoke()
		file, err := ioutil.ReadFile(filepath.Join(homedir.HomeDir(), "go/src/github.com/appscode/kubed/docs/examples/config-syncer/config.yaml"))
		Expect(err).NotTo(HaveOccurred())
		secret := &apiv1.Secret{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Secret",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "kubed-config",
				Namespace: "kube-system",
				Labels: map[string]string{
					"app": "kubed",
				},
			},
			Data: map[string][]byte{
				"config.yaml": file,
			},
		}

		_, err = f.KubeClient.CoreV1().Secrets("kube-system").Update(secret)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		tmp, err := f.KubeClient.CoreV1().ConfigMaps(metav1.NamespaceAll).List(metav1.ListOptions{LabelSelector: metav1.FormatLabelSelector(configSelector)})
		Expect(err).NotTo(HaveOccurred())
		for _, value := range tmp.Items {
			err := f.KubeClient.CoreV1().ConfigMaps(value.Namespace).Delete(value.Name, &metav1.DeleteOptions{})
			Expect(err).NotTo(HaveOccurred())
		}

		tmp, err = f.KubeClient.CoreV1().ConfigMaps(metav1.NamespaceAll).List(metav1.ListOptions{LabelSelector: metav1.FormatLabelSelector(configSelector)})
		Expect(err).NotTo(HaveOccurred())

		Expect(len(tmp.Items)).Should(Equal(0))
	})

	Describe("Config-syncer test", func() {
		Context("Config-sync with update config map", func() {
			BeforeEach(func() {
				cfgMap := &apiv1.ConfigMap{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "v1",
						Kind:       "ConfigMap",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      f.App(),
						Namespace: f.Config.TestNamespace,
						Labels: map[string]string{
							"app": f.App(),
						},
					},
					Data: map[string]string{
						"you":   "only",
						"leave": "once",
					},
				}
				c, err := root.KubeClient.CoreV1().ConfigMaps(root.Config.TestNamespace).Create(cfgMap)
				Expect(err).NotTo(HaveOccurred())

				metav1.SetMetaDataAnnotation(&c.ObjectMeta, "kubed.appscode.com/sync", "true")

				c, err = root.KubeClient.CoreV1().ConfigMaps(root.Config.TestNamespace).Update(c)
				Expect(err).NotTo(HaveOccurred())

				configSelector = metav1.SetAsLabelSelector(c.Labels)
			})
			It("Check config-syncer works", shouldNsAndConfigmapEqual)
		})

		Context("Config-sync with create config map", func() {
			BeforeEach(func() {
				cfgMap := &apiv1.ConfigMap{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "v1",
						Kind:       "ConfigMap",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      f.App(),
						Namespace: f.Config.TestNamespace,
						Labels: map[string]string{
							"app": f.App(),
						},
						Annotations: map[string]string{
							"kubed.appscode.com/sync": "true",
						},
					},
					Data: map[string]string{
						"you":   "only",
						"leave": "once",
					},
				}
				c, err := root.KubeClient.CoreV1().ConfigMaps(root.Config.TestNamespace).Create(cfgMap)
				Expect(err).NotTo(HaveOccurred())

				configSelector = metav1.SetAsLabelSelector(c.Labels)
			})
			It("checkout", shouldNsAndConfigmapEqual)
		})
	})
})
