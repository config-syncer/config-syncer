package clientset

import (
	aci "github.com/appscode/searchlight/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

type NodeAlertGetter interface {
	NodeAlerts(namespace string) NodeAlertInterface
}

type NodeAlertInterface interface {
	List(opts metav1.ListOptions) (*aci.NodeAlertList, error)
	Get(name string) (*aci.NodeAlert, error)
	Create(NodeAlert *aci.NodeAlert) (*aci.NodeAlert, error)
	Update(NodeAlert *aci.NodeAlert) (*aci.NodeAlert, error)
	Delete(name string) error
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	UpdateStatus(NodeAlert *aci.NodeAlert) (*aci.NodeAlert, error)
}

type NodeAlertImpl struct {
	r  rest.Interface
	ns string
}

var _ NodeAlertInterface = &NodeAlertImpl{}

func newNodeAlert(c *ExtensionClient, namespace string) *NodeAlertImpl {
	return &NodeAlertImpl{c.restClient, namespace}
}

func (c *NodeAlertImpl) List(opts metav1.ListOptions) (result *aci.NodeAlertList, err error) {
	result = &aci.NodeAlertList{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource(aci.ResourceTypeNodeAlert).
		VersionedParams(&opts, ExtendedCodec).
		Do().
		Into(result)
	return
}

func (c *NodeAlertImpl) Get(name string) (result *aci.NodeAlert, err error) {
	result = &aci.NodeAlert{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource(aci.ResourceTypeNodeAlert).
		Name(name).
		Do().
		Into(result)
	return
}

func (c *NodeAlertImpl) Create(alert *aci.NodeAlert) (result *aci.NodeAlert, err error) {
	result = &aci.NodeAlert{}
	err = c.r.Post().
		Namespace(c.ns).
		Resource(aci.ResourceTypeNodeAlert).
		Body(alert).
		Do().
		Into(result)
	return
}

func (c *NodeAlertImpl) Update(alert *aci.NodeAlert) (result *aci.NodeAlert, err error) {
	result = &aci.NodeAlert{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource(aci.ResourceTypeNodeAlert).
		Name(alert.Name).
		Body(alert).
		Do().
		Into(result)
	return
}

func (c *NodeAlertImpl) Delete(name string) (err error) {
	return c.r.Delete().
		Namespace(c.ns).
		Resource(aci.ResourceTypeNodeAlert).
		Name(name).
		Do().
		Error()
}

func (c *NodeAlertImpl) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return c.r.Get().
		Prefix("watch").
		Namespace(c.ns).
		Resource(aci.ResourceTypeNodeAlert).
		VersionedParams(&opts, ExtendedCodec).
		Watch()
}

func (c *NodeAlertImpl) UpdateStatus(alert *aci.NodeAlert) (result *aci.NodeAlert, err error) {
	result = &aci.NodeAlert{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource(aci.ResourceTypeNodeAlert).
		Name(alert.Name).
		SubResource("status").
		Body(alert).
		Do().
		Into(result)
	return
}
