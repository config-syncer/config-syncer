package e2e

import (
	"time"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/kubed/test/framework"
	prom "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1"
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
	})

	Describe("Kubed api server: Search", func() {
		Context("Dashboard search", func() {
			var (
				svcName    string
				deployName string
			)
			BeforeEach(func() {
				svcName = rand.WithUniqSuffix("kubed-svc")
				service := &core.Service{
					TypeMeta: metav1.TypeMeta{
						APIVersion: core.SchemeGroupVersion.String(),
						Kind:       "Service",
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
						Ports: []core.ServicePort{
							{
								Protocol: core.ProtocolTCP,
								Port:     80,
							},
						},
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
								Containers: []core.Container{
									{
										Name:  "nginx",
										Image: "nginx:1.7.9",
										Ports: []core.ContainerPort{
											{
												ContainerPort: 80,
											},
										},
									},
								},
							},
						},
					},
				}

				_, err := f.KubeClient.ExtensionsV1beta1().Deployments(f.Namespace()).Create(deploy)
				Expect(err).NotTo(HaveOccurred())

				_, err = f.KubeClient.CoreV1().Services(f.Namespace()).Create(service)
				Expect(err).NotTo(HaveOccurred())
				time.Sleep(2 * time.Second)

				servicemonitor := &prom.ServiceMonitor{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "monitoring.coreos.com/v1",
						Kind:       prom.ServiceMonitorsKind,
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      rand.WithUniqSuffix("test-svcmtr"),
						Namespace: f.Namespace(),
						Labels: map[string]string{
							"app": svcName,
						},
					},
					Spec: prom.ServiceMonitorSpec{
						Selector: metav1.LabelSelector{
							MatchLabels: map[string]string{
								"app": svcName,
							},
						},
					},
				}
				_, err = f.KubedOperator.PromClient.ServiceMonitors(f.Namespace()).Create(servicemonitor)
				Expect(err).NotTo(HaveOccurred())

				prometheus := &prom.Prometheus{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "monitoring.coreos.com/v1",
						Kind:       prom.PrometheusesKind,
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      rand.WithUniqSuffix("test-prom"),
						Namespace: f.Namespace(),
						Labels: map[string]string{
							"app": svcName,
						},
					},
					Spec: prom.PrometheusSpec{
						ServiceMonitorSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"app": svcName,
							},
						},
					},
				}
				_, err = f.KubedOperator.PromClient.Prometheuses(f.Namespace()).Create(prometheus)
				Expect(err).NotTo(HaveOccurred())
			})

			It("Checkout reverse index", func() {
				pods, err := f.KubeClient.CoreV1().Pods(f.Namespace()).List(metav1.ListOptions{
					LabelSelector: labels.SelectorFromSet(map[string]string{"app": svcName}).String(),
				})
				Expect(err).NotTo(HaveOccurred())
				path := "/api/v1/namespaces/" + pods.Items[0].Namespace + "/pods/" + pods.Items[0].Name + "/services"
				f.EventuallyReverseIndex(path).Should(BeNumerically("==", 200)) // TODO: check response body

				svcs, err := f.KubeClient.CoreV1().Services(f.Namespace()).List(metav1.ListOptions{
					LabelSelector: labels.SelectorFromSet(map[string]string{"app": svcName}).String(),
				})
				Expect(err).NotTo(HaveOccurred())
				path = "/apis/" + prom.Group + "/" + prom.Version + "/namespaces/" + svcs.Items[0].Namespace + "/services/" + svcName + "/" + prom.ServiceMonitorName
				f.EventuallyReverseIndex(path).Should(BeNumerically("==", 200)) // TODO: check response body

				proms, err := f.KubeClient.CoreV1().Services(f.Namespace()).List(metav1.ListOptions{
					LabelSelector: labels.SelectorFromSet(map[string]string{"app": svcName}).String(),
				})
				Expect(err).NotTo(HaveOccurred())
				path = "/apis/" + prom.Group + "/" + prom.Version + "/namespaces/" + proms.Items[0].Namespace + "/services/" + svcName + "/" + prom.PrometheusName
				f.EventuallyReverseIndex(path).Should(BeNumerically("==", 200)) // TODO: check response body
			})
		})
	})
})
