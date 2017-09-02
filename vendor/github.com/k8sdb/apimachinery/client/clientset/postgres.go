package clientset

import (
	aci "github.com/k8sdb/apimachinery/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

type PostgresNamespacer interface {
	Postgreses(namespace string) PostgresInterface
}

type PostgresInterface interface {
	List(opts metav1.ListOptions) (*aci.PostgresList, error)
	Get(name string) (*aci.Postgres, error)
	Create(postgres *aci.Postgres) (*aci.Postgres, error)
	Update(postgres *aci.Postgres) (*aci.Postgres, error)
	Delete(name string) error
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	UpdateStatus(postgres *aci.Postgres) (*aci.Postgres, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (*aci.Postgres, error)
}

type PostgresImpl struct {
	r  rest.Interface
	ns string
}

var _ PostgresInterface = &PostgresImpl{}

func newPostgres(c *ExtensionClient, namespace string) *PostgresImpl {
	return &PostgresImpl{c.restClient, namespace}
}

func (c *PostgresImpl) List(opts metav1.ListOptions) (result *aci.PostgresList, err error) {
	result = &aci.PostgresList{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource(aci.ResourceTypePostgres).
		VersionedParams(&opts, ExtendedCodec).
		Do().
		Into(result)
	return
}

func (c *PostgresImpl) Get(name string) (result *aci.Postgres, err error) {
	result = &aci.Postgres{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource(aci.ResourceTypePostgres).
		Name(name).
		Do().
		Into(result)
	return
}

func (c *PostgresImpl) Create(postgres *aci.Postgres) (result *aci.Postgres, err error) {
	result = &aci.Postgres{}
	err = c.r.Post().
		Namespace(c.ns).
		Resource(aci.ResourceTypePostgres).
		Body(postgres).
		Do().
		Into(result)
	return
}

func (c *PostgresImpl) Update(postgres *aci.Postgres) (result *aci.Postgres, err error) {
	result = &aci.Postgres{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource(aci.ResourceTypePostgres).
		Name(postgres.Name).
		Body(postgres).
		Do().
		Into(result)
	return
}

func (c *PostgresImpl) Delete(name string) (err error) {
	return c.r.Delete().
		Namespace(c.ns).
		Resource(aci.ResourceTypePostgres).
		Name(name).
		Do().
		Error()
}

func (c *PostgresImpl) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return c.r.Get().
		Prefix("watch").
		Namespace(c.ns).
		Resource(aci.ResourceTypePostgres).
		VersionedParams(&opts, ExtendedCodec).
		Watch()
}

func (c *PostgresImpl) UpdateStatus(postgres *aci.Postgres) (result *aci.Postgres, err error) {
	result = &aci.Postgres{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource(aci.ResourceTypePostgres).
		Name(postgres.Name).
		SubResource("status").
		Body(postgres).
		Do().
		Into(result)
	return
}

func (c *PostgresImpl) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *aci.Postgres, err error) {
	result = &aci.Postgres{}
	err = c.r.Patch(pt).
		Namespace(c.ns).
		Resource(aci.ResourceTypePostgres).
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
