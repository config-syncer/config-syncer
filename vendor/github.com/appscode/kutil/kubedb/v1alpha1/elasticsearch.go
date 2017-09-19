package v1alpha1

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
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/apimachinery/pkg/util/wait"
)

func EnsureElasticsearch(c tcs.KubedbV1alpha1Interface, meta metav1.ObjectMeta, transform func(alert *aci.Elasticsearch) *aci.Elasticsearch) (*aci.Elasticsearch, error) {
	return CreateOrPatchElasticsearch(c, meta, transform)
}

func CreateOrPatchElasticsearch(c tcs.KubedbV1alpha1Interface, meta metav1.ObjectMeta, transform func(alert *aci.Elasticsearch) *aci.Elasticsearch) (*aci.Elasticsearch, error) {
	cur, err := c.Elasticsearchs(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		return c.Elasticsearchs(meta.Namespace).Create(transform(&aci.Elasticsearch{ObjectMeta: meta}))
	} else if err != nil {
		return nil, err
	}
	return PatchElasticsearch(c, cur, transform)
}

func PatchElasticsearch(c tcs.KubedbV1alpha1Interface, cur *aci.Elasticsearch, transform func(*aci.Elasticsearch) *aci.Elasticsearch) (*aci.Elasticsearch, error) {
	curJson, err := json.Marshal(cur)
	if err != nil {
		return nil, err
	}

	modJson, err := json.Marshal(transform(cur))
	if err != nil {
		return nil, err
	}

	patch, err := strategicpatch.CreateTwoWayMergePatch(curJson, modJson, aci.Elasticsearch{})
	if err != nil {
		return nil, err
	}
	if len(patch) == 0 || string(patch) == "{}" {
		return cur, nil
	}
	glog.V(5).Infof("Patching Elasticsearch %s@%s with %s.", cur.Name, cur.Namespace, string(patch))
	result, err := c.Elasticsearchs(cur.Namespace).Patch(cur.Name, types.StrategicMergePatchType, patch)
	return result, err
}

func TryPatchElasticsearch(c tcs.KubedbV1alpha1Interface, meta metav1.ObjectMeta, transform func(*aci.Elasticsearch) *aci.Elasticsearch) (result *aci.Elasticsearch, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.Elasticsearchs(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = PatchElasticsearch(c, cur, transform)
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to patch Elasticsearch %s@%s due to %v.", attempt, cur.Name, cur.Namespace, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("Failed to patch Elasticsearch %s@%s after %d attempts due to %v", meta.Name, meta.Namespace, attempt, err)
	}
	return
}

func TryUpdateElasticsearch(c tcs.KubedbV1alpha1Interface, meta metav1.ObjectMeta, transform func(*aci.Elasticsearch) *aci.Elasticsearch) (result *aci.Elasticsearch, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.Elasticsearchs(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.Elasticsearchs(cur.Namespace).Update(transform(cur))
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to update Elasticsearch %s@%s due to %v.", attempt, cur.Name, cur.Namespace, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("Failed to update Elasticsearch %s@%s after %d attempts due to %v", meta.Name, meta.Namespace, attempt, err)
	}
	return
}
