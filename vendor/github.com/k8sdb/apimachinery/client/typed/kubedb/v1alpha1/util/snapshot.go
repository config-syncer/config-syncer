package util

import (
	"encoding/json"
	"fmt"

	"github.com/appscode/kutil"
	"github.com/golang/glog"
	aci "github.com/k8sdb/apimachinery/apis/kubedb/v1alpha1"
	tcs "github.com/k8sdb/apimachinery/client/typed/kubedb/v1alpha1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/jsonmergepatch"
	"k8s.io/apimachinery/pkg/util/wait"
)

func EnsureSnapshot(c tcs.KubedbV1alpha1Interface, meta metav1.ObjectMeta, transform func(alert *aci.Snapshot) *aci.Snapshot) (*aci.Snapshot, error) {
	return CreateOrPatchSnapshot(c, meta, transform)
}

func CreateOrPatchSnapshot(c tcs.KubedbV1alpha1Interface, meta metav1.ObjectMeta, transform func(alert *aci.Snapshot) *aci.Snapshot) (*aci.Snapshot, error) {
	cur, err := c.Snapshots(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		glog.V(3).Infof("Creating Snapshot %s/%s.", meta.Namespace, meta.Name)
		return c.Snapshots(meta.Namespace).Create(transform(&aci.Snapshot{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Snapshot",
				APIVersion: aci.SchemeGroupVersion.String(),
			},
			ObjectMeta: meta,
		}))
	} else if err != nil {
		return nil, err
	}
	return PatchSnapshot(c, cur, transform)
}

func PatchSnapshot(c tcs.KubedbV1alpha1Interface, cur *aci.Snapshot, transform func(*aci.Snapshot) *aci.Snapshot) (*aci.Snapshot, error) {
	curJson, err := json.Marshal(cur)
	if err != nil {
		return nil, err
	}

	modJson, err := json.Marshal(transform(cur))
	if err != nil {
		return nil, err
	}

	patch, err := jsonmergepatch.CreateThreeWayJSONMergePatch(curJson, modJson, curJson)
	if err != nil {
		return nil, err
	}
	if len(patch) == 0 || string(patch) == "{}" {
		return cur, nil
	}
	glog.V(3).Infof("Patching Snapshot %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	result, err := c.Snapshots(cur.Namespace).Patch(cur.Name, types.MergePatchType, patch)
	return result, err
}

func TryPatchSnapshot(c tcs.KubedbV1alpha1Interface, meta metav1.ObjectMeta, transform func(*aci.Snapshot) *aci.Snapshot) (result *aci.Snapshot, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.Snapshots(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = PatchSnapshot(c, cur, transform)
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to patch Snapshot %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to patch Snapshot %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}

func TryUpdateSnapshot(c tcs.KubedbV1alpha1Interface, meta metav1.ObjectMeta, transform func(*aci.Snapshot) *aci.Snapshot) (result *aci.Snapshot, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.Snapshots(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.Snapshots(cur.Namespace).Update(transform(cur))
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to update Snapshot %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to update Snapshot %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}
