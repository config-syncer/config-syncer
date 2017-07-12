package clientset

import (
	aci "github.com/k8sdb/apimachinery/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

type ElasticNamespacer interface {
	Elastics(namespace string) ElasticInterface
}

type ElasticInterface interface {
	List(opts metav1.ListOptions) (*aci.ElasticList, error)
	Get(name string) (*aci.Elastic, error)
	Create(elastic *aci.Elastic) (*aci.Elastic, error)
	Update(elastic *aci.Elastic) (*aci.Elastic, error)
	Delete(name string) error
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	UpdateStatus(elastic *aci.Elastic) (*aci.Elastic, error)
}

type ElasticImpl struct {
	r  rest.Interface
	ns string
}

var _ ElasticInterface = &ElasticImpl{}

func newElastic(c *ExtensionClient, namespace string) *ElasticImpl {
	return &ElasticImpl{c.restClient, namespace}
}

func (c *ElasticImpl) List(opts metav1.ListOptions) (result *aci.ElasticList, err error) {
	result = &aci.ElasticList{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource(aci.ResourceTypeElastic).
		VersionedParams(&opts, ExtendedCodec).
		Do().
		Into(result)
	return
}

func (c *ElasticImpl) Get(name string) (result *aci.Elastic, err error) {
	result = &aci.Elastic{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource(aci.ResourceTypeElastic).
		Name(name).
		Do().
		Into(result)
	return
}

func (c *ElasticImpl) Create(elastic *aci.Elastic) (result *aci.Elastic, err error) {
	result = &aci.Elastic{}
	err = c.r.Post().
		Namespace(c.ns).
		Resource(aci.ResourceTypeElastic).
		Body(elastic).
		Do().
		Into(result)
	return
}

func (c *ElasticImpl) Update(elastic *aci.Elastic) (result *aci.Elastic, err error) {
	result = &aci.Elastic{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource(aci.ResourceTypeElastic).
		Name(elastic.Name).
		Body(elastic).
		Do().
		Into(result)
	return
}

func (c *ElasticImpl) Delete(name string) (err error) {
	return c.r.Delete().
		Namespace(c.ns).
		Resource(aci.ResourceTypeElastic).
		Name(name).
		Do().
		Error()
}

func (c *ElasticImpl) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return c.r.Get().
		Prefix("watch").
		Namespace(c.ns).
		Resource(aci.ResourceTypeElastic).
		VersionedParams(&opts, ExtendedCodec).
		Watch()
}

func (c *ElasticImpl) UpdateStatus(elastic *aci.Elastic) (result *aci.Elastic, err error) {
	result = &aci.Elastic{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource(aci.ResourceTypeElastic).
		Name(elastic.Name).
		SubResource("status").
		Body(elastic).
		Do().
		Into(result)
	return
}
