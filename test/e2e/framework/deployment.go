package framework

import (
	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/go/types"
	apps "k8s.io/api/apps/v1beta1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	TestSourceDataVolumeName = "source-data"
	TestSourceDataMountPath  = "/source/data"
)

func (fi *Invocation) Deployment() *apps.Deployment {
	return &apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rand.WithUniqSuffix("busybox"),
			Namespace: fi.namespace,
			Labels: map[string]string{
				"app": fi.app,
			},
		},
		Spec: apps.DeploymentSpec{
			Replicas: types.Int32P(1),
			Template: fi.PodTemplate(),
		},
	}
}

func (f *Framework) CreateDeployment(obj apps.Deployment) (*apps.Deployment, error) {
	return f.KubeClient.AppsV1beta1().Deployments(obj.Namespace).Create(&obj)
}

func (f *Framework) DeleteDeployment(meta metav1.ObjectMeta) error {
	return f.KubeClient.AppsV1beta1().Deployments(meta.Namespace).Delete(meta.Name, deleteInBackground())
}

func (fi *Invocation) PodTemplate() core.PodTemplateSpec {
	return core.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app": fi.app,
			},
		},
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:            "busybox",
					Image:           "busybox",
					ImagePullPolicy: core.PullIfNotPresent,
					Command: []string{
						"sleep",
						"3600",
					},
					VolumeMounts: []core.VolumeMount{
						{
							Name:      TestSourceDataVolumeName,
							MountPath: TestSourceDataMountPath,
						},
					},
				},
			},
			Volumes: []core.Volume{
				{
					Name: TestSourceDataVolumeName,
					VolumeSource: core.VolumeSource{
						GitRepo: &core.GitRepoVolumeSource{
							Repository: "https://github.com/appscode/stash-data.git",
						},
					},
				},
			},
		},
	}
}
