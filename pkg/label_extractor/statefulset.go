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

func (l *ExtractDockerLabel) ExtractFromStatefulSetHandler() cache.ResourceEventHandler {
	return &ExtractFromStatefulSet{l}
}

type ExtractFromStatefulSet struct {
	*ExtractDockerLabel
}

var _ cache.ResourceEventHandler = &ExtractFromStatefulSet{}

func (sts *ExtractFromStatefulSet) OnAdd(obj interface{}) {
	sts.lock.RLock()

	if !sts.enable {
		return
	}
	sts.lock.RUnlock()

	if res, ok := obj.(*v1beta1.StatefulSet); ok {
		if err := sts.AnnotateStatefulSet(res); err != nil {
			log.Errorln(err)
		}
	}
}

func (sts *ExtractFromStatefulSet) OnUpdate(oldObj, newObj interface{}) {
	sts.lock.RLock()
	if !sts.enable {
		return
	}
	sts.lock.RUnlock()

	oldRes, ok := oldObj.(*v1beta1.StatefulSet)
	if !ok {
		return
	}
	newRes, ok := newObj.(*v1beta1.StatefulSet)
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
		if err := sts.AnnotateStatefulSet(newRes); err != nil {
			log.Errorln(err)
		}
	}
}

func (sts *ExtractFromStatefulSet) OnDelete(obj interface{}) {}

// This method takes a statefulset <sts> and checks if there exists any labels in container images
// at PodTemplateSpec. If exists then add them to annotation of the <sts>. It uses the secrets
// provided at 'imagePullSecrets' for getting labels from images
func (l *ExtractDockerLabel) AnnotateStatefulSet(sts *v1beta1.StatefulSet) error {
	log.Infof("Annotating StatefulSet %s...............\n", sts.Name)
	secretNames := getAllSecrets(sts.Spec.Template.Spec.ImagePullSecrets)

	annotations := make(map[string]string)
	for _, cont := range sts.Spec.Template.Spec.Containers {
		image := cont.Image
		repo, tag, _, err := parsers.ParseImageName(image)
		if err != nil {
			return err
		}
		repoName := repo[10:]

		labels, err := l.GetLabels(sts.ObjectMeta.GetNamespace(), repoName, tag, secretNames)
		if err != nil {
			return err
		}

		prefix := "docker.com/" + cont.Name + "-"
		addPrefixToLabels(labels, prefix)
		core_util.UpsertMap(annotations, labels)
	}

	_, status, err := apps_util.PatchStatefulSet(l.kubeClient, sts, func(statefulSet *v1beta1.StatefulSet) *v1beta1.StatefulSet {
		removeOldAnnotations(statefulSet.ObjectMeta.Annotations, "docker.com/")
		statefulSet.ObjectMeta.SetAnnotations(annotations)

		return statefulSet
	})

	log.Infoln("status =", status)
	if err != nil {
		return err
	}

	log.Infof("Annotating StatefulSet %s completed............\n", sts.Name)

	return nil
}
