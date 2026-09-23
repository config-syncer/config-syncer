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

package meta

import (
	"maps"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func LabelsForLabelSelector(sel *metav1.LabelSelector) (map[string]string, bool) {
	if sel == nil {
		return make(map[string]string), true
	}

	labels := make(map[string]string, len(sel.MatchLabels)+len(sel.MatchExpressions))
	maps.Copy(labels, sel.MatchLabels)

	for _, expr := range sel.MatchExpressions {
		switch expr.Operator {
		case metav1.LabelSelectorOpIn:
			if len(expr.Values) > 0 {
				labels[expr.Key] = expr.Values[0]
			}
		case metav1.LabelSelectorOpNotIn:
			if len(expr.Values) > 0 {
				v := expr.Values[0]
				if v == "true" && len(expr.Values) == 1 {
					labels[expr.Key] = "false"
				} else if v == "false" && len(expr.Values) == 1 {
					labels[expr.Key] = "true"
				} else {
					labels[expr.Key] = "not-" + v
				}
			}
		case metav1.LabelSelectorOpExists:
			labels[expr.Key] = ""
		case metav1.LabelSelectorOpDoesNotExist:
			delete(labels, expr.Key)
		}
	}
	return labels, len(sel.MatchExpressions) == 0
}
