package e2e

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/ginkgo/reporters"
	"testing"
	"time"
	"github.com/appscode/kubed/test/framework"
	"github.com/appscode/kubed/pkg/operator"
)

const TestTimeout  = 2 * time.Hour

var (
	root *framework.Framework
)

func TestE2ESuit(t *testing.T)  {
	RegisterFailHandler(Fail)
	SetDefaultEventuallyTimeout(TestTimeout)

	root = framework.New()

	junitReporter := reporters.NewJUnitReporter("report.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Kubed E2e Suite", []Reporter{junitReporter})
}

var _ = BeforeSuite(func() {
	op := &operator.Operator{
		KubeClient:        root.KubeClient,
		VoyagerClient:     root.KubedOperator.VoyagerClient,
		SearchlightClient: root.KubedOperator.SearchlightClient,
		StashClient:       root.KubedOperator.StashClient,
		KubeDBClient:      root.KubedOperator.KubeDBClient,
		Opt: operator.Options{
			KubeConfig: root.Config.KubeConfig,
			ConfigPath: "/srv/kubed/config.yaml",
		},
	}

	By("Ensuring Test Namespace " + root.Config.TestNamespace)
	err := root.EnsureNamespace()
	Expect(err).NotTo(HaveOccurred())

	err = op.Setup()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	root.DeleteNamespace()
})