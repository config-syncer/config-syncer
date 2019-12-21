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
	"fmt"
	"reflect"

	"github.com/onsi/gomega/types"
	core "k8s.io/api/core/v1"
)

func BeEquivalentToConfigMap(expected *core.ConfigMap) types.GomegaMatcher {
	return &configMapMatcher{
		expected: expected,
	}
}

type configMapMatcher struct {
	expected *core.ConfigMap
}

func (matcher *configMapMatcher) Match(actual interface{}) (success bool, err error) {
	found := actual.(*core.ConfigMap)

	if matcher.expected.Name != found.Name {
		return false, err
	}

	if matcher.expected.Namespace != found.Namespace {
		return false, err
	}

	if !reflect.DeepEqual(matcher.expected.Data, found.Data) {
		return false, err
	}
	return true, nil
}

func (matcher *configMapMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\n to  be equvalent to \n\t%#v", actual, matcher.expected)
}

func (matcher *configMapMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\n not to be equivalent to \n\t%#v", actual, matcher.expected)
}
