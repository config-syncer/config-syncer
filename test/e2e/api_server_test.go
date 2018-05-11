package e2e_test

import (
	api "github.com/appscode/kubed/apis/kubed/v1alpha1"
	"github.com/appscode/kubed/test/e2e/framework"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/the-redback/go-oneliners"
	apps "k8s.io/api/apps/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

var _ = Describe("API server", func() {
	var (
		f             *framework.Invocation
		deployment    *apps.Deployment
		clusterConfig api.ClusterConfig
		stopCh        chan struct{}
	)

	BeforeEach(func() {
		f = root.Invoke()
	})

	JustBeforeEach(func() {
		By("Starting Operator")
		stopCh = make(chan struct{})
		err := f.RunOperator(stopCh, clusterConfig)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		close(stopCh)

		err := framework.ResetTestConfigFile()
		Expect(err).NotTo(HaveOccurred())
	})

	FDescribe("Search object", func() {

		BeforeEach(func() {
			By("Creating clusterConfiguration")
			clusterConfig = framework.APIServerClusterConfig()
		})

		Context("Search deployment by name", func() {

			BeforeEach(func() {
				deployment = f.Deployment()
				By("Creating deployment: " + deployment.Name)
				_, err := f.CreateDeployment(*deployment)
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				err := f.DeleteDeployment(deployment.ObjectMeta)
				Expect(err).NotTo(HaveOccurred())
				f.WaitUntilDeploymentTerminated(deployment.ObjectMeta)
			})

			It("SearchResult should have deployment", func() {

				time.Sleep(time.Minute*3)
				By("Searching deployment by name")
				result, err := f.KubedClient.KubedV1alpha1().SearchResults(deployment.Namespace).Get(deployment.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				oneliners.PrettyJson(result, "SearchResult")

				//TODO: do rest of the test
			})
		})
	})
})
