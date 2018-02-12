package label_extractor

import (
	"github.com/appscode/go/log"
	apps_util "github.com/appscode/kutil/apps/v1beta1"
	core_util "github.com/appscode/kutil/core/v1"
	"k8s.io/api/apps/v1beta1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kubernetes/pkg/util/parsers"
)

func (l *ExtractDockerLabel) ExtractFromDeploymentHandler() cache.ResourceEventHandler {
	return &ExtractFromDeployment{l}
}

type ExtractFromDeployment struct {
	*ExtractDockerLabel
}

var _ cache.ResourceEventHandler = &ExtractFromDeployment{}

func (d *ExtractFromDeployment) OnAdd(obj interface{}) {
	d.lock.RLock()

	if !d.enable {
		return
	}
	d.lock.RUnlock()

	if res, ok := obj.(*v1beta1.Deployment); ok {
		if err := d.AnnotateDeployment(res); err != nil {
			log.Errorln(err)
		}
	}
}

func (d *ExtractFromDeployment) OnUpdate(oldObj, newObj interface{}) {
	d.lock.RLock()
	if !d.enable {
		return
	}
	d.lock.RUnlock()

	oldRes, ok := oldObj.(*v1beta1.Deployment)
	if !ok {
		return
	}
	newRes, ok := newObj.(*v1beta1.Deployment)
	if !ok {
		return
	}

	// if container images don't match then annotate
	oldContainers := sets.String{}
	for _, cont := range oldRes.Spec.Template.Spec.Containers {
		oldContainers.Insert(cont.Image)
	}
	newContainers := sets.String{}
	for _, cont := range newRes.Spec.Template.Spec.Containers {
		newContainers.Insert(cont.Image)
	}
	if oldContainers.Equal(newContainers) == false {
		if err := d.AnnotateDeployment(newRes); err != nil {
			log.Errorln(err)
		}
	}
}

func (d *ExtractFromDeployment) OnDelete(obj interface{}) {}

// This method takes a deployment <deploy> and checks if there exists any labels in container images
// at PodTemplateSpec. If exists then add them to annotation of the <deploy>. It uses the secrets
// provided at 'imagePullSecrets' for getting labels from images
func (l *ExtractDockerLabel) AnnotateDeployment(deploy *v1beta1.Deployment) error {
	log.Infof("Annotating Deployment %s...............\n", deploy.Name)
	secretNames := getAllSecrets(deploy.Spec.Template.Spec.ImagePullSecrets)

	annotations := make(map[string]string)
	for _, cont := range deploy.Spec.Template.Spec.Containers {
		image := cont.Image
		repo, tag, _, err := parsers.ParseImageName(image)
		if err != nil {
			return err
		}
		repoName := repo[10:]

		labels, err := l.GetLabels(deploy.ObjectMeta.GetNamespace(), repoName, tag, secretNames)
		if err != nil {
			return err
		}

		prefix := "docker.com/" + cont.Name + "-"
		addPrefixToLabels(labels, prefix)
		core_util.UpsertMap(annotations, labels)
	}

	_, status, err := apps_util.PatchDeployment(l.kubeClient, deploy, func(deployment *v1beta1.Deployment) *v1beta1.Deployment {
		removeOldAnnotations(deployment.ObjectMeta.Annotations, "docker.com/")
		deployment.ObjectMeta.SetAnnotations(annotations)

		return deployment
	})

	log.Infoln("status =", status)
	if err != nil {
		return err
	}

	log.Infof("Annotating Deployment %s completed............\n", deploy.Name)

	return nil
}
