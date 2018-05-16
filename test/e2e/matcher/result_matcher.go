package matcher

import (
	"encoding/json"
	"fmt"
	"reflect"

	api "github.com/appscode/kubed/apis/kubed/v1alpha1"
	wapi "github.com/appscode/kubernetes-webhook-util/apis/workload/v1"
	wcs "github.com/appscode/kubernetes-webhook-util/client/workload/v1"
	"github.com/onsi/gomega/types"
	"k8s.io/apimachinery/pkg/runtime"
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
