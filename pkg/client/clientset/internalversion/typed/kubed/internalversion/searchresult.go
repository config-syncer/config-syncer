/*
Copyright 2018 The Kubed Authors.

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
package internalversion

import (
	kubed "github.com/appscode/kubed/pkg/apis/kubed"
	scheme "github.com/appscode/kubed/pkg/client/clientset/internalversion/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// SearchResultsGetter has a method to return a SearchResultInterface.
// A group's client should implement this interface.
type SearchResultsGetter interface {
	SearchResults(namespace string) SearchResultInterface
}

// SearchResultInterface has methods to work with SearchResult resources.
type SearchResultInterface interface {
	Create(*kubed.SearchResult) (*kubed.SearchResult, error)
	Update(*kubed.SearchResult) (*kubed.SearchResult, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*kubed.SearchResult, error)
	List(opts v1.ListOptions) (*kubed.SearchResultList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *kubed.SearchResult, err error)
	SearchResultExpansion
}

// searchResults implements SearchResultInterface
type searchResults struct {
	client rest.Interface
	ns     string
}

// newSearchResults returns a SearchResults
func newSearchResults(c *KubedClient, namespace string) *searchResults {
	return &searchResults{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the searchResult, and returns the corresponding searchResult object, and an error if there is any.
func (c *searchResults) Get(name string, options v1.GetOptions) (result *kubed.SearchResult, err error) {
	result = &kubed.SearchResult{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("searchresults").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of SearchResults that match those selectors.
func (c *searchResults) List(opts v1.ListOptions) (result *kubed.SearchResultList, err error) {
	result = &kubed.SearchResultList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("searchresults").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested searchResults.
func (c *searchResults) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("searchresults").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a searchResult and creates it.  Returns the server's representation of the searchResult, and an error, if there is any.
func (c *searchResults) Create(searchResult *kubed.SearchResult) (result *kubed.SearchResult, err error) {
	result = &kubed.SearchResult{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("searchresults").
		Body(searchResult).
		Do().
		Into(result)
	return
}

// Update takes the representation of a searchResult and updates it. Returns the server's representation of the searchResult, and an error, if there is any.
func (c *searchResults) Update(searchResult *kubed.SearchResult) (result *kubed.SearchResult, err error) {
	result = &kubed.SearchResult{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("searchresults").
		Name(searchResult.Name).
		Body(searchResult).
		Do().
		Into(result)
	return
}

// Delete takes name of the searchResult and deletes it. Returns an error if one occurs.
func (c *searchResults) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("searchresults").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *searchResults) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("searchresults").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched searchResult.
func (c *searchResults) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *kubed.SearchResult, err error) {
	result = &kubed.SearchResult{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("searchresults").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
