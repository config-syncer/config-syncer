package util

import (
	"encoding/json"
	"fmt"

	"github.com/appscode/kutil"
	"github.com/golang/glog"
	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	cs "github.com/kubedb/apimachinery/client/typed/kubedb/v1alpha1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/jsonmergepatch"
	"k8s.io/apimachinery/pkg/util/wait"
)

func CreateOrPatchDormantDatabase(c cs.KubedbV1alpha1Interface, meta metav1.ObjectMeta, transform func(*api.DormantDatabase) *api.DormantDatabase) (*api.DormantDatabase, kutil.VerbType, error) {
	cur, err := c.DormantDatabases(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		glog.V(3).Infof("Creating DormantDatabase %s/%s.", meta.Namespace, meta.Name)
		out, err := c.DormantDatabases(meta.Namespace).Create(transform(&api.DormantDatabase{
			TypeMeta: metav1.TypeMeta{
				Kind:       "DormantDatabase",
				APIVersion: api.SchemeGroupVersion.String(),
			},
			ObjectMeta: meta,
		}))
		return out, kutil.VerbCreated, err
	} else if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	return PatchDormantDatabase(c, cur, transform)
}

func PatchDormantDatabase(c cs.KubedbV1alpha1Interface, cur *api.DormantDatabase, transform func(*api.DormantDatabase) *api.DormantDatabase) (*api.DormantDatabase, kutil.VerbType, error) {
	curJson, err := json.Marshal(cur)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	modJson, err := json.Marshal(transform(cur.DeepCopy()))
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	patch, err := jsonmergepatch.CreateThreeWayJSONMergePatch(curJson, modJson, curJson)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	if len(patch) == 0 || string(patch) == "{}" {
		return cur, kutil.VerbUnchanged, nil
	}
	glog.V(3).Infof("Patching DormantDatabase %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	out, err := c.DormantDatabases(cur.Namespace).Patch(cur.Name, types.MergePatchType, patch)
	return out, kutil.VerbPatched, err
}

func TryUpdateDormantDatabase(c cs.KubedbV1alpha1Interface, meta metav1.ObjectMeta, transform func(*api.DormantDatabase) *api.DormantDatabase) (result *api.DormantDatabase, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.DormantDatabases(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.DormantDatabases(cur.Namespace).Update(transform(cur.DeepCopy()))
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to update DormantDatabase %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to update DormantDatabase %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}

func DeleteDormantDatabase(c cs.KubedbV1alpha1Interface, meta metav1.ObjectMeta) (err error) {
	err = c.DormantDatabases(meta.Namespace).Delete(meta.Name, nil)
	if err != nil {
		return
	}
	return wait.PollImmediate(kutil.RetryInterval, kutil.ReadinessTimeout, func() (bool, error) {
		_, err := c.DormantDatabases(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if err != nil {
			if kerr.IsNotFound(err) {
				return true, nil
			}
		}
		return false, nil
	})
}
