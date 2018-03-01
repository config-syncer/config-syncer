package label_extractor

import (
	"github.com/appscode/go/log"
	core_util "github.com/appscode/kutil/core/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kubernetes/pkg/util/parsers"
)

func (l *ExtractDockerLabel) ExtractFromReplicationControllerHandler() cache.ResourceEventHandler {
	return &ExtractFromReplicationController{l}
}

type ExtractFromReplicationController struct {
	*ExtractDockerLabel
}

var _ cache.ResourceEventHandler = &ExtractFromReplicationController{}

func (rc *ExtractFromReplicationController) OnAdd(obj interface{}) {
	rc.lock.RLock()

	if !rc.enable {
		return
	}
	rc.lock.RUnlock()

	if res, ok := obj.(*core.ReplicationController); ok {
		if err := rc.AnnotateReplicationController(res); err != nil {
			log.Errorln(err)
		}
	}
}

func (rc *ExtractFromReplicationController) OnUpdate(oldObj, newObj interface{}) {
	rc.lock.RLock()
	if !rc.enable {
		return
	}
	rc.lock.RUnlock()

	oldRes, ok := oldObj.(*core.ReplicationController)
	if !ok {
		return
	}
	newRes, ok := newObj.(*core.ReplicationController)
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
		if err := rc.AnnotateReplicationController(newRes); err != nil {
			log.Errorln(err)
		}
	}
}

func (rc *ExtractFromReplicationController) OnDelete(obj interface{}) {}

// This method takes a replicationcontroller <rc> and checks if there exists any labels in container
// images at PodTemplateSpec. If exists then add them to annotation of the <rc>. It uses the secrets
// provided at 'imagePullSecrets' for getting labels from images
func (l *ExtractDockerLabel) AnnotateReplicationController(rc *core.ReplicationController) error {
	log.Infof("Annotating ReplicationController %s...............\n", rc.Name)
	secretNames := getAllSecrets(rc.Spec.Template.Spec.ImagePullSecrets)

	annotations := make(map[string]string)
	for _, cont := range rc.Spec.Template.Spec.Containers {
		image := cont.Image
		repo, tag, _, err := parsers.ParseImageName(image)
		if err != nil {
			return err
		}
		repoName := repo[10:]

		labels, err := l.GetLabels(rc.ObjectMeta.GetNamespace(), repoName, tag, secretNames)
		if err != nil {
			return err
		}

		prefix := "docker.com/" + cont.Name + "-"
		addPrefixToLabels(labels, prefix)
		core_util.UpsertMap(annotations, labels)
	}

	_, status, err := core_util.PatchRC(l.kubeClient, rc, func(replicationController *core.ReplicationController) *core.ReplicationController {
		removeOldAnnotations(replicationController.ObjectMeta.Annotations, "docker.com/")
		replicationController.ObjectMeta.SetAnnotations(annotations)

		return replicationController
	})

	log.Infoln("status =", status)
	if err != nil {
		return err
	}

	log.Infof("Annotating Replication Controller %s completed............\n", rc.Name)

	return nil
}
