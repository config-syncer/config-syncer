package indexers

import "net/http"

func (ri *ReverseIndexer) Handlers() http.Handler {
	return ri.apiHandler
}

type reverseIndexAPIHandlers struct{}

func (ri *reverseIndexAPIHandlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
