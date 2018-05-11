package e2e

import (
	"os"
	"testing"
	"time"

	"github.com/appscode/kubed/test/framework"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	rbac "k8s.io/api/rbac/v1"
	genericapiserver "k8s.io/apiserver/pkg/server"
)

const TestTimeout = 3 * time.Minute

var (
	root            *framework.Framework
	userRule        *rbac.ClusterRole
	userRoleBinding *rbac.ClusterRoleBinding
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

	By("Creating initial kubed configuration file")
	err = framework.APIServerClusterConfig().Save(framework.KubedTestConfigFileDir)
	Expect(err).NotTo(HaveOccurred())

	By("Registering API service")
	err = root.Invoke().RegisterAPIService()
	Expect(err).NotTo(HaveOccurred())

	root.KubedServer, err = root.NewTestKubedServer()
	Expect(err).NotTo(HaveOccurred())

	//By("Creating UserRole")
	//userRule = root.CreateUserRole()
	//_, err = root.KubeClient.RbacV1().ClusterRoles().Create(userRule)
	//Expect(err).NotTo(HaveOccurred())

	//By("Binding user role to " + framework.USER_ANONYMOUS)
	//userRoleBinding = root.UserRoleBinding()
	//_, err = root.KubeClient.RbacV1().ClusterRoleBindings().Create(userRoleBinding)
	//Expect(err).NotTo(HaveOccurred())

	By("Starting API server")
	stopCh := genericapiserver.SetupSignalHandler()
	go root.KubedServer.GenericAPIServer.PrepareRun().Run(stopCh)

	By("Waiting for API server to be ready")
	root.EventuallyAPIServerReady().Should(Succeed())
	time.Sleep(time.Second * 5)
})

var _ = AfterSuite(func() {
	By("Cleaning API service stuff")
	root.Invoke().DeleteAPIService()
	//root.DeleteClusterRole(userRule.ObjectMeta)
	//root.DeleteClusterRoleBinding(userRoleBinding.ObjectMeta)
	root.DeleteNamespace()
	os.Remove(framework.KubedTestConfigFileDir)
})
