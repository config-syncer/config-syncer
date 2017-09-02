package clientset

import (
	aci "github.com/k8sdb/apimachinery/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

type ElasticsearchNamespacer interface {
	Elasticsearches(namespace string) ElasticsearchInterface
}

type ElasticsearchInterface interface {
	List(opts metav1.ListOptions) (*aci.ElasticsearchList, error)
	Get(name string) (*aci.Elasticsearch, error)
	Create(elastic *aci.Elasticsearch) (*aci.Elasticsearch, error)
	Update(elastic *aci.Elasticsearch) (*aci.Elasticsearch, error)
	Delete(name string) error
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	UpdateStatus(elastic *aci.Elasticsearch) (*aci.Elasticsearch, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (*aci.Elasticsearch, error)
}

type ElasticsearchImpl struct {
	r  rest.Interface
	ns string
}

var _ ElasticsearchInterface = &ElasticsearchImpl{}

func newElastic(c *ExtensionClient, namespace string) *ElasticsearchImpl {
	return &ElasticsearchImpl{c.restClient, namespace}
}

func (c *ElasticsearchImpl) List(opts metav1.ListOptions) (result *aci.ElasticsearchList, err error) {
	result = &aci.ElasticsearchList{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource(aci.ResourceTypeElasticsearch).
		VersionedParams(&opts, ExtendedCodec).
		Do().
		Into(result)
	return
}

func (c *ElasticsearchImpl) Get(name string) (result *aci.Elasticsearch, err error) {
	result = &aci.Elasticsearch{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource(aci.ResourceTypeElasticsearch).
		Name(name).
		Do().
		Into(result)
	return
}

func (c *ElasticsearchImpl) Create(elastic *aci.Elasticsearch) (result *aci.Elasticsearch, err error) {
	result = &aci.Elasticsearch{}
	err = c.r.Post().
		Namespace(c.ns).
		Resource(aci.ResourceTypeElasticsearch).
		Body(elastic).
		Do().
		Into(result)
	return
}

func (c *ElasticsearchImpl) Update(elastic *aci.Elasticsearch) (result *aci.Elasticsearch, err error) {
	result = &aci.Elasticsearch{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource(aci.ResourceTypeElasticsearch).
		Name(elastic.Name).
		Body(elastic).
		Do().
		Into(result)
	return
}

func (c *ElasticsearchImpl) Delete(name string) (err error) {
	return c.r.Delete().
		Namespace(c.ns).
		Resource(aci.ResourceTypeElasticsearch).
		Name(name).
		Do().
		Error()
}

func (c *ElasticsearchImpl) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return c.r.Get().
		Prefix("watch").
		Namespace(c.ns).
		Resource(aci.ResourceTypeElasticsearch).
		VersionedParams(&opts, ExtendedCodec).
		Watch()
}

func (c *ElasticsearchImpl) UpdateStatus(elastic *aci.Elasticsearch) (result *aci.Elasticsearch, err error) {
	result = &aci.Elasticsearch{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource(aci.ResourceTypeElasticsearch).
		Name(elastic.Name).
		SubResource("status").
		Body(elastic).
		Do().
		Into(result)
	return
}

func (c *ElasticsearchImpl) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *aci.Elasticsearch, err error) {
	result = &aci.Elasticsearch{}
	err = c.r.Patch(pt).
		Namespace(c.ns).
		Resource(aci.ResourceTypeElasticsearch).
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
