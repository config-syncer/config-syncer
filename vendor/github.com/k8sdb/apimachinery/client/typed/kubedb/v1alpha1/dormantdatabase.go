/*
Copyright 2017 The KubeDB Authors.

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
	v1alpha1 "github.com/k8sdb/apimachinery/apis/kubedb/v1alpha1"
	scheme "github.com/k8sdb/apimachinery/client/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// DormantDatabasesGetter has a method to return a DormantDatabaseInterface.
// A group's client should implement this interface.
type DormantDatabasesGetter interface {
	DormantDatabases(namespace string) DormantDatabaseInterface
}

// DormantDatabaseInterface has methods to work with DormantDatabase resources.
type DormantDatabaseInterface interface {
	Create(*v1alpha1.DormantDatabase) (*v1alpha1.DormantDatabase, error)
	Update(*v1alpha1.DormantDatabase) (*v1alpha1.DormantDatabase, error)
	UpdateStatus(*v1alpha1.DormantDatabase) (*v1alpha1.DormantDatabase, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.DormantDatabase, error)
	List(opts v1.ListOptions) (*v1alpha1.DormantDatabaseList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.DormantDatabase, err error)
	DormantDatabaseExpansion
}

// dormantDatabases implements DormantDatabaseInterface
type dormantDatabases struct {
	client rest.Interface
	ns     string
}

// newDormantDatabases returns a DormantDatabases
func newDormantDatabases(c *KubedbV1alpha1Client, namespace string) *dormantDatabases {
	return &dormantDatabases{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the dormantDatabase, and returns the corresponding dormantDatabase object, and an error if there is any.
func (c *dormantDatabases) Get(name string, options v1.GetOptions) (result *v1alpha1.DormantDatabase, err error) {
	result = &v1alpha1.DormantDatabase{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("dormantdatabases").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of DormantDatabases that match those selectors.
func (c *dormantDatabases) List(opts v1.ListOptions) (result *v1alpha1.DormantDatabaseList, err error) {
	result = &v1alpha1.DormantDatabaseList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("dormantdatabases").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested dormantDatabases.
func (c *dormantDatabases) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("dormantdatabases").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a dormantDatabase and creates it.  Returns the server's representation of the dormantDatabase, and an error, if there is any.
func (c *dormantDatabases) Create(dormantDatabase *v1alpha1.DormantDatabase) (result *v1alpha1.DormantDatabase, err error) {
	result = &v1alpha1.DormantDatabase{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("dormantdatabases").
		Body(dormantDatabase).
		Do().
		Into(result)
	return
}

// Update takes the representation of a dormantDatabase and updates it. Returns the server's representation of the dormantDatabase, and an error, if there is any.
func (c *dormantDatabases) Update(dormantDatabase *v1alpha1.DormantDatabase) (result *v1alpha1.DormantDatabase, err error) {
	result = &v1alpha1.DormantDatabase{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("dormantdatabases").
		Name(dormantDatabase.Name).
		Body(dormantDatabase).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *dormantDatabases) UpdateStatus(dormantDatabase *v1alpha1.DormantDatabase) (result *v1alpha1.DormantDatabase, err error) {
	result = &v1alpha1.DormantDatabase{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("dormantdatabases").
		Name(dormantDatabase.Name).
		SubResource("status").
		Body(dormantDatabase).
		Do().
		Into(result)
	return
}

// Delete takes name of the dormantDatabase and deletes it. Returns an error if one occurs.
func (c *dormantDatabases) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("dormantdatabases").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *dormantDatabases) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("dormantdatabases").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched dormantDatabase.
func (c *dormantDatabases) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.DormantDatabase, err error) {
	result = &v1alpha1.DormantDatabase{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("dormantdatabases").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
