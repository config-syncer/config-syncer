/*
Copyright The Config Syncer Authors.

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

package syncer

import (
	"context"
	"strings"

	"gomodules.xyz/pointer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	"kmodules.xyz/client-go/meta"
)

type SyncOptions struct {
	NamespaceSelector *string // if nil, delete from cluster
	Contexts          sets.String
}

func GetSyncOptions(annotations map[string]string) SyncOptions {
	opts := SyncOptions{}
	if v, err := meta.GetStringValue(annotations, ConfigSyncKey); err == nil {
		if v == "true" {
			opts.NamespaceSelector = pointer.StringP(labels.Everything().String())
		} else {
			opts.NamespaceSelector = &v
		}
	}
	if contexts, _ := meta.GetStringValue(annotations, ConfigSyncContexts); contexts != "" {
		opts.Contexts = sets.NewString(strings.Split(contexts, ",")...)
	}
	return opts
}

func NamespacesForSelector(kc kubernetes.Interface, selector string) (sets.String, error) {
	namespaces, err := kc.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
		LabelSelector: selector,
	})
	if err != nil {
		return nil, err
	}
	ns := sets.NewString()
	for _, obj := range namespaces.Items {
		ns.Insert(obj.Name)
	}
	return ns, nil
}
