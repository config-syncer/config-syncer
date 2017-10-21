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

// ElasticsearchsGetter has a method to return a ElasticsearchInterface.
// A group's client should implement this interface.
type ElasticsearchsGetter interface {
	Elasticsearchs(namespace string) ElasticsearchInterface
}

// ElasticsearchInterface has methods to work with Elasticsearch resources.
type ElasticsearchInterface interface {
	Create(*v1alpha1.Elasticsearch) (*v1alpha1.Elasticsearch, error)
	Update(*v1alpha1.Elasticsearch) (*v1alpha1.Elasticsearch, error)
	UpdateStatus(*v1alpha1.Elasticsearch) (*v1alpha1.Elasticsearch, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Elasticsearch, error)
	List(opts v1.ListOptions) (*v1alpha1.ElasticsearchList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Elasticsearch, err error)
	ElasticsearchExpansion
}

// elasticsearchs implements ElasticsearchInterface
type elasticsearchs struct {
	client rest.Interface
	ns     string
}

// newElasticsearchs returns a Elasticsearchs
func newElasticsearchs(c *KubedbV1alpha1Client, namespace string) *elasticsearchs {
	return &elasticsearchs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the elasticsearch, and returns the corresponding elasticsearch object, and an error if there is any.
func (c *elasticsearchs) Get(name string, options v1.GetOptions) (result *v1alpha1.Elasticsearch, err error) {
	result = &v1alpha1.Elasticsearch{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("elasticsearchs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Elasticsearchs that match those selectors.
func (c *elasticsearchs) List(opts v1.ListOptions) (result *v1alpha1.ElasticsearchList, err error) {
	result = &v1alpha1.ElasticsearchList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("elasticsearchs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested elasticsearchs.
func (c *elasticsearchs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("elasticsearchs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a elasticsearch and creates it.  Returns the server's representation of the elasticsearch, and an error, if there is any.
func (c *elasticsearchs) Create(elasticsearch *v1alpha1.Elasticsearch) (result *v1alpha1.Elasticsearch, err error) {
	result = &v1alpha1.Elasticsearch{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("elasticsearchs").
		Body(elasticsearch).
		Do().
		Into(result)
	return
}

// Update takes the representation of a elasticsearch and updates it. Returns the server's representation of the elasticsearch, and an error, if there is any.
func (c *elasticsearchs) Update(elasticsearch *v1alpha1.Elasticsearch) (result *v1alpha1.Elasticsearch, err error) {
	result = &v1alpha1.Elasticsearch{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("elasticsearchs").
		Name(elasticsearch.Name).
		Body(elasticsearch).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *elasticsearchs) UpdateStatus(elasticsearch *v1alpha1.Elasticsearch) (result *v1alpha1.Elasticsearch, err error) {
	result = &v1alpha1.Elasticsearch{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("elasticsearchs").
		Name(elasticsearch.Name).
		SubResource("status").
		Body(elasticsearch).
		Do().
		Into(result)
	return
}

// Delete takes name of the elasticsearch and deletes it. Returns an error if one occurs.
func (c *elasticsearchs) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("elasticsearchs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *elasticsearchs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("elasticsearchs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched elasticsearch.
func (c *elasticsearchs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Elasticsearch, err error) {
	result = &v1alpha1.Elasticsearch{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("elasticsearchs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
