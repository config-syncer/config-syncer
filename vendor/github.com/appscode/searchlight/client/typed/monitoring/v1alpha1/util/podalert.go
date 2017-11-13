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

func EnsurePodAlert(c cs.MonitoringV1alpha1Interface, meta metav1.ObjectMeta, transform func(alert *api.PodAlert) *api.PodAlert) (*api.PodAlert, error) {
	return CreateOrPatchPodAlert(c, meta, transform)
}

func CreateOrPatchPodAlert(c cs.MonitoringV1alpha1Interface, meta metav1.ObjectMeta, transform func(alert *api.PodAlert) *api.PodAlert) (*api.PodAlert, error) {
	cur, err := c.PodAlerts(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		glog.V(3).Infof("Creating PodAlert %s/%s.", meta.Namespace, meta.Name)
		return c.PodAlerts(meta.Namespace).Create(transform(&api.PodAlert{
			TypeMeta: metav1.TypeMeta{
				Kind:       "PodAlert",
				APIVersion: api.SchemeGroupVersion.String(),
			},
			ObjectMeta: meta,
		}))
	} else if err != nil {
		return nil, err
	}
	return PatchPodAlert(c, cur, transform)
}

func PatchPodAlert(c cs.MonitoringV1alpha1Interface, cur *api.PodAlert, transform func(*api.PodAlert) *api.PodAlert) (*api.PodAlert, error) {
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
	glog.V(3).Infof("Patching PodAlert %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	result, err := c.PodAlerts(cur.Namespace).Patch(cur.Name, types.MergePatchType, patch)
	return result, err
}

func TryPatchPodAlert(c cs.MonitoringV1alpha1Interface, meta metav1.ObjectMeta, transform func(*api.PodAlert) *api.PodAlert) (result *api.PodAlert, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.PodAlerts(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = PatchPodAlert(c, cur, transform)
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to patch PodAlert %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to patch PodAlert %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}

func TryUpdatePodAlert(c cs.MonitoringV1alpha1Interface, meta metav1.ObjectMeta, transform func(*api.PodAlert) *api.PodAlert) (result *api.PodAlert, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.PodAlerts(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.PodAlerts(cur.Namespace).Update(transform(cur.DeepCopy()))
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to update PodAlert %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to update PodAlert %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}
