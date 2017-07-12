package clientset

import (
	tapi "github.com/appscode/searchlight/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

type PodAlertGetter interface {
	PodAlerts(namespace string) PodAlertInterface
}

type PodAlertInterface interface {
	List(opts metav1.ListOptions) (*tapi.PodAlertList, error)
	Get(name string) (*tapi.PodAlert, error)
	Create(Alert *tapi.PodAlert) (*tapi.PodAlert, error)
	Update(Alert *tapi.PodAlert) (*tapi.PodAlert, error)
	Delete(name string) error
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	UpdateStatus(Alert *tapi.PodAlert) (*tapi.PodAlert, error)
}

type PodAlertImpl struct {
	r  rest.Interface
	ns string
}

var _ PodAlertInterface = &PodAlertImpl{}

func newPodAlert(c *ExtensionClient, namespace string) *PodAlertImpl {
	return &PodAlertImpl{c.restClient, namespace}
}

func (c *PodAlertImpl) List(opts metav1.ListOptions) (result *tapi.PodAlertList, err error) {
	result = &tapi.PodAlertList{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource(tapi.ResourceTypePodAlert).
		VersionedParams(&opts, ExtendedCodec).
		Do().
		Into(result)
	return
}

func (c *PodAlertImpl) Get(name string) (result *tapi.PodAlert, err error) {
	result = &tapi.PodAlert{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource(tapi.ResourceTypePodAlert).
		Name(name).
		Do().
		Into(result)
	return
}

func (c *PodAlertImpl) Create(alert *tapi.PodAlert) (result *tapi.PodAlert, err error) {
	result = &tapi.PodAlert{}
	err = c.r.Post().
		Namespace(c.ns).
		Resource(tapi.ResourceTypePodAlert).
		Body(alert).
		Do().
		Into(result)
	return
}

func (c *PodAlertImpl) Update(alert *tapi.PodAlert) (result *tapi.PodAlert, err error) {
	result = &tapi.PodAlert{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource(tapi.ResourceTypePodAlert).
		Name(alert.Name).
		Body(alert).
		Do().
		Into(result)
	return
}

func (c *PodAlertImpl) Delete(name string) (err error) {
	return c.r.Delete().
		Namespace(c.ns).
		Resource(tapi.ResourceTypePodAlert).
		Name(name).
		Do().
		Error()
}

func (c *PodAlertImpl) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return c.r.Get().
		Prefix("watch").
		Namespace(c.ns).
		Resource(tapi.ResourceTypePodAlert).
		VersionedParams(&opts, ExtendedCodec).
		Watch()
}

func (c *PodAlertImpl) UpdateStatus(alert *tapi.PodAlert) (result *tapi.PodAlert, err error) {
	result = &tapi.PodAlert{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource(tapi.ResourceTypePodAlert).
		Name(alert.Name).
		SubResource("status").
		Body(alert).
		Do().
		Into(result)
	return
}
