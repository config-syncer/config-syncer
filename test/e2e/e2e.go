package e2e

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/appscode/kubed/test/framework"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

const TestTimeout = 3 * time.Minute

var (
	root *framework.Framework
)

func RunE2ETestSuit(t *testing.T) {
	RegisterFailHandler(Fail)
	SetDefaultEventuallyTimeout(TestTimeout)

	root = framework.New()

	junitReporter := reporters.NewJUnitReporter("report.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Kubed E2E Suite", []Reporter{junitReporter})
}

var _ = BeforeSuite(func() {
	By("Ensuring Test Namespace " + root.Config.TestNamespace)
	err := root.EnsureNamespace()
	Expect(err).NotTo(HaveOccurred())

	By("Creating CRDs")
	err = root.EnsureCreatedCRDs()
	Expect(err).NotTo(HaveOccurred())

	By("Creating kubed configuration file dir")
	err = os.MkdirAll(filepath.Dir(framework.KubedTestConfigFileDir), 0755)
	Expect(err).NotTo(HaveOccurred())

	By("Registering API service")
	err = root.Invoke().RegisterAPIService()
	Expect(err).NotTo(HaveOccurred())

	root.KubedServer, err = root.NewTestKubedServer()
	Expect(err).NotTo(HaveOccurred())

	By("Starting API server")
	stopCh := genericapiserver.SetupSignalHandler()
	go root.KubedServer.GenericAPIServer.PrepareRun().Run(stopCh)

	By("Waiting for API server to be ready")
	root.EventuallyAPIServerReady().Should(Succeed())
})

var _ = AfterSuite(func() {
	By("Cleaning API service stuff")
	root.Invoke().DeleteAPIService()

	root.DeleteNamespace()
	//os.RemoveAll(framework.KubedTestConfigFileDir)
})
