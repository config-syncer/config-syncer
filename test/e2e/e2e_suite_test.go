package e2e_test

import (
	"os"
	"testing"
	"time"

	prom_util "github.com/appscode/kube-mon/prometheus/v1"
	"github.com/appscode/kubed/test/e2e/framework"
	"github.com/appscode/kutil/tools/clientcmd"
	searchlightcheme "github.com/appscode/searchlight/client/clientset/versioned/scheme"
	stashscheme "github.com/appscode/stash/client/clientset/versioned/scheme"
	voyagerscheme "github.com/appscode/voyager/client/clientset/versioned/scheme"
	kubedbscheme "github.com/kubedb/apimachinery/client/clientset/versioned/scheme"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	rbac "k8s.io/api/rbac/v1"
	genericapiserver "k8s.io/apiserver/pkg/server"
	logs "github.com/appscode/go/log/golog"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
)

const TestTimeout = 3 * time.Minute

var (
	root            *framework.Framework
	userRule        *rbac.ClusterRole
	userRoleBinding *rbac.ClusterRoleBinding
)

func TestE2E(t *testing.T) {
	logs.InitLogs()
	RegisterFailHandler(Fail)
	SetDefaultEventuallyTimeout(TestTimeout)
	junitReporter := reporters.NewJUnitReporter("report.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Kubed E2E Suite", []Reporter{junitReporter})
}

var _ = BeforeSuite(func() {
	voyagerscheme.AddToScheme(clientsetscheme.Scheme)
	searchlightcheme.AddToScheme(clientsetscheme.Scheme)
	stashscheme.AddToScheme(clientsetscheme.Scheme)
	kubedbscheme.AddToScheme(clientsetscheme.Scheme)
	prom_util.AddToScheme(clientsetscheme.Scheme)

	clientConfig, err := clientcmd.BuildConfigFromContext(options.KubeConfig, options.KubeContext)
	Expect(err).NotTo(HaveOccurred())

	root = framework.New(clientConfig)

	By("Using Namespace " + root.Namespace())
	err = root.EnsureNamespace()
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

	root.KubedServer, err = root.NewTestKubedServer(options.KubeConfig)
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
