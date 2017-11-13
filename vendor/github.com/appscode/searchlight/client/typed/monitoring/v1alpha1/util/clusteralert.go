package util

import (
	"encoding/json"
	"fmt"

	"github.com/appscode/kutil"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	cs "github.com/appscode/searchlight/client/typed/monitoring/v1alpha1"
	"github.com/golang/glog"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/jsonmergepatch"
	"k8s.io/apimachinery/pkg/util/wait"
)

func EnsureClusterAlert(c cs.MonitoringV1alpha1Interface, meta metav1.ObjectMeta, transform func(alert *api.ClusterAlert) *api.ClusterAlert) (*api.ClusterAlert, error) {
	return CreateOrPatchClusterAlert(c, meta, transform)
}

func CreateOrPatchClusterAlert(c cs.MonitoringV1alpha1Interface, meta metav1.ObjectMeta, transform func(alert *api.ClusterAlert) *api.ClusterAlert) (*api.ClusterAlert, error) {
	cur, err := c.ClusterAlerts(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		glog.V(3).Infof("Creating ClusterAlert %s/%s.", meta.Namespace, meta.Name)
		return c.ClusterAlerts(meta.Namespace).Create(transform(&api.ClusterAlert{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ClusterAlert",
				APIVersion: api.SchemeGroupVersion.String(),
			},
			ObjectMeta: meta,
		}))
	} else if err != nil {
		return nil, err
	}
	return PatchClusterAlert(c, cur, transform)
}

func PatchClusterAlert(c cs.MonitoringV1alpha1Interface, cur *api.ClusterAlert, transform func(*api.ClusterAlert) *api.ClusterAlert) (*api.ClusterAlert, error) {
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
	glog.V(3).Infof("Patching ClusterAlert %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	result, err := c.ClusterAlerts(cur.Namespace).Patch(cur.Name, types.MergePatchType, patch)
	return result, err
}

func TryPatchClusterAlert(c cs.MonitoringV1alpha1Interface, meta metav1.ObjectMeta, transform func(*api.ClusterAlert) *api.ClusterAlert) (result *api.ClusterAlert, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.ClusterAlerts(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = PatchClusterAlert(c, cur, transform)
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to patch ClusterAlert %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to patch ClusterAlert %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}

func TryUpdateClusterAlert(c cs.MonitoringV1alpha1Interface, meta metav1.ObjectMeta, transform func(*api.ClusterAlert) *api.ClusterAlert) (result *api.ClusterAlert, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.ClusterAlerts(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.ClusterAlerts(cur.Namespace).Update(transform(cur.DeepCopy()))
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to update ClusterAlert %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to update ClusterAlert %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}
