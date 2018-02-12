package framework

import (
	"strings"

	. "github.com/onsi/gomega"
	apps_v1 "k8s.io/api/apps/v1beta1"
	batch_v1 "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	ext_v1 "k8s.io/api/extensions/v1beta1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func int32Ptr(i int32) *int32 { return &i }

func newObejectMeta(name, namespace string, labels map[string]string) meta_v1.ObjectMeta {
	return meta_v1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Annotations: map[string]string{
			"docker.com/hi-hello": "hello",
		},
		Labels: labels,
	}
}

func (f *Invocation) NewDeployment(
	name, namespace string,
	labels map[string]string,
	containers []core.Container) *apps_v1.Deployment {
	return &apps_v1.Deployment{
		ObjectMeta: newObejectMeta(name, namespace, labels),
		Spec: apps_v1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &meta_v1.LabelSelector{
				MatchLabels: labels,
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: meta_v1.ObjectMeta{
					Labels: labels,
				},
				Spec: core.PodSpec{
					Containers: containers,
					ImagePullSecrets: []core.LocalObjectReference{
						{
							Name: name,
						},
					},
				},
			},
		},
	}
}

func (f *Invocation) NewReplicationController(
	name, namespace string,
	labels map[string]string,
	containers []core.Container) *core.ReplicationController {
	return &core.ReplicationController{
		ObjectMeta: newObejectMeta(name, namespace, labels),
		Spec: core.ReplicationControllerSpec{
			Replicas: int32Ptr(1),
			Selector: labels,
			Template: &core.PodTemplateSpec{
				ObjectMeta: meta_v1.ObjectMeta{
					Labels: labels,
				},
				Spec: core.PodSpec{
					Containers: containers,
					ImagePullSecrets: []core.LocalObjectReference{
						{
							Name: name,
						},
					},
				},
			},
		},
	}
}

func (f *Invocation) NewReplicaSet(
	name, namespace string,
	labels map[string]string,
	containers []core.Container) *ext_v1.ReplicaSet {
	return &ext_v1.ReplicaSet{
		ObjectMeta: newObejectMeta(name, namespace, labels),
		Spec: ext_v1.ReplicaSetSpec{
			Replicas: int32Ptr(1),
			Selector: &meta_v1.LabelSelector{
				MatchLabels: labels,
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: meta_v1.ObjectMeta{
					Labels: labels,
				},
				Spec: core.PodSpec{
					Containers: containers,
					ImagePullSecrets: []core.LocalObjectReference{
						{
							Name: name,
						},
					},
				},
			},
		},
	}
}

func (f *Invocation) NewDaemonSet(
	name, namespace string,
	labels map[string]string,
	containers []core.Container) *ext_v1.DaemonSet {
	return &ext_v1.DaemonSet{
		ObjectMeta: newObejectMeta(name, namespace, labels),
		Spec: ext_v1.DaemonSetSpec{
			Selector: &meta_v1.LabelSelector{
				MatchLabels: labels,
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: meta_v1.ObjectMeta{
					Labels: labels,
				},
				Spec: core.PodSpec{
					Containers: containers,
					ImagePullSecrets: []core.LocalObjectReference{
						{
							Name: name,
						},
					},
				},
			},
		},
	}
}

func (f *Invocation) NewJob(
	name, namespace string,
	labels map[string]string,
	containers []core.Container) *batch_v1.Job {
	return &batch_v1.Job{
		ObjectMeta: newObejectMeta(name, namespace, labels),
		Spec: batch_v1.JobSpec{
			Template: core.PodTemplateSpec{
				ObjectMeta: meta_v1.ObjectMeta{
					Labels: labels,
				},
				Spec: core.PodSpec{
					Containers: containers,
					ImagePullSecrets: []core.LocalObjectReference{
						{
							Name: name,
						},
					},
					RestartPolicy: "Never",
				},
			},
		},
	}
}

func (f *Invocation) NewService(
	name, namespace string,
	labels map[string]string) *core.Service {
	return &core.Service{
		ObjectMeta: newObejectMeta(name, namespace, labels),
		Spec: core.ServiceSpec{
			Ports: []core.ServicePort{
				{
					Protocol: core.ProtocolTCP,
					Port:     80,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 80,
					},
				},
			},
			Type:     core.ServiceTypeNodePort,
			Selector: labels,
		},
	}
}

func (f *Invocation) NewStatefulSet(
	name, namespace string,
	labels map[string]string,
	containers []core.Container,
	svcName string) *apps_v1.StatefulSet {
	return &apps_v1.StatefulSet{
		ObjectMeta: newObejectMeta(name, namespace, labels),
		Spec: apps_v1.StatefulSetSpec{
			ServiceName: name,
			Replicas:    int32Ptr(1),
			Selector: &meta_v1.LabelSelector{
				MatchLabels: labels,
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: meta_v1.ObjectMeta{
					Labels: labels,
				},
				Spec: core.PodSpec{
					Containers: containers,
					ImagePullSecrets: []core.LocalObjectReference{
						{
							Name: name,
						},
					},
				},
			},
		},
	}
}

func (f *Invocation) NewSecret(name, namespace, data string, labels map[string]string) *core.Secret {
	return &core.Secret{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		StringData: map[string]string{
			"config.json": data,
		},
	}
}

func (f *Invocation) EventuallyAnnotationsFromDeployment(name, namespace, prefix string) GomegaAsyncAssertion {
	return Eventually(func() map[string]string {
		return f.AnnotaionsFromDeployment(name, namespace, prefix)
	})
}

func (f *Invocation) EventuallyAnnotationsFromReplicationController(name, namespace, prefix string) GomegaAsyncAssertion {
	return Eventually(func() map[string]string {
		return f.AnnotaionsFromReplicationController(name, namespace, prefix)
	})
}

func (f *Invocation) EventuallyAnnotationsFromReplicaSet(name, namespace, prefix string) GomegaAsyncAssertion {
	return Eventually(func() map[string]string {
		return f.AnnotaionsFromReplicaSet(name, namespace, prefix)
	})
}

func (f *Invocation) EventuallyAnnotationsFromDaemonSet(name, namespace, prefix string) GomegaAsyncAssertion {
	return Eventually(func() map[string]string {
		return f.AnnotaionsFromDaemonSet(name, namespace, prefix)
	})
}

func (f *Invocation) EventuallyAnnotationsFromJob(name, namespace, prefix string) GomegaAsyncAssertion {
	return Eventually(func() map[string]string {
		return f.AnnotaionsFromJob(name, namespace, prefix)
	})
}

func (f *Invocation) EventuallyAnnotationsFromStatefulSet(name, namespace, prefix string) GomegaAsyncAssertion {
	return Eventually(func() map[string]string {
		return f.AnnotaionsFromStatefulSet(name, namespace, prefix)
	})
}

func annotaionsWhoseKeyHasPrefix(annotations map[string]string, prefix string) map[string]string {
	res := map[string]string{}
	if annotations == nil {
		return res
	}

	for key, val := range annotations {
		if strings.HasPrefix(key, prefix) {
			res[key] = val
		}
	}

	return res
}

func (f *Invocation) AnnotaionsFromDeployment(deployName, namespace, prefix string) map[string]string {
	deploy, err := f.KubeClient.AppsV1beta1().Deployments(namespace).Get(deployName, meta_v1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())

	return annotaionsWhoseKeyHasPrefix(deploy.ObjectMeta.Annotations, prefix)
}

func (f *Invocation) AnnotaionsFromReplicationController(rcName, namespace, prefix string) map[string]string {
	rc, err := f.KubeClient.CoreV1().ReplicationControllers(namespace).Get(rcName, meta_v1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())

	return annotaionsWhoseKeyHasPrefix(rc.ObjectMeta.Annotations, prefix)
}

func (f *Invocation) AnnotaionsFromReplicaSet(rsName, namespace, prefix string) map[string]string {
	rs, err := f.KubeClient.ExtensionsV1beta1().ReplicaSets(namespace).Get(rsName, meta_v1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())

	return annotaionsWhoseKeyHasPrefix(rs.ObjectMeta.Annotations, prefix)
}

func (f *Invocation) AnnotaionsFromDaemonSet(dsName, namespace, prefix string) map[string]string {
	ds, err := f.KubeClient.ExtensionsV1beta1().DaemonSets(namespace).Get(dsName, meta_v1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())

	return annotaionsWhoseKeyHasPrefix(ds.ObjectMeta.Annotations, prefix)
}

func (f *Invocation) AnnotaionsFromJob(jobName, namespace, prefix string) map[string]string {
	job, err := f.KubeClient.BatchV1().Jobs(namespace).Get(jobName, meta_v1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())

	return annotaionsWhoseKeyHasPrefix(job.ObjectMeta.Annotations, prefix)
}

func (f *Invocation) AnnotaionsFromStatefulSet(stsName, namespace, prefix string) map[string]string {
	sts, err := f.KubeClient.AppsV1beta1().StatefulSets(namespace).Get(stsName, meta_v1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())

	return annotaionsWhoseKeyHasPrefix(sts.ObjectMeta.Annotations, prefix)
}

func (f *Invocation) DeleteAllDeployments() {
	deployments, err := f.KubeClient.AppsV1beta1().Deployments(meta_v1.NamespaceAll).List(meta_v1.ListOptions{
		LabelSelector: labels.Set{
			"app": f.App(),
		}.String(),
	})
	Expect(err).NotTo(HaveOccurred())

	for _, deploy := range deployments.Items {
		err := f.KubeClient.AppsV1beta1().Deployments(deploy.Namespace).Delete(deploy.Name, &meta_v1.DeleteOptions{})
		if kerr.IsNotFound(err) {
			err = nil
		}
		Expect(err).NotTo(HaveOccurred())
	}
}

func (f *Invocation) DeleteAllReplicationControllers() {
	replicationcontrollers, err := f.KubeClient.CoreV1().ReplicationControllers(meta_v1.NamespaceAll).List(meta_v1.ListOptions{
		LabelSelector: labels.Set{
			"app": f.App(),
		}.String(),
	})
	Expect(err).NotTo(HaveOccurred())

	for _, rc := range replicationcontrollers.Items {
		err := f.KubeClient.CoreV1().ReplicationControllers(rc.Namespace).Delete(rc.Name, &meta_v1.DeleteOptions{})
		if kerr.IsNotFound(err) {
			err = nil
		}
		Expect(err).NotTo(HaveOccurred())
	}
}

func (f *Invocation) DeleteAllReplicasets() {
	replicasets, err := f.KubeClient.ExtensionsV1beta1().ReplicaSets(meta_v1.NamespaceAll).List(meta_v1.ListOptions{
		LabelSelector: labels.Set{
			"app": f.App(),
		}.String(),
	})
	Expect(err).NotTo(HaveOccurred())

	for _, rs := range replicasets.Items {
		err := f.KubeClient.ExtensionsV1beta1().ReplicaSets(rs.Namespace).Delete(rs.Name, &meta_v1.DeleteOptions{})
		if kerr.IsNotFound(err) {
			err = nil
		}
		Expect(err).NotTo(HaveOccurred())
	}
}

func (f *Invocation) DeleteAllDaemonSet() {
	daemonsets, err := f.KubeClient.ExtensionsV1beta1().DaemonSets(meta_v1.NamespaceAll).List(meta_v1.ListOptions{
		LabelSelector: labels.Set{
			"app": f.App(),
		}.String(),
	})
	Expect(err).NotTo(HaveOccurred())

	for _, ds := range daemonsets.Items {
		err := f.KubeClient.ExtensionsV1beta1().DaemonSets(ds.Namespace).Delete(ds.Name, &meta_v1.DeleteOptions{})
		if kerr.IsNotFound(err) {
			err = nil
		}
		Expect(err).NotTo(HaveOccurred())
	}
}

func (f *Invocation) DeleteAllJobs() {
	jobs, err := f.KubeClient.BatchV1().Jobs(meta_v1.NamespaceAll).List(meta_v1.ListOptions{
		LabelSelector: labels.Set{
			"app": f.App(),
		}.String(),
	})
	Expect(err).NotTo(HaveOccurred())

	for _, job := range jobs.Items {
		err := f.KubeClient.BatchV1().Jobs(job.Namespace).Delete(job.Name, &meta_v1.DeleteOptions{})
		if kerr.IsNotFound(err) {
			err = nil
		}
		Expect(err).NotTo(HaveOccurred())
	}
}

func (f *Invocation) DeleteAllStatefulSets() {
	statefulsets, err := f.KubeClient.AppsV1beta1().StatefulSets(meta_v1.NamespaceAll).List(meta_v1.ListOptions{
		LabelSelector: labels.Set{
			"app": f.App(),
		}.String(),
	})
	Expect(err).NotTo(HaveOccurred())

	for _, sts := range statefulsets.Items {
		err := f.KubeClient.AppsV1beta1().StatefulSets(sts.Namespace).Delete(sts.Name, &meta_v1.DeleteOptions{})
		if kerr.IsNotFound(err) {
			err = nil
		}
		Expect(err).NotTo(HaveOccurred())
	}
}

func (f *Invocation) DeleteAllServices() {
	services, err := f.KubeClient.CoreV1().Services(meta_v1.NamespaceAll).List(meta_v1.ListOptions{
		LabelSelector: labels.Set{
			"app": f.App(),
		}.String(),
	})
	Expect(err).NotTo(HaveOccurred())

	for _, svc := range services.Items {
		err := f.KubeClient.CoreV1().Services(svc.Namespace).Delete(svc.Name, &meta_v1.DeleteOptions{})
		if kerr.IsNotFound(err) {
			err = nil
		}
		Expect(err).NotTo(HaveOccurred())
	}
}
