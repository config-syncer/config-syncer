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

func (l *ExtractDockerLabel) ExtractFromReplicaSetHandler() cache.ResourceEventHandler {
	return &ExtractFromReplicaSet{l}
}

type ExtractFromReplicaSet struct {
	*ExtractDockerLabel
}

var _ cache.ResourceEventHandler = &ExtractFromReplicaSet{}

func (rs *ExtractFromReplicaSet) OnAdd(obj interface{}) {
	rs.lock.RLock()

	if !rs.enable {
		return
	}
	rs.lock.RUnlock()

	if res, ok := obj.(*ext_v1.ReplicaSet); ok {
		if err := rs.AnnotateReplicaSet(res); err != nil {
			log.Errorln(err)
		}
	}
}

func (rs *ExtractFromReplicaSet) OnUpdate(oldObj, newObj interface{}) {
	rs.lock.RLock()
	if !rs.enable {
		return
	}
	rs.lock.RUnlock()

	oldRes, ok := oldObj.(*ext_v1.ReplicaSet)
	if !ok {
		return
	}
	newRes, ok := newObj.(*ext_v1.ReplicaSet)
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
		if err := rs.AnnotateReplicaSet(newRes); err != nil {
			log.Errorln(err)
		}
	}
}

func (rs *ExtractFromReplicaSet) OnDelete(obj interface{}) {}

// This method takes a replicaset <rs> and checks if there exists any labels in container images
// at PodTemplateSpec. If exists then add them to annotation of the <rs>. It uses the secrets
// provided at 'imagePullSecrets' for getting labels from images
func (l *ExtractDockerLabel) AnnotateReplicaSet(rs *ext_v1.ReplicaSet) error {
	log.Infof("Annotating ReplicaSet %s...............\n", rs.Name)
	secretNames := getAllSecrets(rs.Spec.Template.Spec.ImagePullSecrets)

	annotations := make(map[string]string)
	for _, cont := range rs.Spec.Template.Spec.Containers {
		image := cont.Image
		repo, tag, _, err := parsers.ParseImageName(image)
		if err != nil {
			return err
		}
		repoName := repo[10:]

		labels, err := l.GetLabels(rs.ObjectMeta.GetNamespace(), repoName, tag, secretNames)
		if err != nil {
			return err
		}

		prefix := "docker.com/" + cont.Name + "-"
		addPrefixToLabels(labels, prefix)
		core_util.UpsertMap(annotations, labels)
	}

	_, status, err := PatchRS(l.kubeClient, rs, func(replicaset *ext_v1.ReplicaSet) *ext_v1.ReplicaSet {
		removeOldAnnotations(replicaset.ObjectMeta.Annotations, "docker.com/")
		replicaset.ObjectMeta.SetAnnotations(annotations)

		return replicaset
	})

	log.Infoln("status =", status)
	if err != nil {
		return err
	}

	log.Infof("Annotating ReplicaSet %s completed............\n", rs.Name)

	return nil
}

func PatchRS(
	c kubernetes.Interface,
	cur *ext_v1.ReplicaSet,
	transform func(*ext_v1.ReplicaSet) *ext_v1.ReplicaSet) (*ext_v1.ReplicaSet, kutil.VerbType, error) {

	curJson, err := json.Marshal(cur)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	modJson, err := json.Marshal(transform(cur.DeepCopy()))
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	patch, err := strategicpatch.CreateTwoWayMergePatch(curJson, modJson, ext_v1.ReplicaSet{})
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	if len(patch) == 0 || string(patch) == "{}" {
		return cur, kutil.VerbUnchanged, nil
	}
	glog.V(3).Infof("Patching ReplicaSet %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	out, err := c.ExtensionsV1beta1().ReplicaSets(cur.Namespace).Patch(cur.Name, types.StrategicMergePatchType, patch)
	return out, kutil.VerbPatched, err
}
