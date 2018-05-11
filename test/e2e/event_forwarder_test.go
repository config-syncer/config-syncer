package e2e_test

import (
	"net/http"
	"os"
	"syscall"

	api "github.com/appscode/kubed/apis/kubed/v1alpha1"
	"github.com/appscode/kubed/test/e2e/framework"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
)

var _ = Describe("Event-forwarder", func() {
	var (
		f              *framework.Invocation
		pvc            *core.PersistentVolumeClaim
		pod            *core.Pod
		clusterConfig  api.ClusterConfig
		stopCh         chan struct{}
		stopServer     chan os.Signal
		requests       []*http.Request
		notifierSecret *core.Secret
	)

	BeforeEach(func() {
		f = root.Invoke()
	})

	JustBeforeEach(func() {
		if missing, _ := BeZero().Match(notifierSecret); missing {
			Skip("Missing notifier secret")
		}

		By("Starting Operator")
		stopCh = make(chan struct{})
		err := f.RunOperator(stopCh, clusterConfig)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		close(stopCh)
		err := f.DeleteSecret(notifierSecret.ObjectMeta)
		Expect(err).NotTo(HaveOccurred())

		err = framework.ResetTestConfigFile()
		Expect(err).NotTo(HaveOccurred())

		err = f.WaitUntilSecretDeleted(notifierSecret.ObjectMeta)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Event-forwarder for Webhook Server", func() {

		BeforeEach(func() {
			requests = make([]*http.Request, 0)
			stopServer = make(chan os.Signal, 1)

			By("Starting Webhook Server")
			f.RunWebhookServer(stopServer, &requests)

			notifierSecret = f.SecretForWebhookNotifier()

			By("Creating notifier secret: " + notifierSecret.Name)
			err := f.CreateSecret(*notifierSecret)
			Expect(err).NotTo(HaveOccurred())

			By("Creating clusterConfiguration")
			clusterConfig = f.EventForwarderClusterConfig()
			clusterConfig.NotifierSecretName = notifierSecret.Name
			clusterConfig.EventForwarder.Receivers = framework.WebhookReceiver()
		})

		AfterEach(func() {
			By("Closing Webhook Server")
			stopServer <- os.Signal(syscall.SIGINT)
			defer close(stopServer)
		})

		Context("PVC add eventer", func() {

			BeforeEach(func() {
				pvc = f.NewPersistentVolumeClaim()
			})

			AfterEach(func() {
				err := f.DeletePersistentVolumeClaim(pvc.ObjectMeta)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should send notification to Webhook Server", func() {

				By("Creating PVC: " + pvc.Name)
				err := f.CreatePersistentVolumeClaim(*pvc)
				Expect(err).NotTo(HaveOccurred())

				By("Waiting for notification in Webhook Server")
				f.EventuallyNotifiedToWebhookServer(&requests, "PersistentVolumeClaim").Should(BeTrue())
			})
		})

		Context("Pod Warning Event", func() {

			BeforeEach(func() {
				pod = f.NewPod()
			})

			AfterEach(func() {
				err := f.DeletePod(pod.ObjectMeta)
				Expect(err).NotTo(HaveOccurred())

				err = f.WaitUntilPodTerminated(pod.ObjectMeta)
				Expect(err).NotTo(HaveOccurred())
			})

			It("Check warning event", func() {

				By("Creating pod: " + pod.Name)
				_, err := f.CreatePod(pod)
				Expect(err).NotTo(HaveOccurred())

				By("Waiting for notification in Webhook Server")
				f.EventuallyNotifiedToWebhookServer(&requests, pod.Name).Should(BeTrue())
			})
		})
	})
})
