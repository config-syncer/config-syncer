package label_extractor

import (
	"encoding/json"

	"github.com/appscode/go/log"
	"github.com/appscode/kutil"
	core_util "github.com/appscode/kutil/core/v1"
	"github.com/golang/glog"
	batch "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kubernetes/pkg/util/parsers"
)

func (l *ExtractDockerLabel) ExtractFromJobHandler() cache.ResourceEventHandler {
	return &ExtractFromJob{l}
}

type ExtractFromJob struct {
	*ExtractDockerLabel
}

var _ cache.ResourceEventHandler = &ExtractFromJob{}

func (j *ExtractFromJob) OnAdd(obj interface{}) {
	j.lock.RLock()

	if !j.enable {
		return
	}
	j.lock.RUnlock()

	if res, ok := obj.(*batch.Job); ok {
		if err := j.AnnotateJob(res); err != nil {
			log.Errorln(err)
		}
	}
}

func (j *ExtractFromJob) OnUpdate(oldObj, newObj interface{}) {
	j.lock.RLock()
	if !j.enable {
		return
	}
	j.lock.RUnlock()

	oldRes, ok := oldObj.(*batch.Job)
	if !ok {
		return
	}
	newRes, ok := newObj.(*batch.Job)
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
		if err := j.AnnotateJob(newRes); err != nil {
			log.Errorln(err)
		}
	}
}

func (j *ExtractFromJob) OnDelete(obj interface{}) {}

// This method takes a job <job> and checks if there exists any labels in container images
// at PodTemplateSpec. If exists then add them to annotation of the <job>. It uses the secrets
// provided at 'imagePullSecrets' for getting labels from images
func (l *ExtractDockerLabel) AnnotateJob(job *batch.Job) error {
	log.Infof("Annotating Job %s...............\n", job.Name)
	secretNames := getAllSecrets(job.Spec.Template.Spec.ImagePullSecrets)

	annotations := make(map[string]string)
	for _, cont := range job.Spec.Template.Spec.Containers {
		image := cont.Image
		repo, tag, _, err := parsers.ParseImageName(image)
		if err != nil {
			return err
		}
		repoName := repo[10:]

		labels, err := l.GetLabels(job.ObjectMeta.GetNamespace(), repoName, tag, secretNames)
		if err != nil {
			return err
		}

		prefix := "docker.com/" + cont.Name + "-"
		addPrefixToLabels(labels, prefix)
		core_util.UpsertMap(annotations, labels)
	}

	_, status, err := PatchJob(l.kubeClient, job, func(curJob *batch.Job) *batch.Job {
		removeOldAnnotations(curJob.ObjectMeta.Annotations, "docker.com/")
		curJob.ObjectMeta.SetAnnotations(annotations)

		return curJob
	})

	log.Infoln("status =", status)
	if err != nil {
		return err
	}

	log.Infof("Annotating Job %s completed............\n", job.Name)

	return nil
}

func PatchJob(
	c kubernetes.Interface,
	cur *batch.Job,
	transform func(*batch.Job) *batch.Job) (*batch.Job, kutil.VerbType, error) {

	curJson, err := json.Marshal(cur)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	modJson, err := json.Marshal(transform(cur.DeepCopy()))
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	patch, err := strategicpatch.CreateTwoWayMergePatch(curJson, modJson, batch.Job{})
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	if len(patch) == 0 || string(patch) == "{}" {
		return cur, kutil.VerbUnchanged, nil
	}
	glog.V(3).Infof("Patching Job %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	out, err := c.BatchV1().Jobs(cur.Namespace).Patch(cur.Name, types.StrategicMergePatchType, patch)
	return out, kutil.VerbPatched, err
}
