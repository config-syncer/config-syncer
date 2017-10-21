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

// ClusterAlertsGetter has a method to return a ClusterAlertInterface.
// A group's client should implement this interface.
type ClusterAlertsGetter interface {
	ClusterAlerts(namespace string) ClusterAlertInterface
}

// ClusterAlertInterface has methods to work with ClusterAlert resources.
type ClusterAlertInterface interface {
	Create(*v1alpha1.ClusterAlert) (*v1alpha1.ClusterAlert, error)
	Update(*v1alpha1.ClusterAlert) (*v1alpha1.ClusterAlert, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.ClusterAlert, error)
	List(opts v1.ListOptions) (*v1alpha1.ClusterAlertList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.ClusterAlert, err error)
	ClusterAlertExpansion
}

// clusterAlerts implements ClusterAlertInterface
type clusterAlerts struct {
	client rest.Interface
	ns     string
}

// newClusterAlerts returns a ClusterAlerts
func newClusterAlerts(c *MonitoringV1alpha1Client, namespace string) *clusterAlerts {
	return &clusterAlerts{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the clusterAlert, and returns the corresponding clusterAlert object, and an error if there is any.
func (c *clusterAlerts) Get(name string, options v1.GetOptions) (result *v1alpha1.ClusterAlert, err error) {
	result = &v1alpha1.ClusterAlert{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("clusteralerts").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ClusterAlerts that match those selectors.
func (c *clusterAlerts) List(opts v1.ListOptions) (result *v1alpha1.ClusterAlertList, err error) {
	result = &v1alpha1.ClusterAlertList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("clusteralerts").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested clusterAlerts.
func (c *clusterAlerts) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("clusteralerts").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a clusterAlert and creates it.  Returns the server's representation of the clusterAlert, and an error, if there is any.
func (c *clusterAlerts) Create(clusterAlert *v1alpha1.ClusterAlert) (result *v1alpha1.ClusterAlert, err error) {
	result = &v1alpha1.ClusterAlert{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("clusteralerts").
		Body(clusterAlert).
		Do().
		Into(result)
	return
}

// Update takes the representation of a clusterAlert and updates it. Returns the server's representation of the clusterAlert, and an error, if there is any.
func (c *clusterAlerts) Update(clusterAlert *v1alpha1.ClusterAlert) (result *v1alpha1.ClusterAlert, err error) {
	result = &v1alpha1.ClusterAlert{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("clusteralerts").
		Name(clusterAlert.Name).
		Body(clusterAlert).
		Do().
		Into(result)
	return
}

// Delete takes name of the clusterAlert and deletes it. Returns an error if one occurs.
func (c *clusterAlerts) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("clusteralerts").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *clusterAlerts) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("clusteralerts").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched clusterAlert.
func (c *clusterAlerts) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.ClusterAlert, err error) {
	result = &v1alpha1.ClusterAlert{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("clusteralerts").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
