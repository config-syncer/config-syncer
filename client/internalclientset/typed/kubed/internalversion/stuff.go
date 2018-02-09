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
	kubed "github.com/appscode/kubed/apis/kubed"
	scheme "github.com/appscode/kubed/client/internalclientset/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rest "k8s.io/client-go/rest"
)

// StuffsGetter has a method to return a StuffInterface.
// A group's client should implement this interface.
type StuffsGetter interface {
	Stuffs(namespace string) StuffInterface
}

// StuffInterface has methods to work with Stuff resources.
type StuffInterface interface {
	Get(name string, options v1.GetOptions) (*kubed.Stuff, error)
	StuffExpansion
}

// stuffs implements StuffInterface
type stuffs struct {
	client rest.Interface
	ns     string
}

// newStuffs returns a Stuffs
func newStuffs(c *KubedClient, namespace string) *stuffs {
	return &stuffs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the stuff, and returns the corresponding stuff object, and an error if there is any.
func (c *stuffs) Get(name string, options v1.GetOptions) (result *kubed.Stuff, err error) {
	result = &kubed.Stuff{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("stuffs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}
