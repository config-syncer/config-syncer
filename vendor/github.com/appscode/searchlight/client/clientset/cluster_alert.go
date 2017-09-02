package clientset

import (
	aci "github.com/appscode/searchlight/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

type ClusterAlertGetter interface {
	ClusterAlerts(namespace string) ClusterAlertInterface
}

type ClusterAlertInterface interface {
	List(opts metav1.ListOptions) (*aci.ClusterAlertList, error)
	Get(name string) (*aci.ClusterAlert, error)
	Create(ClusterAlert *aci.ClusterAlert) (*aci.ClusterAlert, error)
	Update(ClusterAlert *aci.ClusterAlert) (*aci.ClusterAlert, error)
	Delete(name string) error
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	UpdateStatus(ClusterAlert *aci.ClusterAlert) (*aci.ClusterAlert, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (*aci.ClusterAlert, error)
}

type ClusterAlertImpl struct {
	r  rest.Interface
	ns string
}

var _ ClusterAlertInterface = &ClusterAlertImpl{}

func newClusterAlert(c *ExtensionClient, namespace string) *ClusterAlertImpl {
	return &ClusterAlertImpl{c.restClient, namespace}
}

func (c *ClusterAlertImpl) List(opts metav1.ListOptions) (result *aci.ClusterAlertList, err error) {
	result = &aci.ClusterAlertList{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource(aci.ResourceTypeClusterAlert).
		VersionedParams(&opts, ExtendedCodec).
		Do().
		Into(result)
	return
}

func (c *ClusterAlertImpl) Get(name string) (result *aci.ClusterAlert, err error) {
	result = &aci.ClusterAlert{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource(aci.ResourceTypeClusterAlert).
		Name(name).
		Do().
		Into(result)
	return
}

func (c *ClusterAlertImpl) Create(alert *aci.ClusterAlert) (result *aci.ClusterAlert, err error) {
	result = &aci.ClusterAlert{}
	err = c.r.Post().
		Namespace(c.ns).
		Resource(aci.ResourceTypeClusterAlert).
		Body(alert).
		Do().
		Into(result)
	return
}

func (c *ClusterAlertImpl) Update(alert *aci.ClusterAlert) (result *aci.ClusterAlert, err error) {
	result = &aci.ClusterAlert{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource(aci.ResourceTypeClusterAlert).
		Name(alert.Name).
		Body(alert).
		Do().
		Into(result)
	return
}

func (c *ClusterAlertImpl) Delete(name string) (err error) {
	return c.r.Delete().
		Namespace(c.ns).
		Resource(aci.ResourceTypeClusterAlert).
		Name(name).
		Do().
		Error()
}

func (c *ClusterAlertImpl) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return c.r.Get().
		Prefix("watch").
		Namespace(c.ns).
		Resource(aci.ResourceTypeClusterAlert).
		VersionedParams(&opts, ExtendedCodec).
		Watch()
}

func (c *ClusterAlertImpl) UpdateStatus(alert *aci.ClusterAlert) (result *aci.ClusterAlert, err error) {
	result = &aci.ClusterAlert{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource(aci.ResourceTypeClusterAlert).
		Name(alert.Name).
		SubResource("status").
		Body(alert).
		Do().
		Into(result)
	return
}

func (c *ClusterAlertImpl) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *aci.ClusterAlert, err error) {
	result = &aci.ClusterAlert{}
	err = c.r.Patch(pt).
		Namespace(c.ns).
		Resource(aci.ResourceTypeClusterAlert).
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
