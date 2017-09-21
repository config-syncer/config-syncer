package v1alpha1

import (
	"encoding/json"
	"fmt"

	"github.com/appscode/kutil"
	aci "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	tcs "github.com/appscode/searchlight/client/typed/monitoring/v1alpha1"
	"github.com/golang/glog"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/jsonmergepatch"
	"k8s.io/apimachinery/pkg/util/wait"
)

func EnsureNodeAlert(c tcs.MonitoringV1alpha1Interface, meta metav1.ObjectMeta, transform func(alert *aci.NodeAlert) *aci.NodeAlert) (*aci.NodeAlert, error) {
	return CreateOrPatchNodeAlert(c, meta, transform)
}

func CreateOrPatchNodeAlert(c tcs.MonitoringV1alpha1Interface, meta metav1.ObjectMeta, transform func(alert *aci.NodeAlert) *aci.NodeAlert) (*aci.NodeAlert, error) {
	cur, err := c.NodeAlerts(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		return c.NodeAlerts(meta.Namespace).Create(transform(&aci.NodeAlert{ObjectMeta: meta}))
	} else if err != nil {
		return nil, err
	}
	return PatchNodeAlert(c, cur, transform)
}

func PatchNodeAlert(c tcs.MonitoringV1alpha1Interface, cur *aci.NodeAlert, transform func(*aci.NodeAlert) *aci.NodeAlert) (*aci.NodeAlert, error) {
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
	glog.V(5).Infof("Patching NodeAlert %s@%s with %s.", cur.Name, cur.Namespace, string(patch))
	result, err := c.NodeAlerts(cur.Namespace).Patch(cur.Name, types.MergePatchType, patch)
	return result, err
}

func TryPatchNodeAlert(c tcs.MonitoringV1alpha1Interface, meta metav1.ObjectMeta, transform func(*aci.NodeAlert) *aci.NodeAlert) (result *aci.NodeAlert, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.NodeAlerts(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = PatchNodeAlert(c, cur, transform)
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to patch NodeAlert %s@%s due to %v.", attempt, cur.Name, cur.Namespace, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to patch NodeAlert %s@%s after %d attempts due to %v", meta.Name, meta.Namespace, attempt, err)
	}
	return
}

func TryUpdateNodeAlert(c tcs.MonitoringV1alpha1Interface, meta metav1.ObjectMeta, transform func(*aci.NodeAlert) *aci.NodeAlert) (result *aci.NodeAlert, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.NodeAlerts(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.NodeAlerts(cur.Namespace).Update(transform(cur))
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to update NodeAlert %s@%s due to %v.", attempt, cur.Name, cur.Namespace, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to update NodeAlert %s@%s after %d attempts due to %v", meta.Name, meta.Namespace, attempt, err)
	}
	return
}
