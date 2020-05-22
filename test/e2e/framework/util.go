/*
Copyright The Kubed Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package framework

import (
	"context"
	"fmt"

	"github.com/appscode/go/types"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	kutil "kmodules.xyz/client-go"
	"kmodules.xyz/client-go/tools/exec"
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
	return wait.PollImmediate(kutil.RetryInterval, kutil.ReadinessTimeout, func() (done bool, err error) {
		if obj, err := fi.KubeClient.AppsV1().Deployments(meta.Namespace).Get(context.TODO(), meta.Name, metav1.GetOptions{}); err == nil {
			return types.Int32(obj.Spec.Replicas) == obj.Status.ReadyReplicas, nil
		}
		return false, nil
	})
}

func (fi *Invocation) WaitUntilDeploymentTerminated(meta metav1.ObjectMeta) error {
	return wait.PollImmediate(kutil.RetryInterval, kutil.GCTimeout, func() (done bool, err error) {
		if pods, err := fi.KubeClient.CoreV1().Pods(meta.Namespace).List(context.TODO(), metav1.ListOptions{}); err == nil {
			return len(pods.Items) == 0, nil
		}
		return false, nil
	})
}

func (fi *Invocation) RemoveFromOperatorPod(dir string) error {
	pod, err := fi.OperatorPod()
	if err != nil {
		return err
	}

	_, err = exec.ExecIntoPod(fi.ClientConfig, pod, exec.Command("rm", "-rf", dir))
	if err != nil {
		return err
	}

	return nil
}

func (fi *Invocation) OperatorPod() (*core.Pod, error) {
	pods, err := fi.KubeClient.CoreV1().Pods(OperatorNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, pod := range pods.Items {
		for _, c := range pod.Spec.Containers {
			if c.Name == ContainerOperator {
				return &pod, nil
			}
		}
	}

	return nil, fmt.Errorf("pod not found")
}

func (fi *Invocation) DeleteService(meta metav1.ObjectMeta) error {
	return fi.KubeClient.CoreV1().Services(meta.Namespace).Delete(context.TODO(), meta.Name, *deleteInBackground())
}

func (fi *Invocation) DeleteEndpoints(meta metav1.ObjectMeta) error {
	return fi.KubeClient.CoreV1().Endpoints(meta.Namespace).Delete(context.TODO(), meta.Name, *deleteInBackground())
}
