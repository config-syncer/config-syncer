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

func EnsurePostgres(c tcs.KubedbV1alpha1Interface, meta metav1.ObjectMeta, transform func(alert *aci.Postgres) *aci.Postgres) (*aci.Postgres, error) {
	return CreateOrPatchPostgres(c, meta, transform)
}

func CreateOrPatchPostgres(c tcs.KubedbV1alpha1Interface, meta metav1.ObjectMeta, transform func(alert *aci.Postgres) *aci.Postgres) (*aci.Postgres, error) {
	cur, err := c.Postgreses(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		glog.V(3).Infof("Creating Postgres %s/%s.", meta.Namespace, meta.Name)
		return c.Postgreses(meta.Namespace).Create(transform(&aci.Postgres{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Postgres",
				APIVersion: aci.SchemeGroupVersion.String(),
			},
			ObjectMeta: meta,
		}))
	} else if err != nil {
		return nil, err
	}
	return PatchPostgres(c, cur, transform)
}

func PatchPostgres(c tcs.KubedbV1alpha1Interface, cur *aci.Postgres, transform func(*aci.Postgres) *aci.Postgres) (*aci.Postgres, error) {
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
	glog.V(3).Infof("Patching Postgres %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	result, err := c.Postgreses(cur.Namespace).Patch(cur.Name, types.MergePatchType, patch)
	return result, err
}

func TryPatchPostgres(c tcs.KubedbV1alpha1Interface, meta metav1.ObjectMeta, transform func(*aci.Postgres) *aci.Postgres) (result *aci.Postgres, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.Postgreses(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = PatchPostgres(c, cur, transform)
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to patch Postgres %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to patch Postgres %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}

func TryUpdatePostgres(c tcs.KubedbV1alpha1Interface, meta metav1.ObjectMeta, transform func(*aci.Postgres) *aci.Postgres) (result *aci.Postgres, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.Postgreses(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.Postgreses(cur.Namespace).Update(transform(cur))
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to update Postgres %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to update Postgres %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}
