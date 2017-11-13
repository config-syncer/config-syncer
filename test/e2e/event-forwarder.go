package e2e

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/appscode/kubed/test/framework"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Event-forwarder", func() {
	var (
		f        *framework.Invocation
		requests []*http.Request
		s        *http.Server
		handler  = func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%q", r.URL)
		}
	)
	BeforeEach(func() {
		f = root.Invoke()
		mux := http.NewServeMux()

		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			requests = append(requests, r)
			fmt.Fprintf(w, "%q", r.URL)
		})

		s = &http.Server{
			Addr:           ":8181",
			Handler:        mux,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		go s.ListenAndServe()
		notifierConfig := &core.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: metav1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "notifier-config",
				Namespace: metav1.NamespaceSystem,
			},
			Data: map[string][]byte{
				"WEBHOOK_URL": []byte("http://localhost:8181"),
			},
		}

		_, err := f.KubeClient.CoreV1().Secrets(metav1.NamespaceSystem).Get("notifier-config", metav1.GetOptions{})
		if err != nil {
			_, err = f.KubeClient.CoreV1().Secrets(metav1.NamespaceSystem).Create(notifierConfig)
		} else {
			_, err = f.KubeClient.CoreV1().Secrets(metav1.NamespaceSystem).Update(notifierConfig)
		}
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Checkout event-forward", func() {
		Context("Pvc add eventer", func() {
			BeforeEach(func() {
				myclaim := &core.PersistentVolumeClaim{
					TypeMeta: metav1.TypeMeta{
						APIVersion: metav1.SchemeGroupVersion.String(),
						Kind:       "PersistentVolumeClaim",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "myclaim",
						Namespace: metav1.NamespaceSystem,
					},
					Spec: core.PersistentVolumeClaimSpec{
						AccessModes: append([]core.PersistentVolumeAccessMode{}, "ReadWriteOnce"),
						Resources: core.ResourceRequirements{
							Requests: core.ResourceList{
								"storage": resource.Quantity{},
							},
						},
					},
				}
				_, err := f.KubeClient.CoreV1().PersistentVolumeClaims(metav1.NamespaceSystem).Get("myclaim", metav1.GetOptions{})
				if err == nil {
					err = f.KubeClient.CoreV1().PersistentVolumeClaims(metav1.NamespaceSystem).Delete("myclaim", &metav1.DeleteOptions{})
				}
				_, err = f.KubeClient.CoreV1().PersistentVolumeClaims(metav1.NamespaceSystem).Create(myclaim)
				Expect(err).NotTo(HaveOccurred())
			})

			FIt("Check notify kubed", func() {
				Eventually(func() bool {
					for _, val := range requests {
						wr := httptest.NewRecorder()
						handler(wr, val)
						result := wr.Result()
						bit, err := ioutil.ReadAll(result.Body)
						Expect(err).NotTo(HaveOccurred())
						respStr := string(bit)
						if strings.Contains(respStr, "PersistentVolumeClaim") && result.StatusCode == 200 {
							return true
						}
					}
					return false
				}).Should(BeTrue())
				Expect(0).Should(BeZero())
			})
		})
		Context("Warning Event", func() {
			BeforeEach(func() {
				wPod := &core.Pod{
					TypeMeta: metav1.TypeMeta{
						APIVersion: metav1.SchemeGroupVersion.String(),
						Kind:       "Pod",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "busybox",
						Namespace: metav1.NamespaceSystem,
					},
					Spec: core.PodSpec{
						RestartPolicy: "Never",
						Containers: append([]core.Container{}, core.Container{
							Name:            "busybox",
							Image:           "busybox",
							ImagePullPolicy: "IfNotPresent",
							Command:         []string{"bad", "3600"},
						}),
					},
				}
				_, err := f.KubeClient.CoreV1().Pods(metav1.NamespaceSystem).Create(wPod)
				Expect(err).NotTo(HaveOccurred())
			})
			AfterEach(func() {
				err := f.KubeClient.CoreV1().Pods(metav1.NamespaceSystem).Delete("busybox", &metav1.DeleteOptions{})
				Expect(err).NotTo(HaveOccurred())
			})

			FIt("Check warning event", func() {
				Eventually(func() bool {
					for _, val := range requests {
						wr := httptest.NewRecorder()
						handler(wr, val)
						resp := wr.Result()
						byt, err := ioutil.ReadAll(resp.Body)
						Expect(err).NotTo(HaveOccurred())
						respStr := string(byt)
						if resp.StatusCode == 200 && strings.Contains(respStr, "busybox") {
							return true
						}
					}
					return false
				}).Should(BeTrue())
			})
		})
	})
})
