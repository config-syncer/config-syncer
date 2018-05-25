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
