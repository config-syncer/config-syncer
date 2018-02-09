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
	v1alpha1 "github.com/appscode/kubed/pkg/apis/kubed/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	testing "k8s.io/client-go/testing"
)

// FakeStuffs implements StuffInterface
type FakeStuffs struct {
	Fake *FakeKubed
	ns   string
}

var stuffsResource = schema.GroupVersionResource{Group: "kubed.appscode.com", Version: "", Resource: "stuffs"}

var stuffsKind = schema.GroupVersionKind{Group: "kubed.appscode.com", Version: "", Kind: "Stuff"}

// Get takes name of the searchResult, and returns the corresponding searchResult object, and an error if there is any.
func (c *FakeStuffs) Get(name string, options v1.GetOptions) (result *v1alpha1.Stuff, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(stuffsResource, c.ns, name), &v1alpha1.Stuff{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Stuff), err
}
