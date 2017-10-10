package e2e

import (
	. "github.com/onsi/gomega"
	. "github.com/onsi/ginkgo"
	"github.com/appscode/kubed/test/framework"
	"path/filepath"
	"io/ioutil"
	"k8s.io/client-go/util/homedir"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/labels"
	"os/exec"
	"strings"
	"github.com/appscode/go/crypto/rand"
	"net/http"
)

var _ = Describe("Kubed api server", func() {
	var (
		f *framework.Invocation
	)
	BeforeEach(func() {
		Expect(0).Should(Equal(0))
		f = root.Invoke()
		file, err := ioutil.ReadFile(filepath.Join(homedir.HomeDir(), "go/src/github.com/appscode/kubed/docs/examples/apiserver/config.yaml"))
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
				kubedSvc, err := f.KubeClient.CoreV1().Services("kube-system").Get("kubed-operator", metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())

				kubedSvc.Spec.Type = "LoadBalancer"
				_, err = f.KubeClient.CoreV1().Services("kube-system").Update(kubedSvc)
				Expect(err).NotTo(HaveOccurred())

				Eventually(func() error {
					var outputs []byte
					outputs, err = exec.Command(
						"/usr/local/bin/minikube",
						"service",
						"kubed-operator",
						"--url",
						"-n",
						"kube-system",
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
				}, "5m", "10s").Should(BeNil())

				svcName = rand.WithUniqSuffix("kubed-svc")
				service := &apiv1.Service{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "v1",
						Kind:       "Pod",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      svcName,
						Namespace: f.Namespace(),
						Labels: map[string]string{
							"app": svcName,
						},
					},
					Spec: apiv1.ServiceSpec{
						Selector: map[string]string{
							"app": svcName,
						},
						Ports: append([]apiv1.ServicePort{}, apiv1.ServicePort{
							Protocol: "TCP",
							Port:     80,
						}),
					},
				}

				deployName = rand.WithUniqSuffix("kubed-deploy")
				deploy := &v1beta1.Deployment{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "extensions/v1beta1",
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
						Template: apiv1.PodTemplateSpec{
							ObjectMeta: metav1.ObjectMeta{
								Labels: map[string]string{
									"app": svcName,
								},
							},
							Spec: apiv1.PodSpec{
								Containers: append([]apiv1.Container{}, apiv1.Container{
									Name:  "nginx",
									Image: "nginx:1.7.9",
									Ports: append([]apiv1.ContainerPort{}, apiv1.ContainerPort{
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
				Eventually(func() int {
					resp, err := http.DefaultClient.Do(request)
					Expect(err).NotTo(HaveOccurred())
					_, err = ioutil.ReadAll(resp.Body)
					Expect(err).NotTo(HaveOccurred())
					return resp.StatusCode
				}).Should(BeNumerically("==", 200))
			})
		})
	})
})
