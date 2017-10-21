/*
Copyright 2017 The Searchlight Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	v1alpha1 "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	scheme "github.com/appscode/searchlight/client/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// NodeAlertsGetter has a method to return a NodeAlertInterface.
// A group's client should implement this interface.
type NodeAlertsGetter interface {
	NodeAlerts(namespace string) NodeAlertInterface
}

// NodeAlertInterface has methods to work with NodeAlert resources.
type NodeAlertInterface interface {
	Create(*v1alpha1.NodeAlert) (*v1alpha1.NodeAlert, error)
	Update(*v1alpha1.NodeAlert) (*v1alpha1.NodeAlert, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.NodeAlert, error)
	List(opts v1.ListOptions) (*v1alpha1.NodeAlertList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.NodeAlert, err error)
	NodeAlertExpansion
}

// nodeAlerts implements NodeAlertInterface
type nodeAlerts struct {
	client rest.Interface
	ns     string
}

// newNodeAlerts returns a NodeAlerts
func newNodeAlerts(c *MonitoringV1alpha1Client, namespace string) *nodeAlerts {
	return &nodeAlerts{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the nodeAlert, and returns the corresponding nodeAlert object, and an error if there is any.
func (c *nodeAlerts) Get(name string, options v1.GetOptions) (result *v1alpha1.NodeAlert, err error) {
	result = &v1alpha1.NodeAlert{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("nodealerts").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of NodeAlerts that match those selectors.
func (c *nodeAlerts) List(opts v1.ListOptions) (result *v1alpha1.NodeAlertList, err error) {
	result = &v1alpha1.NodeAlertList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("nodealerts").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested nodeAlerts.
func (c *nodeAlerts) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("nodealerts").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a nodeAlert and creates it.  Returns the server's representation of the nodeAlert, and an error, if there is any.
func (c *nodeAlerts) Create(nodeAlert *v1alpha1.NodeAlert) (result *v1alpha1.NodeAlert, err error) {
	result = &v1alpha1.NodeAlert{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("nodealerts").
		Body(nodeAlert).
		Do().
		Into(result)
	return
}

// Update takes the representation of a nodeAlert and updates it. Returns the server's representation of the nodeAlert, and an error, if there is any.
func (c *nodeAlerts) Update(nodeAlert *v1alpha1.NodeAlert) (result *v1alpha1.NodeAlert, err error) {
	result = &v1alpha1.NodeAlert{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("nodealerts").
		Name(nodeAlert.Name).
		Body(nodeAlert).
		Do().
		Into(result)
	return
}

// Delete takes name of the nodeAlert and deletes it. Returns an error if one occurs.
func (c *nodeAlerts) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("nodealerts").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *nodeAlerts) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("nodealerts").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched nodeAlert.
func (c *nodeAlerts) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.NodeAlert, err error) {
	result = &v1alpha1.NodeAlert{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("nodealerts").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
