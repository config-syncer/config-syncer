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
