package e2e

import (
	"os"

	"github.com/appscode/kubed/test/framework"
	"github.com/appscode/kutil/meta"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/api/apps/v1beta1"
	batch_v1 "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	ext_v1 "k8s.io/api/extensions/v1beta1"
)

var _ = Describe("Extract Docker Label", func() {
	var (
		f *framework.Invocation

		labels          map[string]string
		name, namespace string
		containers1     []core.Container
		containers2     []core.Container
		containers3     []core.Container
		data            string
		skip            bool

		secret                    *core.Secret
		service, svc              *core.Service
		deployment, deploy        *v1beta1.Deployment
		replicationcontroller, rc *core.ReplicationController
		replicaset, rs            *ext_v1.ReplicaSet
		daemonset, ds             *ext_v1.DaemonSet
		job, newJob               *batch_v1.Job
		statefulset, sts          *v1beta1.StatefulSet
	)

	BeforeEach(func() {
		f = root.Invoke()
		name = f.App()
		namespace = "kube-system"
		labels = map[string]string{
			"app": f.App(),
		}
		data = ""
		skip = false

		if val, ok := os.LookupEnv("DOCKER_CFG"); !ok {
			skip = true
		} else {
			data = val
			secret = f.NewSecret(name, namespace, data, labels)
		}

		containers1 = []core.Container{
			{
				Name:  "labels",
				Image: "shudipta/labels",
				Ports: []core.ContainerPort{
					{
						ContainerPort: 80,
					},
				},
			},
		}
		containers2 = []core.Container{
			{
				Name:  "book-server",
				Image: "shudipta/book_server:v1",
				Ports: []core.ContainerPort{
					{
						ContainerPort: 10000,
					},
				},
			},
		}
		containers3 = []core.Container{
			{
				Name:  "guard",
				Image: "nightfury1204/guard:azure",
				Ports: []core.ContainerPort{
					{
						ContainerPort: 10000,
					},
				},
			},
		}
	})

	Describe("Adding Annotaions in Deployment", func() {
		JustBeforeEach(func() {
			By("Creating secret")
			_, err := root.KubeClient.CoreV1().Secrets(secret.Namespace).Create(secret)
			Expect(err).NotTo(HaveOccurred())
			//f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))

			By("Creating deployment")
			deploy, err = root.KubeClient.AppsV1beta1().Deployments(deployment.Namespace).Create(deployment)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			f.DeleteAllSecrets()
			f.DeleteAllDeployments()
		})

		Context("When docker image contains labels", func() {
			BeforeEach(func() {
				deployment = f.NewDeployment(name, namespace, labels, containers1)
			})

			It("Should add annotations", func() {
				if skip {
					Skip("environment var \"DOCKER_CFG\" not found")
				}

				imageAnnotations := map[string]string{
					"docker.com/labels-git-commit": "unkown",
				}

				By("\"git-commit\" annotations should be present")
				f.EventuallyAnnotationsFromDeployment(deploy.Name, deploy.Namespace, "docker.com/").
					Should(Equal(imageAnnotations))

				By("\"hello\" annotation shouldn't be present")
				annotations := f.AnnotaionsFromDeployment(deploy.Name, deploy.Namespace, "docker.com/")
				value, err := meta.GetStringValue(annotations, "docker.com/hi-hello")
				Expect(err).To(HaveOccurred())
				Expect(value).To(Equal(""))
			})
		})

		Context("When docker image doesn't contain any labels", func() {
			BeforeEach(func() {
				deployment = f.NewDeployment(name, namespace, labels, containers2)
			})

			It("Shouldn't add annotations", func() {
				if skip {
					Skip("environment var \"DOCKER_CFG\" not found")
				}

				imageAnnotations := map[string]string{}
				//time.Sleep(time.Second * 10)
				By("\"git-commit\" annotations should be present")
				f.EventuallyAnnotationsFromDeployment(deploy.Name, deploy.Namespace, "docker.com/").
					Should(Equal(imageAnnotations))
				//annotations := f.AnnotaionsFromDeployment(deploy.Name, deploy.Namespace, "docker.com/")
				//log.Infoln("annotations =", annotations)
				//Expect(annotations).To(Equal(imageAnnotations))
			})
		})

		Context("When docker images aren't found with the provided secrets", func() {
			BeforeEach(func() {
				deployment = f.NewDeployment(name, namespace, labels, containers3)
			})

			It("Shouldn't found images", func() {
				if skip {
					Skip("environment var \"DOCKER_CFG\" not found")
				}

				imageAnnotations := map[string]string{}
				By("\"git-commit\" annotations should be present")
				f.EventuallyAnnotationsFromDeployment(deploy.Name, deploy.Namespace, "docker.com/").
					Should(Equal(imageAnnotations))
			})
		})
	})

	Describe("Adding Annotaions in ReplicationController", func() {
		JustBeforeEach(func() {
			By("Creating secret")
			_, err := root.KubeClient.CoreV1().Secrets(secret.Namespace).Create(secret)
			Expect(err).NotTo(HaveOccurred())
			//f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))

			By("Creating replicationcontroller")
			rc, err = root.KubeClient.CoreV1().
				ReplicationControllers(replicationcontroller.Namespace).
				Create(replicationcontroller)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			f.DeleteAllSecrets()
			f.DeleteAllReplicationControllers()
		})

		Context("When docker image contains labels", func() {
			BeforeEach(func() {
				replicationcontroller = f.NewReplicationController(name, namespace, labels, containers1)
			})

			It("Should add annotations", func() {
				if skip {
					Skip("environment var \"DOCKER_CFG\" not found")
				}

				imageAnnotations := map[string]string{
					"docker.com/labels-git-commit": "unkown",
				}

				By("\"git-commit\" annotations should be present")
				f.EventuallyAnnotationsFromReplicationController(rc.Name, rc.Namespace, "docker.com/").
					Should(Equal(imageAnnotations))

				By("\"hello\" annotation shouldn't be present")
				annotations := f.AnnotaionsFromReplicationController(rc.Name, rc.Namespace, "docker.com/")
				value, err := meta.GetStringValue(annotations, "docker.com/hi-hello")
				Expect(err).To(HaveOccurred())
				Expect(value).To(Equal(""))
			})
		})

		Context("When docker image doesn't contain any labels", func() {
			BeforeEach(func() {
				replicationcontroller = f.NewReplicationController(name, namespace, labels, containers2)
			})

			It("Shouldn't add annotations", func() {
				if skip {
					Skip("environment var \"DOCKER_CFG\" not found")
				}

				imageAnnotations := map[string]string{}
				By("\"git-commit\" annotations should be present")
				f.EventuallyAnnotationsFromReplicationController(rc.Name, rc.Namespace, "docker.com/").
					Should(Equal(imageAnnotations))
			})
		})

		Context("When docker images aren't found with the provided secrets", func() {
			BeforeEach(func() {
				replicationcontroller = f.NewReplicationController(name, namespace, labels, containers3)
			})

			It("Shouldn't found images", func() {
				if skip {
					Skip("environment var \"DOCKER_CFG\" not found")
				}

				imageAnnotations := map[string]string{}
				By("\"git-commit\" annotations should be present")
				f.EventuallyAnnotationsFromReplicationController(rc.Name, rc.Namespace, "docker.com/").
					Should(Equal(imageAnnotations))
			})
		})
	})

	Describe("Adding Annotaions in ReplicaSet", func() {
		JustBeforeEach(func() {
			By("Creating secret")
			_, err := root.KubeClient.CoreV1().Secrets(secret.Namespace).Create(secret)
			Expect(err).NotTo(HaveOccurred())
			//f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))

			By("Creating replicaset")
			rs, err = root.KubeClient.ExtensionsV1beta1().ReplicaSets(replicaset.Namespace).Create(replicaset)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			f.DeleteAllSecrets()
			f.DeleteAllReplicasets()
		})

		Context("When docker image contains labels", func() {
			BeforeEach(func() {
				replicaset = f.NewReplicaSet(name, namespace, labels, containers1)
			})

			It("Should add annotations", func() {
				if skip {
					Skip("environment var \"DOCKER_CFG\" not found")
				}

				imageAnnotations := map[string]string{
					"docker.com/labels-git-commit": "unkown",
				}

				By("\"git-commit\" annotations should be present")
				f.EventuallyAnnotationsFromReplicaSet(rs.Name, rs.Namespace, "docker.com/").
					Should(Equal(imageAnnotations))

				By("\"hello\" annotation shouldn't be present")
				annotations := f.AnnotaionsFromReplicaSet(rs.Name, rs.Namespace, "docker.com/")
				value, err := meta.GetStringValue(annotations, "docker.com/hi-hello")
				Expect(err).To(HaveOccurred())
				Expect(value).To(Equal(""))
			})
		})

		Context("When docker image doesn't contain any labels", func() {
			BeforeEach(func() {
				replicaset = f.NewReplicaSet(name, namespace, labels, containers2)
			})

			It("Shouldn't add annotations", func() {
				if skip {
					Skip("environment var \"DOCKER_CFG\" not found")
				}

				imageAnnotations := map[string]string{}
				By("\"git-commit\" annotations should be present")
				f.EventuallyAnnotationsFromReplicaSet(rs.Name, rs.Namespace, "docker.com/").
					Should(Equal(imageAnnotations))
			})
		})

		Context("When docker images aren't found with the provided secrets", func() {
			BeforeEach(func() {
				replicaset = f.NewReplicaSet(name, namespace, labels, containers3)
			})

			It("Shouldn't found images", func() {
				if skip {
					Skip("environment var \"DOCKER_CFG\" not found")
				}

				imageAnnotations := map[string]string{}
				By("\"git-commit\" annotations should be present")
				f.EventuallyAnnotationsFromReplicaSet(rs.Name, rs.Namespace, "docker.com/").
					Should(Equal(imageAnnotations))
			})
		})
	})

	Describe("Adding Annotaions in DaemonSet", func() {
		JustBeforeEach(func() {
			By("Creating secret")
			_, err := root.KubeClient.CoreV1().Secrets(secret.Namespace).Create(secret)
			Expect(err).NotTo(HaveOccurred())
			//f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))

			By("Creating daemonset")
			ds, err = root.KubeClient.ExtensionsV1beta1().DaemonSets(daemonset.Namespace).Create(daemonset)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			f.DeleteAllSecrets()
			f.DeleteAllDaemonSet()
		})

		Context("When docker image contains labels", func() {
			BeforeEach(func() {
				daemonset = f.NewDaemonSet(name, namespace, labels, containers1)
			})

			It("Should add annotations", func() {
				if skip {
					Skip("environment var \"DOCKER_CFG\" not found")
				}

				imageAnnotations := map[string]string{
					"docker.com/labels-git-commit": "unkown",
				}

				By("\"git-commit\" annotations should be present")
				f.EventuallyAnnotationsFromDaemonSet(ds.Name, ds.Namespace, "docker.com/").
					Should(Equal(imageAnnotations))

				By("\"hello\" annotation shouldn't be present")
				annotations := f.AnnotaionsFromDaemonSet(ds.Name, ds.Namespace, "docker.com/")
				value, err := meta.GetStringValue(annotations, "docker.com/hi-hello")
				Expect(err).To(HaveOccurred())
				Expect(value).To(Equal(""))
			})
		})

		Context("When docker image doesn't contain any labels", func() {
			BeforeEach(func() {
				daemonset = f.NewDaemonSet(name, namespace, labels, containers2)
			})

			It("Shouldn't add annotations", func() {
				if skip {
					Skip("environment var \"DOCKER_CFG\" not found")
				}

				imageAnnotations := map[string]string{}
				By("\"git-commit\" annotations should be present")
				f.EventuallyAnnotationsFromDaemonSet(ds.Name, ds.Namespace, "docker.com/").
					Should(Equal(imageAnnotations))
			})
		})

		Context("When docker images aren't found with the provided secrets", func() {
			BeforeEach(func() {
				daemonset = f.NewDaemonSet(name, namespace, labels, containers3)
			})

			It("Shouldn't found images", func() {
				if skip {
					Skip("environment var \"DOCKER_CFG\" not found")
				}

				imageAnnotations := map[string]string{}
				By("\"git-commit\" annotations should be present")
				f.EventuallyAnnotationsFromDaemonSet(ds.Name, ds.Namespace, "docker.com/").
					Should(Equal(imageAnnotations))
			})
		})
	})

	Describe("Adding Annotaions in Job", func() {
		JustBeforeEach(func() {
			By("Creating secret")
			_, err := root.KubeClient.CoreV1().Secrets(secret.Namespace).Create(secret)
			Expect(err).NotTo(HaveOccurred())
			//f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))

			By("Creating job")
			newJob, err = root.KubeClient.BatchV1().Jobs(job.Namespace).Create(job)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			f.DeleteAllSecrets()
			f.DeleteAllJobs()
		})

		Context("When docker image contains labels", func() {
			BeforeEach(func() {
				job = f.NewJob(name, namespace, labels, containers1)
			})

			It("Should add annotations", func() {
				if skip {
					Skip("environment var \"DOCKER_CFG\" not found")
				}

				imageAnnotations := map[string]string{
					"docker.com/labels-git-commit": "unkown",
				}

				By("\"git-commit\" annotations should be present")
				f.EventuallyAnnotationsFromJob(newJob.Name, newJob.Namespace, "docker.com/").
					Should(Equal(imageAnnotations))

				By("\"hello\" annotation shouldn't be present")
				annotations := f.AnnotaionsFromJob(newJob.Name, newJob.Namespace, "docker.com/")
				value, err := meta.GetStringValue(annotations, "docker.com/hi-hello")
				Expect(err).To(HaveOccurred())
				Expect(value).To(Equal(""))
			})
		})

		Context("When docker image doesn't contain any labels", func() {
			BeforeEach(func() {
				job = f.NewJob(name, namespace, labels, containers2)
			})

			It("Shouldn't add annotations", func() {
				if skip {
					Skip("environment var \"DOCKER_CFG\" not found")
				}

				imageAnnotations := map[string]string{}
				By("\"git-commit\" annotations should be present")
				f.EventuallyAnnotationsFromJob(newJob.Name, newJob.Namespace, "docker.com/").
					Should(Equal(imageAnnotations))
			})
		})

		Context("When docker images aren't found with the provided secrets", func() {
			BeforeEach(func() {
				job = f.NewJob(name, namespace, labels, containers3)
			})

			It("Shouldn't found images", func() {
				if skip {
					Skip("environment var \"DOCKER_CFG\" not found")
				}

				imageAnnotations := map[string]string{}
				By("\"git-commit\" annotations should be present")
				f.EventuallyAnnotationsFromJob(newJob.Name, newJob.Namespace, "docker.com/").
					Should(Equal(imageAnnotations))
			})
		})
	})

	Describe("Adding Annotaions in StatefulSet", func() {
		JustBeforeEach(func() {
			By("Creating secret")
			_, err := root.KubeClient.CoreV1().Secrets(secret.Namespace).Create(secret)
			Expect(err).NotTo(HaveOccurred())
			//f.EventuallyNumOfSecrets(f.Namespace()).Should(BeNumerically("==", 1))

			By("Creating service")
			svc, err = root.KubeClient.CoreV1().Services(service.Namespace).Create(service)
			Expect(err).NotTo(HaveOccurred())

			By("Creating statefulset")
			sts, err = root.KubeClient.AppsV1beta1().StatefulSets(statefulset.Namespace).Create(statefulset)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			f.DeleteAllSecrets()
			f.DeleteAllServices()
			f.DeleteAllStatefulSets()
		})

		Context("When docker image contains labels", func() {
			BeforeEach(func() {
				service = f.NewService(name, namespace, labels)
				statefulset = f.NewStatefulSet(name, namespace, labels, containers1, service.Name)
			})

			It("Should add annotations", func() {
				if skip {
					Skip("environment var \"DOCKER_CFG\" not found")
				}

				imageAnnotations := map[string]string{
					"docker.com/labels-git-commit": "unkown",
				}
				By("\"git-commit\" annotations should be present")
				f.EventuallyAnnotationsFromStatefulSet(sts.Name, sts.Namespace, "docker.com/").
					Should(Equal(imageAnnotations))

				By("\"hello\" annotation shouldn't be present")
				annotations := f.AnnotaionsFromStatefulSet(sts.Name, sts.Namespace, "docker.com/")
				value, err := meta.GetStringValue(annotations, "docker.com/hi-hello")
				Expect(err).To(HaveOccurred())
				Expect(value).To(Equal(""))
			})
		})

		Context("When docker image doesn't contain any labels", func() {
			BeforeEach(func() {
				service = f.NewService(name, namespace, labels)
				statefulset = f.NewStatefulSet(name, namespace, labels, containers2, service.Name)
			})

			It("Shouldn't add annotations", func() {
				if skip {
					Skip("environment var \"DOCKER_CFG\" not found")
				}

				imageAnnotations := map[string]string{}
				By("\"git-commit\" annotations should be present")
				f.EventuallyAnnotationsFromStatefulSet(sts.Name, sts.Namespace, "docker.com/").
					Should(Equal(imageAnnotations))
			})
		})

		Context("When docker images aren't found with the provided secrets", func() {
			BeforeEach(func() {
				service = f.NewService(name, namespace, labels)
				statefulset = f.NewStatefulSet(name, namespace, labels, containers3, service.Name)
			})

			It("Shouldn't found images", func() {
				if skip {
					Skip("environment var \"DOCKER_CFG\" not found")
				}

				imageAnnotations := map[string]string{}
				By("\"git-commit\" annotations should be present")
				f.EventuallyAnnotationsFromStatefulSet(sts.Name, sts.Namespace, "docker.com/").
					Should(Equal(imageAnnotations))
			})
		})
	})
})
