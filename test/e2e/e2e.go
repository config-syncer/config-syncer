package e2e

import (
	"testing"
	"time"

	"github.com/appscode/kubed/test/framework"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

const TestTimeout = 1 * time.Hour

var (
	root *framework.Framework
)

func TestE2ESuit(t *testing.T) {
	RegisterFailHandler(Fail)
	SetDefaultEventuallyTimeout(TestTimeout)

	root = framework.New()

	junitReporter := reporters.NewJUnitReporter("report.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Kubed E2e Suite", []Reporter{junitReporter})
}

var _ = BeforeSuite(func() {
	By("Ensuring Test Namespace " + root.Config.TestNamespace)
	err := root.EnsureNamespace()
	Expect(err).NotTo(HaveOccurred())
	err = root.EnsureCreatedCRDs()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	root.DeleteNamespace()
})
