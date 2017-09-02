package clientset

import (
	aci "github.com/appscode/searchlight/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

type PodAlertGetter interface {
	PodAlerts(namespace string) PodAlertInterface
}

type PodAlertInterface interface {
	List(opts metav1.ListOptions) (*aci.PodAlertList, error)
	Get(name string) (*aci.PodAlert, error)
	Create(Alert *aci.PodAlert) (*aci.PodAlert, error)
	Update(Alert *aci.PodAlert) (*aci.PodAlert, error)
	Delete(name string) error
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	UpdateStatus(Alert *aci.PodAlert) (*aci.PodAlert, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (*aci.PodAlert, error)
}

type PodAlertImpl struct {
	r  rest.Interface
	ns string
}

var _ PodAlertInterface = &PodAlertImpl{}

func newPodAlert(c *ExtensionClient, namespace string) *PodAlertImpl {
	return &PodAlertImpl{c.restClient, namespace}
}

func (c *PodAlertImpl) List(opts metav1.ListOptions) (result *aci.PodAlertList, err error) {
	result = &aci.PodAlertList{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource(aci.ResourceTypePodAlert).
		VersionedParams(&opts, ExtendedCodec).
		Do().
		Into(result)
	return
}

func (c *PodAlertImpl) Get(name string) (result *aci.PodAlert, err error) {
	result = &aci.PodAlert{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource(aci.ResourceTypePodAlert).
		Name(name).
		Do().
		Into(result)
	return
}

func (c *PodAlertImpl) Create(alert *aci.PodAlert) (result *aci.PodAlert, err error) {
	result = &aci.PodAlert{}
	err = c.r.Post().
		Namespace(c.ns).
		Resource(aci.ResourceTypePodAlert).
		Body(alert).
		Do().
		Into(result)
	return
}

func (c *PodAlertImpl) Update(alert *aci.PodAlert) (result *aci.PodAlert, err error) {
	result = &aci.PodAlert{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource(aci.ResourceTypePodAlert).
		Name(alert.Name).
		Body(alert).
		Do().
		Into(result)
	return
}

func (c *PodAlertImpl) Delete(name string) (err error) {
	return c.r.Delete().
		Namespace(c.ns).
		Resource(aci.ResourceTypePodAlert).
		Name(name).
		Do().
		Error()
}

func (c *PodAlertImpl) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return c.r.Get().
		Prefix("watch").
		Namespace(c.ns).
		Resource(aci.ResourceTypePodAlert).
		VersionedParams(&opts, ExtendedCodec).
		Watch()
}

func (c *PodAlertImpl) UpdateStatus(alert *aci.PodAlert) (result *aci.PodAlert, err error) {
	result = &aci.PodAlert{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource(aci.ResourceTypePodAlert).
		Name(alert.Name).
		SubResource("status").
		Body(alert).
		Do().
		Into(result)
	return
}

func (c *PodAlertImpl) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *aci.PodAlert, err error) {
	result = &aci.PodAlert{}
	err = c.r.Patch(pt).
		Namespace(c.ns).
		Resource(aci.ResourceTypePodAlert).
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
