package resource

import (
	"errors"

	api "github.com/appscode/kubed/pkg/apis/kubed/v1alpha1"
	"github.com/blevesearch/bleve"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/request"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
)

var _ rest.Getter = &ResourceIndexer{}
var _ rest.GroupVersionKindProvider = &ResourceIndexer{}

func (ri *ResourceIndexer) NewREST() rest.Storage {
	return ri
}

func (ri *ResourceIndexer) New() runtime.Object {
	return &api.Stuff{}
}

func (ri *ResourceIndexer) GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind {
	return api.SchemeGroupVersion.WithKind("Stuff")
}

func (ri *ResourceIndexer) Get(ctx apirequest.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	ns, ok := request.NamespaceFrom(ctx)
	if !ok {
		return nil, errors.New("missing namespace")
	}
	if len(name) == 0 {
		return nil, errors.New("missing search query")
	}

	index, err := ri.indexFor(ns)
	if err != nil {
		return nil, err
	}

	req := bleve.NewSearchRequest(bleve.NewMatchQuery(name))
	result, err := index.Search(req)
	if err != nil {
		return nil, err
	}

	resp := &api.Stuff{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "kubed.appscode.com/v1alpha1",
			Kind:       "Stuff",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Hits:     make([]api.ResultEntry, 0, result.Total),
		Total:    result.Total,
		MaxScore: result.MaxScore,
		Took:     metav1.Duration{Duration: result.Took},
	}

	for _, hit := range result.Hits {
		raw, err := index.GetInternal([]byte(hit.ID))
		if err != nil {
			// log.Errorf("failed to load document with id %s. Reason: %s", hit.ID, err)
			continue
		}
		resp.Hits = append(resp.Hits, api.ResultEntry{
			Object: runtime.RawExtension{Raw: raw},
			Score:  hit.Score,
		})
	}
	return resp, nil
}
