package e2e

import (
	"fmt"
	"net"
	"strconv"

	api "github.com/appscode/kubed/apis/kubed/v1alpha1"
	"github.com/appscode/kubed/test/framework"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Snapshots", func() {
	var (
		f             *framework.Invocation
		cred          core.Secret
		stopCh        chan struct{}
		clusterConfig api.ClusterConfig
		backend       *api.Backend
	)

	BeforeEach(func() {
		f = root.Invoke()
	})

	AfterEach(func() {
		close(stopCh)
	})

	JustBeforeEach(func() {
		if missing, _ := BeZero().Match(cred); missing {
			Skip("Missing backend credential")
		}

		err := f.CreateSecret(cred)
		Expect(err).NotTo(HaveOccurred())

		operatorConfig := f.NewTestOperatorConfig()
		f.KubedServer.Operator, err = operatorConfig.New()
		Expect(err).NotTo(HaveOccurred())

		f.KubedServer.Operator.ClusterConfig = clusterConfig
		f.KubedServer.Operator.OperatorNamespace = f.Namespace()

		err = f.CreateBucketIfNotExist(clusterConfig.Snapshotter.Backend)
		Expect(err).NotTo(HaveOccurred())

		stopCh = make(chan struct{})
		go f.KubedServer.Operator.Run(stopCh, true)
	})

	shouldTakeClusterSnapshot := func() {
		f.EventuallyBackupSnapshot(*backend).ShouldNot(BeEmpty())
	}

	Describe("Snapshots operations", func() {
		Context(`"Minio" backend`, func() {
			AfterEach(func() {
				f.DeleteMinioServer()
			})

			BeforeEach(func() {
				minikubeIP := net.IP{192, 168, 99, 100}

				By("Creating Minio server with cacert")
				_, err := f.CreateMinioServer(true, []net.IP{minikubeIP})
				Expect(err).NotTo(HaveOccurred())

				msvc, err := f.KubeClient.CoreV1().Services(f.Namespace()).Get("minio-service", metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				minioServiceNodePort := strconv.Itoa(int(msvc.Spec.Ports[0].NodePort))
				minioEndpoint := fmt.Sprintf("https://" + minikubeIP.String() + ":" + minioServiceNodePort)

				cred = f.SecretForMinioBackend(true)

				backend = framework.NewMinioBackend("kubed-test", "demo", minioEndpoint, cred.Name)
				clusterConfig = framework.SnapshotClusterConfig(backend)
			})
			It(`should backup cluster Snapshot`, shouldTakeClusterSnapshot)

			It(`should backup cluster Snapshot`, shouldTakeClusterSnapshot)

		})
	})

})
