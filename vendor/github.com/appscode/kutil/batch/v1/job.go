package v1

import (
	"encoding/json"
	"fmt"

	"github.com/appscode/kutil"
	"github.com/golang/glog"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	batch "k8s.io/client-go/pkg/apis/batch/v1"
)

func EnsureJob(c clientset.Interface, meta metav1.ObjectMeta, transform func(*batch.Job) *batch.Job) (*batch.Job, error) {
	return CreateOrPatchJob(c, meta, transform)
}

func CreateOrPatchJob(c clientset.Interface, meta metav1.ObjectMeta, transform func(*batch.Job) *batch.Job) (*batch.Job, error) {
	cur, err := c.BatchV1().Jobs(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		return c.BatchV1().Jobs(meta.Namespace).Create(transform(&batch.Job{ObjectMeta: meta}))
	} else if err != nil {
		return nil, err
	}
	return PatchJob(c, cur, transform)
}

func PatchJob(c clientset.Interface, cur *batch.Job, transform func(*batch.Job) *batch.Job) (*batch.Job, error) {
	curJson, err := json.Marshal(cur)
	if err != nil {
		return nil, err
	}

	modJson, err := json.Marshal(transform(cur))
	if err != nil {
		return nil, err
	}

	patch, err := strategicpatch.CreateTwoWayMergePatch(curJson, modJson, batch.Job{})
	if err != nil {
		return nil, err
	}
	if len(patch) == 0 || string(patch) == "{}" {
		return cur, nil
	}
	glog.V(5).Infof("Patching Job %s@%s with %s.", cur.Name, cur.Namespace, string(patch))
	return c.BatchV1().Jobs(cur.Namespace).Patch(cur.Name, types.StrategicMergePatchType, patch)
}

func TryPatchJob(c clientset.Interface, meta metav1.ObjectMeta, transform func(*batch.Job) *batch.Job) (result *batch.Job, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.BatchV1().Jobs(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = PatchJob(c, cur, transform)
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to patch Job %s@%s due to %v.", attempt, cur.Name, cur.Namespace, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to patch Job %s@%s after %d attempts due to %v", meta.Name, meta.Namespace, attempt, err)
	}
	return
}

func TryUpdateJob(c clientset.Interface, meta metav1.ObjectMeta, transform func(*batch.Job) *batch.Job) (result *batch.Job, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.BatchV1().Jobs(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.BatchV1().Jobs(cur.Namespace).Update(transform(cur))
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to update Job %s@%s due to %v.", attempt, cur.Name, cur.Namespace, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to update Job %s@%s after %d attempts due to %v", meta.Name, meta.Namespace, attempt, err)
	}
	return
}
