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
	v1alpha1 "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	scheme "github.com/kubedb/apimachinery/client/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// RedisesGetter has a method to return a RedisInterface.
// A group's client should implement this interface.
type RedisesGetter interface {
	Redises(namespace string) RedisInterface
}

// RedisInterface has methods to work with Redis resources.
type RedisInterface interface {
	Create(*v1alpha1.Redis) (*v1alpha1.Redis, error)
	Update(*v1alpha1.Redis) (*v1alpha1.Redis, error)
	UpdateStatus(*v1alpha1.Redis) (*v1alpha1.Redis, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Redis, error)
	List(opts v1.ListOptions) (*v1alpha1.RedisList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Redis, err error)
	RedisExpansion
}

// redises implements RedisInterface
type redises struct {
	client rest.Interface
	ns     string
}

// newRedises returns a Redises
func newRedises(c *KubedbV1alpha1Client, namespace string) *redises {
	return &redises{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the redis, and returns the corresponding redis object, and an error if there is any.
func (c *redises) Get(name string, options v1.GetOptions) (result *v1alpha1.Redis, err error) {
	result = &v1alpha1.Redis{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("redises").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Redises that match those selectors.
func (c *redises) List(opts v1.ListOptions) (result *v1alpha1.RedisList, err error) {
	result = &v1alpha1.RedisList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("redises").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested redises.
func (c *redises) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("redises").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a redis and creates it.  Returns the server's representation of the redis, and an error, if there is any.
func (c *redises) Create(redis *v1alpha1.Redis) (result *v1alpha1.Redis, err error) {
	result = &v1alpha1.Redis{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("redises").
		Body(redis).
		Do().
		Into(result)
	return
}

// Update takes the representation of a redis and updates it. Returns the server's representation of the redis, and an error, if there is any.
func (c *redises) Update(redis *v1alpha1.Redis) (result *v1alpha1.Redis, err error) {
	result = &v1alpha1.Redis{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("redises").
		Name(redis.Name).
		Body(redis).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *redises) UpdateStatus(redis *v1alpha1.Redis) (result *v1alpha1.Redis, err error) {
	result = &v1alpha1.Redis{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("redises").
		Name(redis.Name).
		SubResource("status").
		Body(redis).
		Do().
		Into(result)
	return
}

// Delete takes name of the redis and deletes it. Returns an error if one occurs.
func (c *redises) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("redises").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *redises) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("redises").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched redis.
func (c *redises) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Redis, err error) {
	result = &v1alpha1.Redis{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("redises").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
