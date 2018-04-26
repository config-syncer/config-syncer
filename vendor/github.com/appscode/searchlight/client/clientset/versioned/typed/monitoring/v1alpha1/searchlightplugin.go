/*
Copyright 2018 The Searchlight Authors.

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
	scheme "github.com/appscode/searchlight/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// SearchlightPluginsGetter has a method to return a SearchlightPluginInterface.
// A group's client should implement this interface.
type SearchlightPluginsGetter interface {
	SearchlightPlugins() SearchlightPluginInterface
}

// SearchlightPluginInterface has methods to work with SearchlightPlugin resources.
type SearchlightPluginInterface interface {
	Create(*v1alpha1.SearchlightPlugin) (*v1alpha1.SearchlightPlugin, error)
	Update(*v1alpha1.SearchlightPlugin) (*v1alpha1.SearchlightPlugin, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.SearchlightPlugin, error)
	List(opts v1.ListOptions) (*v1alpha1.SearchlightPluginList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.SearchlightPlugin, err error)
	SearchlightPluginExpansion
}

// searchlightPlugins implements SearchlightPluginInterface
type searchlightPlugins struct {
	client rest.Interface
}

// newSearchlightPlugins returns a SearchlightPlugins
func newSearchlightPlugins(c *MonitoringV1alpha1Client) *searchlightPlugins {
	return &searchlightPlugins{
		client: c.RESTClient(),
	}
}

// Get takes name of the searchlightPlugin, and returns the corresponding searchlightPlugin object, and an error if there is any.
func (c *searchlightPlugins) Get(name string, options v1.GetOptions) (result *v1alpha1.SearchlightPlugin, err error) {
	result = &v1alpha1.SearchlightPlugin{}
	err = c.client.Get().
		Resource("searchlightplugins").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of SearchlightPlugins that match those selectors.
func (c *searchlightPlugins) List(opts v1.ListOptions) (result *v1alpha1.SearchlightPluginList, err error) {
	result = &v1alpha1.SearchlightPluginList{}
	err = c.client.Get().
		Resource("searchlightplugins").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested searchlightPlugins.
func (c *searchlightPlugins) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Resource("searchlightplugins").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a searchlightPlugin and creates it.  Returns the server's representation of the searchlightPlugin, and an error, if there is any.
func (c *searchlightPlugins) Create(searchlightPlugin *v1alpha1.SearchlightPlugin) (result *v1alpha1.SearchlightPlugin, err error) {
	result = &v1alpha1.SearchlightPlugin{}
	err = c.client.Post().
		Resource("searchlightplugins").
		Body(searchlightPlugin).
		Do().
		Into(result)
	return
}

// Update takes the representation of a searchlightPlugin and updates it. Returns the server's representation of the searchlightPlugin, and an error, if there is any.
func (c *searchlightPlugins) Update(searchlightPlugin *v1alpha1.SearchlightPlugin) (result *v1alpha1.SearchlightPlugin, err error) {
	result = &v1alpha1.SearchlightPlugin{}
	err = c.client.Put().
		Resource("searchlightplugins").
		Name(searchlightPlugin.Name).
		Body(searchlightPlugin).
		Do().
		Into(result)
	return
}

// Delete takes name of the searchlightPlugin and deletes it. Returns an error if one occurs.
func (c *searchlightPlugins) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("searchlightplugins").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *searchlightPlugins) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Resource("searchlightplugins").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched searchlightPlugin.
func (c *searchlightPlugins) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.SearchlightPlugin, err error) {
	result = &v1alpha1.SearchlightPlugin{}
	err = c.client.Patch(pt).
		Resource("searchlightplugins").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
