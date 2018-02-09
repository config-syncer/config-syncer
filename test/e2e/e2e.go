package e2e

import (
	"testing"
	"time"

	"github.com/appscode/kubed/test/framework"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
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
	err = root.EnsureCreatedCRDs()
	Expect(err).NotTo(HaveOccurred())

	//// configure and run operator
	//opt := server.Options{
	//	KubeConfig:        filepath.Join(homedir.HomeDir(), ".kube/config"),
	//	ConfigPath:        "config.yaml",
	//	APIAddress:        ":8080",
	//	WebAddress:        ":56790",
	//	ScratchDir:        "/tmp/kubed",
	//	OperatorNamespace: root.Namespace(),
	//	ResyncPeriod:      5 * time.Minute,
	//}
	//
	//By("Running kubed operator")
	//go cmds.Run(opt)
})

var _ = AfterSuite(func() {
	root.DeleteNamespace()
})
