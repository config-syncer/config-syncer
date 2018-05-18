package framework

import (
	"path/filepath"

	"github.com/appscode/go/runtime"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"github.com/appscode/go/types"
)

var (
	KubedTestConfigFileDir = filepath.Join(runtime.GOPath(), "src", "github.com", "appscode", "kubed", "test", "e2e", "config.yaml")
)

func deleteInBackground() *metav1.DeleteOptions {
	policy := metav1.DeletePropagationBackground
	return &metav1.DeleteOptions{PropagationPolicy: &policy}
}

func deleteInForeground() *metav1.DeleteOptions {
	policy := metav1.DeletePropagationForeground
	return &metav1.DeleteOptions{PropagationPolicy: &policy}
}

func (fi *Invocation) WaitUntilDeploymentReady(meta metav1.ObjectMeta) error {
	return wait.PollImmediate(interval, timeout, func() (done bool, err error) {
		if obj, err := fi.KubeClient.AppsV1beta1().Deployments(meta.Namespace).Get(meta.Name, metav1.GetOptions{}); err == nil {
			return types.Int32(obj.Spec.Replicas) == obj.Status.ReadyReplicas, nil
		}
		return false, nil
	})
}

func (fi *Invocation) WaitUntilDeploymentTerminated(meta metav1.ObjectMeta) error {
	return wait.PollImmediate(interval, timeout, func() (done bool, err error) {
		if pods, err := fi.KubeClient.CoreV1().Pods(meta.Namespace).List(metav1.ListOptions{}); err == nil {
			return len(pods.Items) == 0, nil
		}
		return false, nil
	})
}
