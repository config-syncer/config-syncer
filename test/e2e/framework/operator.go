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
	api "github.com/appscode/kubed/apis/kubed/v1alpha1"

	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	kutil "kmodules.xyz/client-go"
)

const (
	OperatorName      = "kubed-operator"
	OperatorNamespace = "kube-system"
	ContainerOperator = "operator"
	OperatorConfig    = "kubed-config"
)

func (fi *Invocation) RestartKubedOperator(config *api.ClusterConfig) error {
	meta := metav1.ObjectMeta{
		Name:      OperatorConfig,
		Namespace: OperatorNamespace,
	}

	err := fi.DeleteSecret(meta)
	if err != nil && !kerr.IsNotFound(err) {
		return err
	}

	err = fi.WaitUntilSecretDeleted(meta)
	if err != nil {
		return err
	}

	kubeConfig, err := fi.KubeConfigSecret(config, meta)
	if err != nil {
		return err
	}

	_, err = fi.CreateSecret(kubeConfig)
	if err != nil {
		return err
	}

	pods, err := fi.KubeClient.CoreV1().Pods(OperatorNamespace).List(metav1.ListOptions{LabelSelector: "app=kubed"})
	if err != nil {
		return err
	}
	for _, pod := range pods.Items {
		for _, c := range pod.Spec.Containers {
			if c.Name == ContainerOperator {
				err = fi.KubeClient.CoreV1().Pods(OperatorNamespace).Delete(pod.Name, deleteInBackground())
				if err != nil {
					return err
				}
				err = fi.WaitUntilPodTerminated(pod.ObjectMeta)
				if err != nil {
					return err
				}
				break
			}
		}
	}

	deployment, err := fi.KubeClient.AppsV1().Deployments(OperatorNamespace).Get(OperatorName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	return fi.WaitUntilDeploymentReady(deployment.ObjectMeta)
}

func (fi *Invocation) WaitUntilPodTerminated(meta metav1.ObjectMeta) error {
	return wait.PollImmediate(kutil.RetryInterval, kutil.GCTimeout, func() (done bool, err error) {
		if _, err := fi.KubeClient.CoreV1().Pods(meta.Namespace).Get(meta.Name, metav1.GetOptions{}); err != nil {
			if kerr.IsNotFound(err) {
				return true, nil
			} else {
				return true, err
			}
		}
		return false, nil
	})
}
