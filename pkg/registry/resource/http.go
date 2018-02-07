package resource

import (
	"encoding/json"
	"net/http"

	"github.com/appscode/go/log"
	"github.com/appscode/pat"
	"github.com/blevesearch/bleve"
)

const (
	// Default router prefix of indexers
	indexerHTTPRouterPrefix = "/search"
)

func (ri *Indexer) RegisterRouters(r *pat.PatternServeMux) {
	// Example format /index/namespaces/kube-system/pods/kube-dns
	r.Get(indexerHTTPRouterPrefix, http.HandlerFunc(ri.ServeHTTP))
}

type SearchResults struct {
	*bleve.SearchResult `json:",inline"`
	Results             []json.RawMessage `json:"results"`
}

func (ri *Indexer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ns := r.URL.Query().Get("ns")
	qs := r.URL.Query().Get("q")

	if len(qs) == 0 {
		http.Error(w, "missing search query (q=)", http.StatusBadRequest)
		return
	}
	if len(ns) == 0 {
		http.Error(w, "missing namespace (ns=)", http.StatusBadRequest)
		return
	}

	index, err := ri.indexFor(ns)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req := bleve.NewSearchRequest(bleve.NewMatchQuery(qs))
	result, err := index.Search(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := &SearchResults{
		SearchResult: result,
		Results:      make([]json.RawMessage, 0, result.Total),
	}
	for _, hit := range result.Hits {
		raw, err := index.GetInternal([]byte(hit.ID))
		if err != nil {
			log.Errorf("failed to load document with id %s. Reason: %s", hit.ID, err)
			continue
		}
		resp.Results = append(resp.Results, json.RawMessage(raw))
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
