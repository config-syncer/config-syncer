package indexers

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/appscode/go/errors"
	"github.com/appscode/go/log"
	"github.com/appscode/pat"
	"github.com/blevesearch/bleve"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/tools/cache"
)

const (
	// Default router prefix of indexers
	indexerHTTPRouterPrefix = "/search"
)

type ResourceIndexer struct {
	// Full text indexer client
	index bleve.Index

	enable bool
	lock   sync.RWMutex
}

var _ cache.ResourceEventHandler = &ResourceIndexer{}

func NewResourceIndexer(dst string) (*ResourceIndexer, error) {
	idx, err := ensureIndex(filepath.Join(dst, "resource.indexer"), "search")
	if err != nil {
		return nil, err
	}
	return &ResourceIndexer{index: idx}, nil
}

func (ri *ResourceIndexer) Configure(enable bool) error {
	ri.lock.Lock()
	defer ri.lock.Unlock()

	ri.enable = enable
	return nil
}

func (ri *ResourceIndexer) OnAdd(obj interface{}) {
	ri.lock.RLock()
	defer ri.lock.RUnlock()

	if !ri.enable {
		return
	}

	if err := ri.indexDocument(obj); err != nil {
		log.Errorln(err)
	}
}

func (ri *ResourceIndexer) OnDelete(obj interface{}) {
	ri.lock.RLock()
	defer ri.lock.RUnlock()

	if !ri.enable {
		return
	}

	key := keyFunction(obj)
	if err := ri.index.Delete(key); err != nil {
		log.Errorln(err)
		return
	}
	if err := ri.index.DeleteInternal([]byte(key)); err != nil {
		log.Errorln(err)
		return
	}
}

func (ri *ResourceIndexer) OnUpdate(oldObj, newObj interface{}) {
	ri.lock.RLock()
	defer ri.lock.RUnlock()

	if !ri.enable {
		return
	}

	if err := ri.indexDocument(newObj); err != nil {
		log.Errorln(err)
	}
}

func (ri *ResourceIndexer) indexDocument(obj interface{}) error {
	key := keyFunction(obj)
	err := ri.index.Index(key, obj)
	if err != nil {
		return errors.FromErr(err).WithMessage("Failed to index document").Err()
	}

	data, err := json.Marshal(obj)
	if err != nil {
		return errors.FromErr(err).WithMessage("Failed to marshal internal document").Err()
	}

	err = ri.index.SetInternal([]byte(key), data)
	if err != nil {
		return errors.FromErr(err).WithMessage("Failed store internal document").Err()
	}
	return nil
}

func keyFunction(obj interface{}) string {
	if accessor, err := meta.Accessor(obj); err == nil {
		return string(accessor.GetUID())
	}
	return ""
}

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
		result, err := ri.index.Search(search)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		searchResult := &SearchResults{SearchResult: result, Results: make([]json.RawMessage, 0)}
		for _, hit := range result.Hits {
			raw, err := ri.index.GetInternal([]byte(hit.ID))
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
