/*
Copyright AppsCode Inc. and Contributors

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

package clusterid

import (
	"context"
	"flag"

	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

var clusterName = ""

func AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&clusterName, "cluster-name", clusterName, "Name of cluster used in a multi-cluster setup")
}

func AddGoFlags(fs *flag.FlagSet) {
	fs.StringVar(&clusterName, "cluster-name", clusterName, "Name of cluster used in a multi-cluster setup")
}

func ClusterName() string {
	return clusterName
}

func ClusterUID(client corev1.NamespaceInterface) (string, error) {
	ns, err := client.Get(context.TODO(), metav1.NamespaceSystem, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return string(ns.UID), nil
}
