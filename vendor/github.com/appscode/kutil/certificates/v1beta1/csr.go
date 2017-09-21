package v1beta1

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
	certificates "k8s.io/client-go/pkg/apis/certificates/v1beta1"
)

func EnsureCSR(c clientset.Interface, meta metav1.ObjectMeta, transform func(*certificates.CertificateSigningRequest) *certificates.CertificateSigningRequest) (*certificates.CertificateSigningRequest, error) {
	return CreateOrPatchCSR(c, meta, transform)
}

func CreateOrPatchCSR(c clientset.Interface, meta metav1.ObjectMeta, transform func(*certificates.CertificateSigningRequest) *certificates.CertificateSigningRequest) (*certificates.CertificateSigningRequest, error) {
	cur, err := c.CertificatesV1beta1().CertificateSigningRequests().Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		return c.CertificatesV1beta1().CertificateSigningRequests().Create(transform(&certificates.CertificateSigningRequest{ObjectMeta: meta}))
	} else if err != nil {
		return nil, err
	}
	return PatchCSR(c, cur, transform)
}

func PatchCSR(c clientset.Interface, cur *certificates.CertificateSigningRequest, transform func(*certificates.CertificateSigningRequest) *certificates.CertificateSigningRequest) (*certificates.CertificateSigningRequest, error) {
	curJson, err := json.Marshal(cur)
	if err != nil {
		return nil, err
	}

	modJson, err := json.Marshal(transform(cur))
	if err != nil {
		return nil, err
	}

	patch, err := strategicpatch.CreateTwoWayMergePatch(curJson, modJson, certificates.CertificateSigningRequest{})
	if err != nil {
		return nil, err
	}
	if len(patch) == 0 || string(patch) == "{}" {
		return cur, nil
	}
	glog.V(5).Infof("Patching CertificateSigningRequest %s@%s with %s.", cur.Name, cur.Namespace, string(patch))
	return c.CertificatesV1beta1().CertificateSigningRequests().Patch(cur.Name, types.StrategicMergePatchType, patch)
}

func TryPatchCSR(c clientset.Interface, meta metav1.ObjectMeta, transform func(*certificates.CertificateSigningRequest) *certificates.CertificateSigningRequest) (result *certificates.CertificateSigningRequest, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.CertificatesV1beta1().CertificateSigningRequests().Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = PatchCSR(c, cur, transform)
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to patch CertificateSigningRequest %s@%s due to %v.", attempt, cur.Name, cur.Namespace, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to patch CertificateSigningRequest %s@%s after %d attempts due to %v", meta.Name, meta.Namespace, attempt, err)
	}
	return
}

func TryUpdateCSR(c clientset.Interface, meta metav1.ObjectMeta, transform func(*certificates.CertificateSigningRequest) *certificates.CertificateSigningRequest) (result *certificates.CertificateSigningRequest, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.CertificatesV1beta1().CertificateSigningRequests().Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.CertificatesV1beta1().CertificateSigningRequests().Update(transform(cur))
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to update CertificateSigningRequest %s@%s due to %v.", attempt, cur.Name, cur.Namespace, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to update CertificateSigningRequest %s@%s after %d attempts due to %v", meta.Name, meta.Namespace, attempt, err)
	}
	return
}
