package indexers

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/appscode/log"
	"github.com/appscode/pat"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Default router prefix of indexers
	httpRouterPrefix = "/index"
)

func (ri *ReverseIndexer) RegisterRouters(r *pat.PatternServeMux) {
	// Example format /index/namespaces/kube-system/pods/kube-dns
	r.Get(httpRouterPrefix+"/namespaces/:namespace/:resource/:name", http.HandlerFunc(ri.ServeHTTP))
}

func (ri *ReverseIndexer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Infoln("Request received at", req.URL.Path)
	params, found := pat.FromContext(req.Context())
	if !found {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}
	resource := params.Get(":resource")

	// TODO (@sadlil) Use mapping
	switch resource {
	case "pod", "pods":
		ri.servePodIndex(w, req, params)
		return
	}

	http.Error(w, "Resource not supported", http.StatusNotImplemented)
}

func (ri *ReverseIndexer) servePodIndex(w http.ResponseWriter, req *http.Request, params url.Values) {
	namespace, name := params.Get(":namespace"), params.Get(":name")
	if len(namespace) > 0 && len(name) > 0 {
		key := namespacerKey(v1.ObjectMeta{Name: name, Namespace: namespace})
		if val, ok := ri.podToServiceRecordMap[key]; ok {
			if err := json.NewEncoder(w).Encode(val); err == nil {
				w.Header().Set("Content-Type", "application/json")
				return
			} else {
				http.Error(w, "Server error"+err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.NotFound(w, req)
		}
		return
	}
	http.Error(w, "Bad Request", http.StatusBadRequest)
}
