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
	"net/http"

	. "github.com/onsi/gomega"
)

func (f *Invocation) EventuallyReverseIndex(path string) GomegaAsyncAssertion {
	request, err := http.NewRequest(http.MethodGet, "http://localhost:8080"+path, nil)
	Expect(err).NotTo(HaveOccurred())
	return Eventually(func() int {
		resp, err := http.DefaultClient.Do(request)
		Expect(err).NotTo(HaveOccurred())
		return resp.StatusCode
	}, DefaultEventuallyTimeout, DefaultEventuallyPollingInterval)
}
