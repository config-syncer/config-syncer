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

// MongoDBsGetter has a method to return a MongoDBInterface.
// A group's client should implement this interface.
type MongoDBsGetter interface {
	MongoDBs(namespace string) MongoDBInterface
}

// MongoDBInterface has methods to work with MongoDB resources.
type MongoDBInterface interface {
	Create(*v1alpha1.MongoDB) (*v1alpha1.MongoDB, error)
	Update(*v1alpha1.MongoDB) (*v1alpha1.MongoDB, error)
	UpdateStatus(*v1alpha1.MongoDB) (*v1alpha1.MongoDB, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.MongoDB, error)
	List(opts v1.ListOptions) (*v1alpha1.MongoDBList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.MongoDB, err error)
	MongoDBExpansion
}

// mongoDBs implements MongoDBInterface
type mongoDBs struct {
	client rest.Interface
	ns     string
}

// newMongoDBs returns a MongoDBs
func newMongoDBs(c *KubedbV1alpha1Client, namespace string) *mongoDBs {
	return &mongoDBs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the mongoDB, and returns the corresponding mongoDB object, and an error if there is any.
func (c *mongoDBs) Get(name string, options v1.GetOptions) (result *v1alpha1.MongoDB, err error) {
	result = &v1alpha1.MongoDB{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("mongodbs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of MongoDBs that match those selectors.
func (c *mongoDBs) List(opts v1.ListOptions) (result *v1alpha1.MongoDBList, err error) {
	result = &v1alpha1.MongoDBList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("mongodbs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested mongoDBs.
func (c *mongoDBs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("mongodbs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a mongoDB and creates it.  Returns the server's representation of the mongoDB, and an error, if there is any.
func (c *mongoDBs) Create(mongoDB *v1alpha1.MongoDB) (result *v1alpha1.MongoDB, err error) {
	result = &v1alpha1.MongoDB{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("mongodbs").
		Body(mongoDB).
		Do().
		Into(result)
	return
}

// Update takes the representation of a mongoDB and updates it. Returns the server's representation of the mongoDB, and an error, if there is any.
func (c *mongoDBs) Update(mongoDB *v1alpha1.MongoDB) (result *v1alpha1.MongoDB, err error) {
	result = &v1alpha1.MongoDB{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("mongodbs").
		Name(mongoDB.Name).
		Body(mongoDB).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *mongoDBs) UpdateStatus(mongoDB *v1alpha1.MongoDB) (result *v1alpha1.MongoDB, err error) {
	result = &v1alpha1.MongoDB{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("mongodbs").
		Name(mongoDB.Name).
		SubResource("status").
		Body(mongoDB).
		Do().
		Into(result)
	return
}

// Delete takes name of the mongoDB and deletes it. Returns an error if one occurs.
func (c *mongoDBs) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("mongodbs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *mongoDBs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("mongodbs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched mongoDB.
func (c *mongoDBs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.MongoDB, err error) {
	result = &v1alpha1.MongoDB{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("mongodbs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
