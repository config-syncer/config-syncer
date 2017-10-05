package e2e

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/util/homedir"
	// "github.com/appscode/kubed/test/framework"
	"path/filepath"
	// "github.com/ghodss/yaml"
	"io/ioutil"
	"time"
)

var _ = Describe("Book", func() {
	BeforeEach(func() {

	})
	JustBeforeEach(func() {

	})

	Describe("Create Secret", func() {
		It("Create kubed-config Secret", func() {
			file, err := ioutil.ReadFile(filepath.Join(homedir.HomeDir(), "go/src/github.com/appscode/kubed/docs/examples/config-syncer/config.yaml"))
			Expect(err).NotTo(HaveOccurred())
			cfgMap := apiv1.Secret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind: 		"Secret",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "kubed-config",
					Namespace: "kube-system",
					Labels: map[string]string{
						"app": "kubed",
					},
				},
				Data: map[string][]byte {
					"config.yaml": file,
				},
			}

			secret, err := root.KubeClient.CoreV1().Secrets("kube-system").Update(&cfgMap)
			Expect(err).NotTo(HaveOccurred())
			fmt.Printf("%s Successfully created in %s namespace\n", secret.Name, secret.Namespace)
		})

		It("Create a other config map, which will sync all namespaces", func() {
			cfgMap := apiv1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind: "ConfigMap",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "omni",
					Namespace: root.Config.TestNamespace,
				},
				Data: map[string]string{
					"you": "only",
					"leave": "once",
				},
			}
			cm, err := root.KubeClient.CoreV1().ConfigMaps(root.Config.TestNamespace).Create(&cfgMap)
			Expect(err).NotTo(HaveOccurred())
			fmt.Printf("ConfigMap %s created successfully in %s namespace\n", cm.Name, cm.Namespace)

			cfgMap = apiv1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind: "ConfigMap",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "omni",
					Namespace: root.Config.TestNamespace,
					Annotations: map[string]string {
						"kubed.appscode.com/sync": "true",
					},
				},
				Data: map[string]string{
					"you": "only",
					"leave": "once",
				},
			}
			cm, err = root.KubeClient.CoreV1().ConfigMaps(root.Config.TestNamespace).Update(&cfgMap)
			Expect(err).NotTo(HaveOccurred())
			fmt.Printf("ConfigMap %s successfully updated in %s namespace", cm.Name, cm.Namespace)
		})

		It("Checkout config-syncer works", func() {
			time.Sleep(15 * time.Second)
			ns, err := root.KubeClient.CoreV1().Namespaces().List(metav1.ListOptions{})
			Expect(err).NotTo(HaveOccurred())
			for _, value := range ns.Items {
				tmp, err := root.KubeClient.CoreV1().ConfigMaps(value.Name).Get("omni", metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(tmp.Name).Should(Equal("omni"))
			}
			root.DeleteNamespace()
		})
	})


	AfterEach(func() {
		fmt.Println("Delete namespaces ***********", root.Config.TestNamespace)
	})
})
