package e2e

import (
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/kubed/test/framework"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

var _ = Describe("Kubed api server", func() {
	var (
		f *framework.Invocation
	)
	BeforeEach(func() {
		f = root.Invoke()
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
apiServer:
  address: :8080
  enableReverseIndex: true
  enableSearchIndex: true
`),
			},
		}

		_, err := f.KubeClient.CoreV1().Secrets(metav1.NamespaceSystem).Update(secret)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {

	})

	Describe("Kubed api server: Search", func() {
		Context("Dashboard search", func() {
			var (
				KubedEnpoint []string
				svcName      string
				deployName   string
				request      *http.Request
			)
			BeforeEach(func() {
				kubedSvc, err := f.KubeClient.CoreV1().Services(metav1.NamespaceSystem).Get("kubed-operator", metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())

				kubedSvc.Spec.Type = core.ServiceTypeNodePort
				_, err = f.KubeClient.CoreV1().Services(metav1.NamespaceSystem).Update(kubedSvc)
				Expect(err).NotTo(HaveOccurred())

				svcName = rand.WithUniqSuffix("kubed-svc")
				service := &core.Service{
					TypeMeta: metav1.TypeMeta{
						APIVersion: core.SchemeGroupVersion.String(),
						Kind:       "Pod",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      svcName,
						Namespace: f.Namespace(),
						Labels: map[string]string{
							"app": svcName,
						},
					},
					Spec: core.ServiceSpec{
						Selector: map[string]string{
							"app": svcName,
						},
						Ports: append([]core.ServicePort{}, core.ServicePort{
							Protocol: "TCP",
							Port:     80,
						}),
					},
				}

				deployName = rand.WithUniqSuffix("kubed-deploy")
				deploy := &v1beta1.Deployment{
					TypeMeta: metav1.TypeMeta{
						APIVersion: v1beta1.SchemeGroupVersion.String(),
						Kind:       "Deployment",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      deployName,
						Namespace: f.Namespace(),
						Labels: map[string]string{
							"app": svcName,
						},
					},
					Spec: v1beta1.DeploymentSpec{
						Selector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"app": svcName,
							},
						},
						Template: core.PodTemplateSpec{
							ObjectMeta: metav1.ObjectMeta{
								Labels: map[string]string{
									"app": svcName,
								},
							},
							Spec: core.PodSpec{
								Containers: append([]core.Container{}, core.Container{
									Name:  "nginx",
									Image: "nginx:1.7.9",
									Ports: append([]core.ContainerPort{}, core.ContainerPort{
										ContainerPort: 80,
									}),
								}),
							},
						},
					},
				}

				_, err = f.KubeClient.ExtensionsV1beta1().Deployments(f.Namespace()).Create(deploy)
				Expect(err).NotTo(HaveOccurred())

				_, err = f.KubeClient.CoreV1().Services(f.Namespace()).Create(service)
				Expect(err).NotTo(HaveOccurred())

				time.Sleep(2 * time.Second)

				pods, err := f.KubeClient.CoreV1().Pods(f.Namespace()).List(metav1.ListOptions{
					LabelSelector: labels.SelectorFromSet(map[string]string{"app": svcName}).String(),
				})
				Expect(err).NotTo(HaveOccurred())

				path := "/api/v1/namespaces/" + pods.Items[0].Namespace + "/pods/" + pods.Items[0].Name + "/services"
				Expect(len(KubedEnpoint)).Should(BeNumerically(">=", 1))

				request, err = http.NewRequest(http.MethodGet, KubedEnpoint[0]+path, nil)
				Expect(err).NotTo(HaveOccurred())
			})

			It("Checkout reverse index", func() {
				Eventually(func() error {
					var outputs []byte
					outputs, err := exec.Command(
						"/usr/local/bin/minikube",
						"service",
						"kubed-operator",
						"--url",
						"-n",
						metav1.NamespaceSystem,
					).CombinedOutput()
					if err == nil {
						for _, output := range strings.Split(string(outputs), "\n") {
							if strings.HasPrefix(output, "http") {
								KubedEnpoint = append(KubedEnpoint, output)
							}
						}
						return nil
					}
					return err
				}, framework.DefaultEventuallyTimeout, framework.DefaultEventuallyPollingInterval).Should(BeNil())

				Eventually(func() int {
					resp, err := http.DefaultClient.Do(request)
					Expect(err).NotTo(HaveOccurred())
					return resp.StatusCode
				}).Should(BeNumerically("==", 200)) // TODO: check response body
			})
		})
	})
})
