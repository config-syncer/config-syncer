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
	"testing"
	"time"

	"kubeops.dev/config-syncer/test/e2e/framework"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"gomodules.xyz/logs"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"kmodules.xyz/client-go/tools/clientcmd"
)

const TestTimeout = 3 * time.Minute

var root *framework.Framework

func TestE2E(t *testing.T) {
	logs.InitLogs()
	RegisterFailHandler(Fail)
	SetDefaultEventuallyTimeout(TestTimeout)
	junitReporter := reporters.NewJUnitReporter("junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Config Syncer E2E Suite", []Reporter{junitReporter})
}

var _ = BeforeSuite(func() {
	clientConfig, err := clientcmd.BuildConfigFromContext(options.KubeConfig, options.KubeContext)
	Expect(err).NotTo(HaveOccurred())

	root = framework.New(clientConfig)
	root.KubeConfigPath = options.KubeConfig

	By("Using Namespace " + root.Namespace())
	err = root.EnsureNamespace()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	utilruntime.Must(root.DeleteNamespace(root.Namespace()))
})
