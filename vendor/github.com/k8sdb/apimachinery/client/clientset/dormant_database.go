package clientset

import (
	aci "github.com/k8sdb/apimachinery/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

type DormantDatabaseNamespacer interface {
	DormantDatabases(namespace string) DormantDatabaseInterface
}

type DormantDatabaseInterface interface {
	List(opts metav1.ListOptions) (*aci.DormantDatabaseList, error)
	Get(name string) (*aci.DormantDatabase, error)
	Create(drmn *aci.DormantDatabase) (*aci.DormantDatabase, error)
	Update(drmn *aci.DormantDatabase) (*aci.DormantDatabase, error)
	Delete(name string) error
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	UpdateStatus(drmn *aci.DormantDatabase) (*aci.DormantDatabase, error)
}

type DormantDatabaseImpl struct {
	r  rest.Interface
	ns string
}

var _ DormantDatabaseInterface = &DormantDatabaseImpl{}

func newDormantDatabase(c *ExtensionClient, namespace string) *DormantDatabaseImpl {
	return &DormantDatabaseImpl{c.restClient, namespace}
}

func (c *DormantDatabaseImpl) List(opts metav1.ListOptions) (result *aci.DormantDatabaseList, err error) {
	result = &aci.DormantDatabaseList{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource(aci.ResourceTypeDormantDatabase).
		VersionedParams(&opts, ExtendedCodec).
		Do().
		Into(result)
	return
}

func (c *DormantDatabaseImpl) Get(name string) (result *aci.DormantDatabase, err error) {
	result = &aci.DormantDatabase{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource(aci.ResourceTypeDormantDatabase).
		Name(name).
		Do().
		Into(result)
	return
}

func (c *DormantDatabaseImpl) Create(drmn *aci.DormantDatabase) (result *aci.DormantDatabase, err error) {
	result = &aci.DormantDatabase{}
	err = c.r.Post().
		Namespace(c.ns).
		Resource(aci.ResourceTypeDormantDatabase).
		Body(drmn).
		Do().
		Into(result)
	return
}

func (c *DormantDatabaseImpl) Update(drmn *aci.DormantDatabase) (result *aci.DormantDatabase, err error) {
	result = &aci.DormantDatabase{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource(aci.ResourceTypeDormantDatabase).
		Name(drmn.Name).
		Body(drmn).
		Do().
		Into(result)
	return
}

func (c *DormantDatabaseImpl) Delete(name string) (err error) {
	return c.r.Delete().
		Namespace(c.ns).
		Resource(aci.ResourceTypeDormantDatabase).
		Name(name).
		Do().
		Error()
}

func (c *DormantDatabaseImpl) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return c.r.Get().
		Prefix("watch").
		Namespace(c.ns).
		Resource(aci.ResourceTypeDormantDatabase).
		VersionedParams(&opts, ExtendedCodec).
		Watch()
}

func (c *DormantDatabaseImpl) UpdateStatus(drmn *aci.DormantDatabase) (result *aci.DormantDatabase, err error) {
	result = &aci.DormantDatabase{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource(aci.ResourceTypeDormantDatabase).
		Name(drmn.Name).
		SubResource("status").
		Body(drmn).
		Do().
		Into(result)
	return
}
