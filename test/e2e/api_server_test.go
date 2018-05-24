package e2e_test

import (
	"os"
	"path/filepath"
	"time"

	api "github.com/appscode/kubed/apis/kubed/v1alpha1"
	"github.com/appscode/kubed/test/e2e/framework"
	. "github.com/appscode/kubed/test/e2e/matcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apps "k8s.io/api/apps/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		if !f.SelfHostedOperator {
			os.RemoveAll(filepath.Join("/tmp", "indices"))
		}
	})

	JustBeforeEach(func() {
		if f.SelfHostedOperator {
			By("Restarting kubed operator")
			err:=f.RestartKubedOperator(&clusterConfig)
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
			os.RemoveAll(filepath.Join("/tmp", "indices"))
		}
	})

	Describe("Search object", func() {

		BeforeEach(func() {
			By("Creating clusterConfiguration")
			clusterConfig = framework.APIServerClusterConfig()
		})

		Context("Search deployment by name", func() {

			BeforeEach(func() {
				deployment = f.Deployment()
			})

			AfterEach(func() {
				err := f.DeleteDeployment(deployment.ObjectMeta)
				Expect(err).NotTo(HaveOccurred())
				f.WaitUntilDeploymentTerminated(deployment.ObjectMeta)
			})

			It("SearchResult should have deployment", func() {

				By("Creating deployment: " + deployment.Name)
				_, err := f.CreateDeployment(*deployment)
				Expect(err).NotTo(HaveOccurred())
				f.WaitUntilDeploymentReady(deployment.ObjectMeta)

				// give some time for indexing
				time.Sleep(time.Second * 30)

				By("Searching deployment by name")
				result, err := f.KubedClient.KubedV1alpha1().SearchResults(deployment.Namespace).Get(deployment.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Total).Should(BeNumerically(">", 0))

				dp, err := f.KubeClient.AppsV1beta1().Deployments(deployment.Namespace).Get(deployment.Name, metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())

				By("Checking search result returns the deployment")
				Expect(result).Should(HaveObject(dp))
			})
		})
	})
})
