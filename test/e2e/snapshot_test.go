package e2e

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	api "github.com/appscode/kubed/apis/kubed/v1alpha1"
	"github.com/appscode/kubed/test/framework"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apps "k8s.io/api/apps/v1beta1"
)

const (
	TEST_LOCAL_BACKUP_DIR = "/tmp/kubed/snapshot"
)

var _ = Describe("Snapshotter", func() {
	var (
		f             *framework.Invocation
		cred          core.Secret
		stopCh        chan struct{}
		clusterConfig api.ClusterConfig
		backend       *api.Backend
		deployment	*apps.Deployment
	)

	BeforeEach(func() {
		f = root.Invoke()
	})

	AfterEach(func() {
		close(stopCh)
	})

	JustBeforeEach(func() {
		var err error
		if missing, _ := BeZero().Match(cred); missing && backend.Local == nil {
			Skip("Missing backend credential")
		}

		if backend.Local == nil {
			err := f.CreateSecret(cred)
			Expect(err).NotTo(HaveOccurred())
		}

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

	Describe("Take Snapshot of Cluster in", func() {
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
		})

		Context(`"Local" backend`, func() {
			AfterEach(func() {
				os.RemoveAll(TEST_LOCAL_BACKUP_DIR)
			})

			BeforeEach(func() {
				err:=os.MkdirAll(TEST_LOCAL_BACKUP_DIR,0777)
				Expect(err).NotTo(HaveOccurred())

				backend = framework.NewLocalBackend(TEST_LOCAL_BACKUP_DIR)
				clusterConfig = framework.SnapshotClusterConfig(backend)
			})

			It(`should backup cluster Snapshot`, shouldTakeClusterSnapshot)
		})
	})

	Describe("Sanitize backed up object", func() {
		Context(`"Local" backend`, func() {
			AfterEach(func() {
				os.RemoveAll(TEST_LOCAL_BACKUP_DIR)
				f.DeleteDeployment(deployment.ObjectMeta)
			})

			BeforeEach(func() {
				err:=os.MkdirAll(TEST_LOCAL_BACKUP_DIR,0777)
				Expect(err).NotTo(HaveOccurred())

				backend = framework.NewLocalBackend(TEST_LOCAL_BACKUP_DIR)
				clusterConfig = framework.SnapshotClusterConfig(backend)

				deployment = f.Deployment()
				f.CreateDeployment(*deployment)
				f.WaitUntilDeploymentReady(deployment.ObjectMeta)
			})

			FIt(`should sanitize backed up deployment`, func() {
				shouldTakeClusterSnapshot()
				time.Sleep(time.Minute*1)
			})
		})
	})

})
