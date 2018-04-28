package e2e

import (
	"fmt"
	"strconv"
	"time"

	"github.com/appscode/kubed/test/framework"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Snapshots", func() {
	var (
		f *framework.Invocation
		//cred core.Secret
	)

	BeforeEach(func() {
		f = root.Invoke()
	})
	AfterEach(func() {
		fmt.Println("==========After Each==================")
		f.StopKubedOperator()
		time.Sleep(30 * time.Second)
	})
	JustBeforeEach(func() {
		//if missing, _ := BeZero().Match(cred); missing {
		//	Skip("Missing repository credential")
		//}
		f.RunKubedOperator()
	})

	Describe("Snapshots operations", func() {
		FContext(`"Minio" backend`, func() {
			AfterEach(func() {
				f.DeleteMinioServer()
			})

			BeforeEach(func() {
				By("Creating Minio server without cacert")
				_, err := f.CreateMinioServer(true)
				Expect(err).NotTo(HaveOccurred())

				minikubeIP := "192.168.99.100"
				msvc, err := f.KubeClient.CoreV1().Services(f.Namespace()).Get("minio-service", metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				minioServiceNodePort := strconv.Itoa(int(msvc.Spec.Ports[0].NodePort))
				fmt.Println("Minio server address: https://" + minikubeIP + ":" + minioServiceNodePort)
			})
			It(`should success to perform Snapshot's operations`, func() {
				//TODO: write test
			})

		})
	})

})
