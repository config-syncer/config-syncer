package v1alpha1

import (
	"encoding/json"
	"fmt"

	"github.com/appscode/kutil"
	aci "github.com/appscode/stash/apis/stash/v1alpha1"
	tcs "github.com/appscode/stash/client/typed/stash/v1alpha1"
	"github.com/golang/glog"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/jsonmergepatch"
	"k8s.io/apimachinery/pkg/util/wait"
)

func EnsureRestic(c tcs.StashV1alpha1Interface, meta metav1.ObjectMeta, transform func(alert *aci.Restic) *aci.Restic) (*aci.Restic, error) {
	return CreateOrPatchRestic(c, meta, transform)
}

func CreateOrPatchRestic(c tcs.StashV1alpha1Interface, meta metav1.ObjectMeta, transform func(alert *aci.Restic) *aci.Restic) (*aci.Restic, error) {
	cur, err := c.Restics(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		return c.Restics(meta.Namespace).Create(transform(&aci.Restic{ObjectMeta: meta}))
	} else if err != nil {
		return nil, err
	}
	return PatchRestic(c, cur, transform)
}

func PatchRestic(c tcs.StashV1alpha1Interface, cur *aci.Restic, transform func(*aci.Restic) *aci.Restic) (*aci.Restic, error) {
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
	glog.V(5).Infof("Patching Restic %s@%s with %s.", cur.Name, cur.Namespace, string(patch))
	result, err := c.Restics(cur.Namespace).Patch(cur.Name, types.MergePatchType, patch)
	return result, err
}

func TryPatchRestic(c tcs.StashV1alpha1Interface, meta metav1.ObjectMeta, transform func(*aci.Restic) *aci.Restic) (result *aci.Restic, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.Restics(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = PatchRestic(c, cur, transform)
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to patch Restic %s@%s due to %v.", attempt, cur.Name, cur.Namespace, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to patch Restic %s@%s after %d attempts due to %v", meta.Name, meta.Namespace, attempt, err)
	}
	return
}

func TryUpdateRestic(c tcs.StashV1alpha1Interface, meta metav1.ObjectMeta, transform func(*aci.Restic) *aci.Restic) (result *aci.Restic, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.Restics(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.Restics(cur.Namespace).Update(transform(cur))
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to update Restic %s@%s due to %v.", attempt, cur.Name, cur.Namespace, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to update Restic %s@%s after %d attempts due to %v", meta.Name, meta.Namespace, attempt, err)
	}
	return
}
