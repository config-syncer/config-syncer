package indexers

import (
	"encoding/json"
	"net/http"

	"github.com/appscode/log"
	"github.com/appscode/pat"
	"github.com/blevesearch/bleve"
)

const (
	// Default router prefix of indexers
	indexerHTTPRouterPrefix = "/search"
)

func (ri *ResourceIndexer) RegisterRouters(r *pat.PatternServeMux) {
	// Example format /index/namespaces/kube-system/pods/kube-dns
	r.Get(indexerHTTPRouterPrefix, http.HandlerFunc(ri.ServeHTTP))
}

type SearchResults struct {
	*bleve.SearchResult `json:",inline"`
	Results             []json.RawMessage `json:"results"`
}

func (ri *ResourceIndexer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Infoln("Request received at", req.URL.Path)

	queryString := req.URL.Query().Get("q")
	if len(queryString) > 0 {
		log.Infoln("Query received", queryString)
		q := bleve.NewMatchQuery(queryString)
		search := bleve.NewSearchRequest(q)
		result, err := ri.client.Search(search)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		searchResult := &SearchResults{SearchResult: result, Results: make([]json.RawMessage, 0)}
		for _, hit := range result.Hits {
			raw, err := ri.client.GetInternal([]byte(hit.ID))
			if err != nil {
				log.Errorln("Failed to get internal result", err)
				continue
			}
			searchResult.Results = append(searchResult.Results, json.RawMessage(raw))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(searchResult)
		return
	}
	http.Error(w, "Bad Request", http.StatusBadRequest)
}
