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

func EnsureDormantDatabase(c tcs.KubedbV1alpha1Interface, meta metav1.ObjectMeta, transform func(alert *aci.DormantDatabase) *aci.DormantDatabase) (*aci.DormantDatabase, error) {
	return CreateOrPatchDormantDatabase(c, meta, transform)
}

func CreateOrPatchDormantDatabase(c tcs.KubedbV1alpha1Interface, meta metav1.ObjectMeta, transform func(alert *aci.DormantDatabase) *aci.DormantDatabase) (*aci.DormantDatabase, error) {
	cur, err := c.DormantDatabases(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		glog.V(3).Infof("Creating DormantDatabase %s/%s.", meta.Namespace, meta.Name)
		return c.DormantDatabases(meta.Namespace).Create(transform(&aci.DormantDatabase{
			TypeMeta: metav1.TypeMeta{
				Kind:       "DormantDatabase",
				APIVersion: aci.SchemeGroupVersion.String(),
			},
			ObjectMeta: meta,
		}))
	} else if err != nil {
		return nil, err
	}
	return PatchDormantDatabase(c, cur, transform)
}

func PatchDormantDatabase(c tcs.KubedbV1alpha1Interface, cur *aci.DormantDatabase, transform func(*aci.DormantDatabase) *aci.DormantDatabase) (*aci.DormantDatabase, error) {
	curJson, err := json.Marshal(cur)
	if err != nil {
		return nil, err
	}

	modJson, err := json.Marshal(transform(cur.DeepCopy()))
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
	glog.V(3).Infof("Patching DormantDatabase %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	result, err := c.DormantDatabases(cur.Namespace).Patch(cur.Name, types.MergePatchType, patch)
	return result, err
}

func TryPatchDormantDatabase(c tcs.KubedbV1alpha1Interface, meta metav1.ObjectMeta, transform func(*aci.DormantDatabase) *aci.DormantDatabase) (result *aci.DormantDatabase, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.DormantDatabases(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = PatchDormantDatabase(c, cur, transform)
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to patch DormantDatabase %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to patch DormantDatabase %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}

func TryUpdateDormantDatabase(c tcs.KubedbV1alpha1Interface, meta metav1.ObjectMeta, transform func(*aci.DormantDatabase) *aci.DormantDatabase) (result *aci.DormantDatabase, err error) {
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
