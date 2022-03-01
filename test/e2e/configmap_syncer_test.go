/*
Copyright The Config Syncer Authors.

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
	"context"
	"os"

	"kubeops.dev/config-syncer/pkg/operator"
	"kubeops.dev/config-syncer/pkg/syncer"
	"kubeops.dev/config-syncer/test/e2e/framework"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	core_util "kmodules.xyz/client-go/core/v1"
	"kmodules.xyz/client-go/meta"
)

var _ = Describe("Config-Syncer", func() {
	var (
		f           *framework.Invocation
		cfgMap      *core.ConfigMap
		nsWithLabel *core.Namespace
		config      operator.Config
	)

	BeforeEach(func() {
		f = root.Invoke()
		cfgMap = f.NewConfigMap()
		nsWithLabel = f.NewNamespaceWithLabel()
	})

	AfterEach(func() {
		f.DeleteAllConfigmaps()

		err := f.KubeClient.CoreV1().Namespaces().Delete(context.TODO(), nsWithLabel.Name, metav1.DeleteOptions{})
		if kerr.IsNotFound(err) {
			err = nil
		}
		Expect(err).NotTo(HaveOccurred())
		f.EventuallyNamespaceDeleted(nsWithLabel.Name).Should(BeTrue())
	})

	shouldSyncConfigMapToAllNamespaces := func() {
		By("Creating configMap")
		sourceCM, err := f.KubeClient.CoreV1().ConfigMaps(cfgMap.Namespace).Create(context.TODO(), cfgMap, metav1.CreateOptions{})
		Expect(err).NotTo(HaveOccurred())

		By("Checking configMap has not synced yet")
		f.EventuallyConfigMapNotSynced(sourceCM).Should(BeTrue())

		By("Adding sync annotation")
		sourceCM, _, err = core_util.PatchConfigMap(context.TODO(), f.KubeClient, sourceCM, func(obj *core.ConfigMap) *core.ConfigMap {
			metav1.SetMetaDataAnnotation(&obj.ObjectMeta, syncer.ConfigSyncKey, "")
			return obj
		}, metav1.PatchOptions{})
		Expect(err).ShouldNot(HaveOccurred())

		By("Checking configMap has synced to all namespaces")
		f.EventuallyConfigMapSynced(sourceCM).Should(BeTrue())
	}

	Describe("Across Namespaces", func() {
		BeforeEach(func() {
			config = operator.Config{}
		})

		Context("All Namespaces", func() {
			It("should sync configMap to all namespaces", shouldSyncConfigMapToAllNamespaces)
		})

		Context("New Namespace", func() {
			It("should synced configMap to new namespace", func() {
				shouldSyncConfigMapToAllNamespaces()

				By("Creating new namespace")
				err := f.CreateNamespace(nsWithLabel)
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking new namespace has the configMap")
				f.EventuallyConfigMapSyncedToNamespace(cfgMap, nsWithLabel.Name).Should(BeTrue())
			})
		})

		Context("Remove Sync Annotation", func() {
			It("should delete synced configMaps", func() {
				shouldSyncConfigMapToAllNamespaces()

				By("Removing sync annotation")
				source, err := f.KubeClient.CoreV1().ConfigMaps(cfgMap.Namespace).Get(context.TODO(), cfgMap.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())

				_, _, err = core_util.PatchConfigMap(context.TODO(), f.KubeClient, source, func(obj *core.ConfigMap) *core.ConfigMap {
					obj.Annotations = meta.RemoveKey(obj.Annotations, syncer.ConfigSyncKey)
					return obj
				}, metav1.PatchOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking synced configMaps has been deleted")
				f.EventuallySyncedConfigMapsDeleted(source)
			})
		})

		Context("Source Update ConfigMap Data", func() {
			It("should update synced configMaps", func() {
				shouldSyncConfigMapToAllNamespaces()

				By("Updating source configMap")
				source, err := f.KubeClient.CoreV1().ConfigMaps(cfgMap.Namespace).Get(context.TODO(), cfgMap.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())

				source, _, err = core_util.PatchConfigMap(context.TODO(), f.KubeClient, source, func(obj *core.ConfigMap) *core.ConfigMap {
					if obj.Data == nil {
						obj.Data = map[string]string{}
					}
					obj.Data["data"] = "test"
					return obj
				}, metav1.PatchOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking synced configMaps has been updated")
				f.EventuallySyncedConfigMapsUpdated(source).Should(BeTrue())
			})
		})

		Context("Source Update ConfigMap BinaryData", func() {
			It("should update synced configMaps", func() {
				shouldSyncConfigMapToAllNamespaces()

				By("Updating source configMap")
				source, err := f.KubeClient.CoreV1().ConfigMaps(cfgMap.Namespace).Get(context.TODO(), cfgMap.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())

				source, _, err = core_util.PatchConfigMap(context.TODO(), f.KubeClient, source, func(obj *core.ConfigMap) *core.ConfigMap {
					if obj.BinaryData == nil {
						obj.BinaryData = map[string][]byte{}
					}
					obj.BinaryData["data"] = []byte("test")
					return obj
				}, metav1.PatchOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking synced configMaps has been updated")
				f.EventuallySyncedConfigMapsUpdated(source).Should(BeTrue())
			})
		})

		Context("Backward Compatibility", func() {
			It("should sync configMap to all namespaces", func() {
				By("Creating configMap")
				source, err := f.CreateConfigMap(cfgMap)
				Expect(err).NotTo(HaveOccurred())

				By("Checking configMap has not synced yet")
				f.EventuallyConfigMapNotSynced(source).Should(BeTrue())

				By("Adding sync=true annotation")
				source, _, err = core_util.PatchConfigMap(context.TODO(), f.KubeClient, source, func(obj *core.ConfigMap) *core.ConfigMap {
					metav1.SetMetaDataAnnotation(&obj.ObjectMeta, syncer.ConfigSyncKey, "true")
					return obj
				}, metav1.PatchOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking configMap has synced to all namespaces")
				f.EventuallyConfigMapSynced(source).Should(BeTrue())
			})
		})

		Context("Namespace Selector", func() {
			It("should add configMap to selected namespaces", func() {
				shouldSyncConfigMapToAllNamespaces()

				By("Adding selector annotation")
				source, err := f.KubeClient.CoreV1().ConfigMaps(cfgMap.Namespace).Get(context.TODO(), cfgMap.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())

				source, _, err = core_util.PatchConfigMap(context.TODO(), f.KubeClient, source, func(obj *core.ConfigMap) *core.ConfigMap {
					metav1.SetMetaDataAnnotation(&obj.ObjectMeta, syncer.ConfigSyncKey, "app="+f.App())
					return obj
				}, metav1.PatchOptions{})
				Expect(err).NotTo(HaveOccurred())

				By("Checking configMap has not synced to other namespaces")
				f.EventuallyConfigMapNotSynced(source).Should(BeTrue())

				By("Creating new namespace with label")
				err = f.CreateNamespace(nsWithLabel)
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking configmap synced to new namespace")
				f.EventuallyConfigMapSyncedToNamespace(source, nsWithLabel.Name)

				By("Changing selector annotation")
				_, _, err = core_util.PatchConfigMap(context.TODO(), f.KubeClient, source, func(obj *core.ConfigMap) *core.ConfigMap {
					metav1.SetMetaDataAnnotation(&obj.ObjectMeta, syncer.ConfigSyncKey, "app=do-not-match")
					return obj
				}, metav1.PatchOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking synced configMap has been deleted")
				f.EventuallySyncedConfigMapsDeleted(source)

				By("Removing selector annotation")
				source, err = f.KubeClient.CoreV1().ConfigMaps(source.Namespace).Get(context.TODO(), source.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())

				source, _, err = core_util.PatchConfigMap(context.TODO(), f.KubeClient, source, func(obj *core.ConfigMap) *core.ConfigMap {
					metav1.SetMetaDataAnnotation(&obj.ObjectMeta, syncer.ConfigSyncKey, "")
					return obj
				}, metav1.PatchOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking configMap synced to all namespaces")
				f.EventuallyConfigMapSynced(source).Should(BeTrue())
			})
		})

		Context("Source Deleted", func() {
			It("should delete synced configMaps", func() {
				shouldSyncConfigMapToAllNamespaces()

				By("Deleting source configMap")
				source, err := f.KubeClient.CoreV1().ConfigMaps(cfgMap.Namespace).Get(context.TODO(), cfgMap.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())

				err = f.DeleteConfigMap(source.ObjectMeta)
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking synced configMaps has been deleted")
				f.EventuallySyncedConfigMapsDeleted(source).Should(BeTrue())
			})
		})

		Context("Source Namespace Deleted", func() {
			var sourceNamespace *core.Namespace

			BeforeEach(func() {
				sourceNamespace = f.NewNamespace("source")
				cfgMap.Namespace = sourceNamespace.Name
			})

			It("should delete synced configMaps", func() {
				By("Creating source namespace")
				err := f.CreateNamespace(sourceNamespace)
				Expect(err).NotTo(HaveOccurred())

				shouldSyncConfigMapToAllNamespaces()

				By("Deleting source namespace")
				source, err := f.KubeClient.CoreV1().ConfigMaps(cfgMap.Namespace).Get(context.TODO(), cfgMap.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())

				err = f.DeleteNamespace(sourceNamespace.Name)
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking synced configMaps has been deleted")
				f.EventuallySyncedConfigMapsDeleted(source).Should(BeTrue())
			})
		})
	})

	Describe("Across Cluster", func() {
		Context("ConfigMap Context Syncer Test", func() {
			var (
				kubeConfigPath = "/home/dipta/all/config-syncer-test/kubeconfig"
				ctx            = "gke_tigerworks-kube_us-central1-f_kite"
			)

			BeforeEach(func() {
				config = operator.Config{}
				config.ClusterName = "minikube"
				config.KubeConfigFile = kubeConfigPath

				if _, err := os.Stat(kubeConfigPath); err != nil {
					Skip(`"config" file not found on` + kubeConfigPath)
				}

				By("Creating namespace for context")
				f.EnsureNamespaceForContext(kubeConfigPath, ctx)
			})

			AfterEach(func() {
				By("Deleting namespaces for contexts")
				f.DeleteNamespaceForContext(kubeConfigPath, ctx)
			})

			XIt("Should add configmap to contexts", func() {
				By("Creating source ns in remote cluster")
				f.EnsureNamespaceForContext(kubeConfigPath, ctx)

				By("Creating configmap")
				cfgMap, err := f.KubeClient.CoreV1().ConfigMaps(cfgMap.Namespace).Create(context.TODO(), cfgMap, metav1.CreateOptions{})
				Expect(err).NotTo(HaveOccurred())

				By("Adding sync annotation")
				cfgMap, _, err = core_util.PatchConfigMap(context.TODO(), f.KubeClient, cfgMap, func(obj *core.ConfigMap) *core.ConfigMap {
					metav1.SetMetaDataAnnotation(&obj.ObjectMeta, syncer.ConfigSyncContexts, ctx)
					return obj
				}, metav1.PatchOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking configmap added to contexts")
				f.EventuallyNumOfConfigmapsForContext(kubeConfigPath, ctx).Should(BeNumerically("==", 1))

				By("Removing sync annotation")
				_, _, err = core_util.PatchConfigMap(context.TODO(), f.KubeClient, cfgMap, func(obj *core.ConfigMap) *core.ConfigMap {
					obj.Annotations = meta.RemoveKey(obj.Annotations, syncer.ConfigSyncContexts)
					return obj
				}, metav1.PatchOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				By("Checking configmap removed from contexts")
				f.EventuallyNumOfConfigmapsForContext(kubeConfigPath, ctx).Should(BeNumerically("==", 0))
			})
		})
	})
})
