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
package fake

import (
	kubed "github.com/appscode/kubed/pkg/apis/kubed"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeSearchResults implements SearchResultInterface
type FakeSearchResults struct {
	Fake *FakeKubed
	ns   string
}

var searchresultsResource = schema.GroupVersionResource{Group: "kubed.appscode.com", Version: "", Resource: "searchresults"}

var searchresultsKind = schema.GroupVersionKind{Group: "kubed.appscode.com", Version: "", Kind: "SearchResult"}

// Get takes name of the searchResult, and returns the corresponding searchResult object, and an error if there is any.
func (c *FakeSearchResults) Get(name string, options v1.GetOptions) (result *kubed.SearchResult, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(searchresultsResource, c.ns, name), &kubed.SearchResult{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kubed.SearchResult), err
}

// List takes label and field selectors, and returns the list of SearchResults that match those selectors.
func (c *FakeSearchResults) List(opts v1.ListOptions) (result *kubed.SearchResultList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(searchresultsResource, searchresultsKind, c.ns, opts), &kubed.SearchResultList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &kubed.SearchResultList{}
	for _, item := range obj.(*kubed.SearchResultList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested searchResults.
func (c *FakeSearchResults) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(searchresultsResource, c.ns, opts))

}

// Create takes the representation of a searchResult and creates it.  Returns the server's representation of the searchResult, and an error, if there is any.
func (c *FakeSearchResults) Create(searchResult *kubed.SearchResult) (result *kubed.SearchResult, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(searchresultsResource, c.ns, searchResult), &kubed.SearchResult{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kubed.SearchResult), err
}

// Update takes the representation of a searchResult and updates it. Returns the server's representation of the searchResult, and an error, if there is any.
func (c *FakeSearchResults) Update(searchResult *kubed.SearchResult) (result *kubed.SearchResult, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(searchresultsResource, c.ns, searchResult), &kubed.SearchResult{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kubed.SearchResult), err
}

// Delete takes name of the searchResult and deletes it. Returns an error if one occurs.
func (c *FakeSearchResults) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(searchresultsResource, c.ns, name), &kubed.SearchResult{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSearchResults) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(searchresultsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &kubed.SearchResultList{})
	return err
}

// Patch applies the patch and returns the patched searchResult.
func (c *FakeSearchResults) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *kubed.SearchResult, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(searchresultsResource, c.ns, name, data, subresources...), &kubed.SearchResult{})

	if obj == nil {
		return nil, err
	}
	return obj.(*kubed.SearchResult), err
}
