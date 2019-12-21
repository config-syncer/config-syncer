/*
Copyright The Kubed Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package framework

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/go/log"
	"github.com/appscode/go/types"

	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
)

var (
	handler = func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%q", r.URL)
	}
)

func (f *Invocation) RunWebhookServer(stopCh <-chan os.Signal, requests *[]*http.Request) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		*requests = append(*requests, r)
		fmt.Fprintf(w, "%q", r.URL)
	})

	srv := &http.Server{
		Addr:           ":8181",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// Error starting or closing listener:
			log.Errorln("ListenAndServe error. Reason: %v", err)
		}
	}()

	go func() {
		<-stopCh
		fmt.Println("Clossing webhook server....")
		srv.Shutdown(context.Background())
	}()
}

func (f *Invocation) ServiceForWebhook() *core.Service {
	return &core.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      f.app,
			Namespace: f.namespace,
		},
		Spec: core.ServiceSpec{
			Ports: []core.ServicePort{
				{
					Name:       "webhook",
					Port:       80,
					Protocol:   core.ProtocolTCP,
					TargetPort: intstr.FromInt(8181),
				},
			},
			Type: core.ServiceTypeClusterIP,
		},
	}
}

func (f *Invocation) LocalEndPointForWebhook(service *core.Service) *core.Endpoints {
	return &core.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.Name,
			Namespace: service.Namespace,
		},
		Subsets: []core.EndpointSubset{
			{
				Addresses: []core.EndpointAddress{
					{
						IP: "10.0.2.2",
					},
				},
				Ports: []core.EndpointPort{
					{
						Name:     "webhook",
						Port:     8181,
						Protocol: core.ProtocolTCP,
					},
				},
			},
		},
	}
}

func (f *Invocation) NewPersistentVolumeClaim() *core.PersistentVolumeClaim {
	return &core.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rand.WithUniqSuffix(f.app + "-"),
			Namespace: f.namespace,
		},
		Spec: core.PersistentVolumeClaimSpec{
			AccessModes: []core.PersistentVolumeAccessMode{
				core.ReadWriteOnce,
			},
			Resources: core.ResourceRequirements{
				Requests: core.ResourceList{
					core.ResourceName(core.ResourceStorage): resource.MustParse("2Gi"),
				},
			},
			StorageClassName: types.StringP(StandardStorageClass),
		},
	}
}

func (f *Invocation) NewPod() *core.Pod {
	podTemplate := f.PodTemplate()
	pod := &core.Pod{}
	pod.Name = rand.WithUniqSuffix(f.app + "-")
	pod.Namespace = f.namespace
	pod.Spec = podTemplate.Spec
	pod.Spec.Containers = []core.Container{f.BusyboxContainerWithBadCommand()}

	return pod
}

func (f *Invocation) BusyboxContainerWithBadCommand() core.Container {
	return core.Container{
		Name:            "busybox",
		Image:           "busybox",
		ImagePullPolicy: core.PullIfNotPresent,
		Command:         []string{"bad", "3600"},
	}
}

func (f *Invocation) EventuallyNotifiedToWebhookServer(requests *[]*http.Request, expetedSubstr string) GomegaAsyncAssertion {
	return Eventually(func() bool {
		for _, val := range *requests {
			wr := httptest.NewRecorder()
			handler(wr, val)
			result := wr.Result()
			bit, err := ioutil.ReadAll(result.Body)
			Expect(err).NotTo(HaveOccurred())
			respStr := string(bit)
			if strings.Contains(respStr, expetedSubstr) && result.StatusCode == 200 {
				return true
			}
		}
		return false
	})
}

func (f *Invocation) CreatePod(pod *core.Pod) (*core.Pod, error) {
	return f.KubeClient.CoreV1().Pods(pod.Namespace).Create(pod)
}

func (f *Invocation) DeletePod(meta metav1.ObjectMeta) error {
	return f.KubeClient.CoreV1().Pods(meta.Namespace).Delete(meta.Name, deleteInBackground())
}

func (f *Invocation) DeletePersistentVolumeClaim(meta metav1.ObjectMeta) error {
	return f.KubeClient.CoreV1().PersistentVolumeClaims(meta.Namespace).Delete(meta.Name, deleteInBackground())
}

func (f *Invocation) WaitUntilPodTerminated(meta metav1.ObjectMeta) error {
	return wait.PollImmediate(interval, timeout, func() (done bool, err error) {
		if _, err := f.KubeClient.CoreV1().Pods(meta.Namespace).Get(meta.Name, metav1.GetOptions{}); err != nil {
			if kerr.IsNotFound(err) {
				return true, nil
			} else {
				return true, err
			}
		}
		return false, nil
	})
}
