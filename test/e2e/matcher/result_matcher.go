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

package matcher

import (
	"encoding/json"
	"fmt"
	"reflect"

	api "github.com/appscode/kubed/apis/kubed/v1alpha1"

	"github.com/onsi/gomega/types"
	"k8s.io/apimachinery/pkg/runtime"
	wapi "kmodules.xyz/webhook-runtime/apis/workload/v1"
	wcs "kmodules.xyz/webhook-runtime/client/workload/v1"
)

func HaveObject(expected runtime.Object) types.GomegaMatcher {
	return &objectMatcher{
		expected: expected,
	}
}

type objectMatcher struct {
	expected runtime.Object
}

func (matcher *objectMatcher) Match(actual interface{}) (success bool, err error) {
	result := actual.(*api.SearchResult)
	hits := result.Hits

	expected, err := wcs.ConvertToWorkload(matcher.expected)
	if err != nil {
		return false, err
	}

	for _, hit := range hits {
		found := &wapi.Workload{}
		err = json.Unmarshal(hit.Object.Raw, &found)
		if err != nil {
			continue
		}
		success, err = MatchWorkload(expected, found)
		if success {
			return success, nil
		}
	}
	return false, err
}

func (matcher *objectMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\n to  contain \n\t%#v", actual, matcher.expected)
}

func (matcher *objectMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\n not to contain \n\t%#v", actual, matcher.expected)
}

func MatchWorkload(expected, found *wapi.Workload) (success bool, err error) {
	if expected.Name != found.Name || expected.Namespace != found.Namespace {
		return false, nil
	}
	if !reflect.DeepEqual(expected.Spec, found.Spec) {
		return false, nil
	}

	return true, nil
}
