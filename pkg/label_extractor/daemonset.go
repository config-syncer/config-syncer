package label_extractor

import (
	"encoding/json"

	"github.com/appscode/go/log"
	"github.com/appscode/kutil"
	core_util "github.com/appscode/kutil/core/v1"
	"github.com/golang/glog"
	ext_v1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kubernetes/pkg/util/parsers"
)

func (l *ExtractDockerLabel) ExtractFromDaemonSetHandler() cache.ResourceEventHandler {
	return &ExtractFromDaemonSet{l}
}

type ExtractFromDaemonSet struct {
	*ExtractDockerLabel
}

var _ cache.ResourceEventHandler = &ExtractFromDaemonSet{}

func (ds *ExtractFromDaemonSet) OnAdd(obj interface{}) {
	ds.lock.RLock()

	if !ds.enable {
		return
	}
	ds.lock.RUnlock()

	if res, ok := obj.(*ext_v1.DaemonSet); ok {
		if err := ds.AnnotateDaemonSet(res); err != nil {
			log.Errorln(err)
		}
	}
}

func (ds *ExtractFromDaemonSet) OnUpdate(oldObj, newObj interface{}) {
	ds.lock.RLock()
	if !ds.enable {
		return
	}
	ds.lock.RUnlock()

	oldRes, ok := oldObj.(*ext_v1.DaemonSet)
	if !ok {
		return
	}
	newRes, ok := newObj.(*ext_v1.DaemonSet)
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
		if err := ds.AnnotateDaemonSet(newRes); err != nil {
			log.Errorln(err)
		}
	}
}

func (ds *ExtractFromDaemonSet) OnDelete(obj interface{}) {}

// This method takes a daemonset <ds> and checks if there exists any labels in container images
// at PodTemplateSpec. If exists then add them to annotation of the <ds>. It uses the secrets
// provided at 'imagePullSecrets' for getting labels from images
func (l *ExtractDockerLabel) AnnotateDaemonSet(ds *ext_v1.DaemonSet) error {
	log.Infof("Annotating DaemonSet %s...............\n", ds.Name)
	secretNames := getAllSecrets(ds.Spec.Template.Spec.ImagePullSecrets)

	annotations := make(map[string]string)
	for _, cont := range ds.Spec.Template.Spec.Containers {
		image := cont.Image
		repo, tag, _, err := parsers.ParseImageName(image)
		if err != nil {
			return err
		}
		repoName := repo[10:]

		labels, err := l.GetLabels(ds.ObjectMeta.GetNamespace(), repoName, tag, secretNames)
		if err != nil {
			return err
		}

		prefix := "docker.com/" + cont.Name + "-"
		addPrefixToLabels(labels, prefix)
		core_util.UpsertMap(annotations, labels)
	}

	_, status, err := PatchDS(l.kubeClient, ds, func(daemonset *ext_v1.DaemonSet) *ext_v1.DaemonSet {
		removeOldAnnotations(daemonset.ObjectMeta.Annotations, "docker.com/")
		daemonset.ObjectMeta.SetAnnotations(annotations)

		return daemonset
	})

	log.Infoln("status =", status)
	if err != nil {
		return err
	}

	log.Infof("Annotating DaemonSet %s completed............\n", ds.Name)

	return nil
}

func PatchDS(
	c kubernetes.Interface,
	cur *ext_v1.DaemonSet,
	transform func(*ext_v1.DaemonSet) *ext_v1.DaemonSet) (*ext_v1.DaemonSet, kutil.VerbType, error) {

	curJson, err := json.Marshal(cur)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	modJson, err := json.Marshal(transform(cur.DeepCopy()))
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	patch, err := strategicpatch.CreateTwoWayMergePatch(curJson, modJson, ext_v1.DaemonSet{})
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	if len(patch) == 0 || string(patch) == "{}" {
		return cur, kutil.VerbUnchanged, nil
	}
	glog.V(3).Infof("Patching DaemonSet %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	out, err := c.ExtensionsV1beta1().DaemonSets(cur.Namespace).Patch(cur.Name, types.StrategicMergePatchType, patch)
	return out, kutil.VerbPatched, err
}
