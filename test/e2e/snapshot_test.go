package e2e_test

import (
	"strconv"
	"time"

	"github.com/appscode/kubed/test/e2e/framework"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Snapshots", func() {
	var (
		f    *framework.Invocation
		cred core.Secret
	)

	BeforeEach(func() {
		f = root.Invoke()
	})
	AfterEach(func() {
		time.Sleep(60 * time.Second)
	})
	JustBeforeEach(func() {
		if missing, _ := BeZero().Match(cred); missing {
			Skip("Missing repository credential")
		}
	})

	Describe("Snapshots operations", func() {
		Context(`"Minio" backend`, func() {
			AfterEach(func() {
				f.DeleteMinioServer()
			})

			BeforeEach(func() {
				By("Creating Minio server without cacert")
				_, err = f.CreateMinioServer(false)
				Expect(err).NotTo(HaveOccurred())

				minikubeIP := "192.168.99.100"
				msvc, err := f.KubeClient.CoreV1().Services(f.Namespace()).Get("minio-service", metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				minioServiceNodePort := strconv.Itoa(int(msvc.Spec.Ports[0].NodePort))

			})
			It(`should success to perform Snapshot's operations`, performOperationOnSnapshot)

		})
	})

})
