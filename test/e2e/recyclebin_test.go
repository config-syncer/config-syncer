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

package e2e_test

import (
	"os"
	"time"

	api "github.com/appscode/kubed/apis/kubed/v1alpha1"
	"github.com/appscode/kubed/test/e2e/framework"
	. "github.com/appscode/kubed/test/e2e/matcher"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	core_util "kmodules.xyz/client-go/core/v1"
)

var _ = Describe("RecycleBin", func() {
	var (
		f             *framework.Invocation
		configMap     *core.ConfigMap
		clusterConfig api.ClusterConfig
		stopCh        chan struct{}
	)

	BeforeEach(func() {
		f = root.Invoke()
	})

	JustBeforeEach(func() {
		if f.SelfHostedOperator {
			By("Restarting kubed operator")
			err := f.RestartKubedOperator(&clusterConfig)
			Expect(err).NotTo(HaveOccurred())
		} else {
			By("Starting Kubed")
			stopCh = make(chan struct{})
			err := f.RunKubed(stopCh, clusterConfig)
			Expect(err).NotTo(HaveOccurred())

			By("Waiting for API server to be ready")
			root.EventuallyAPIServerReady().Should(Succeed())
			time.Sleep(time.Second * 5)
		}
	})

	AfterEach(func() {
		if !f.SelfHostedOperator {
			close(stopCh)
		}
		time.Sleep(time.Second * 10)
	})

	Describe("ConfigMap", func() {

		BeforeEach(func() {
			By("Creating clusterConfiguration")
			clusterConfig = framework.RecycleBinClusterConfig()
		})

		Context("recycle ConfigMap ", func() {

			BeforeEach(func() {
				configMap = f.NewConfigMap()
			})

			AfterEach(func() {
				if !f.SelfHostedOperator {
					os.RemoveAll(clusterConfig.RecycleBin.Path)
				}
			})

			It("should store deleted configMap in RecycleBin", func() {

				By("Creating configMap: " + configMap.Name)
				cm, err := f.CreateConfigMap(configMap)
				Expect(err).NotTo(HaveOccurred())

				By("Deleting configMap: " + cm.Name)
				err = f.DeleteConfigMap(cm.ObjectMeta)
				Expect(err).NotTo(HaveOccurred())

				//give some time to recycle configMap
				time.Sleep(time.Second * 10)

				By("Checking configMap stored in RecycleBin")
				deletedConfigMap, err := f.ReadConfigMapFromRecycleBin(clusterConfig.RecycleBin.Path, cm)
				Expect(err).NotTo(HaveOccurred())

				By("Checking recycled configMap is the deleted configMap")
				Expect(deletedConfigMap).Should(BeEquivalentToConfigMap(cm))
			})
		})

		Context("HandleUpdate", func() {

			BeforeEach(func() {
				configMap = f.NewConfigMap()
				clusterConfig.RecycleBin.HandleUpdates = true
			})

			AfterEach(func() {
				if !f.SelfHostedOperator {
					os.RemoveAll(clusterConfig.RecycleBin.Path)
				}
			})

			It("should store old configMap in RecycleBin", func() {

				By("Creating configMap: " + configMap.Name)
				cm, err := f.CreateConfigMap(configMap)
				Expect(err).NotTo(HaveOccurred())

				By("Patching configMap: " + cm.Name)
				patchedConfigMap, _, err := core_util.PatchConfigMap(f.KubeClient, cm, func(in *core.ConfigMap) *core.ConfigMap {
					in.Data["from"] = "here"
					return in
				})
				Expect(err).NotTo(HaveOccurred())

				//give some time to recycle configMap
				time.Sleep(time.Second * 10)

				By("Checking configMap stored in RecycleBin")
				recycledConfigMap, err := f.ReadConfigMapFromRecycleBin(clusterConfig.RecycleBin.Path, cm)
				Expect(err).NotTo(HaveOccurred())

				By("Checking recycled configMap is the old configMap")
				Expect(recycledConfigMap).Should(BeEquivalentToConfigMap(cm))
				Expect(recycledConfigMap).ShouldNot(BeEquivalentToConfigMap(patchedConfigMap))
			})
		})

		Context("TTL timeout", func() {

			BeforeEach(func() {
				if f.SelfHostedOperator {
					Skip("Skipping test. Reason: In Self Hosted Operator mode Trash cleaner run in 1hour interval")
				}
				configMap = f.NewConfigMap()
				clusterConfig.RecycleBin.TTL = metav1.Duration{Duration: time.Minute}
			})

			AfterEach(func() {
				os.RemoveAll(clusterConfig.RecycleBin.Path)
			})

			It("should delete stored configMap from RecycleBin after configured TTL", func() {

				By("Creating configMap: " + configMap.Name)
				cm, err := f.CreateConfigMap(configMap)
				Expect(err).NotTo(HaveOccurred())

				By("Deleting configMap: " + cm.Name)
				err = f.DeleteConfigMap(cm.ObjectMeta)
				Expect(err).NotTo(HaveOccurred())

				//give some time to recycle configMap
				time.Sleep(time.Second * 10)

				By("Checking configMap stored in RecycleBin")
				deletedConfigMap, err := f.ReadConfigMapFromRecycleBin(clusterConfig.RecycleBin.Path, cm)
				Expect(err).NotTo(HaveOccurred())

				By("Checking recycled configMap is the deleted configMap")
				Expect(deletedConfigMap).Should(BeEquivalentToConfigMap(cm))

				By("Waiting for TTL timeout")
				time.Sleep(time.Minute * 2)

				By("Checking configMap deleted from RecycleBin")
				_, err = f.ReadConfigMapFromRecycleBin(clusterConfig.RecycleBin.Path, cm)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
