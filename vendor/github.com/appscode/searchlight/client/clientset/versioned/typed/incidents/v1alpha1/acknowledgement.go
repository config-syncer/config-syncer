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
	v1alpha1 "github.com/appscode/searchlight/apis/incidents/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rest "k8s.io/client-go/rest"
)

// AcknowledgementsGetter has a method to return a AcknowledgementInterface.
// A group's client should implement this interface.
type AcknowledgementsGetter interface {
	Acknowledgements(namespace string) AcknowledgementInterface
}

// AcknowledgementInterface has methods to work with Acknowledgement resources.
type AcknowledgementInterface interface {
	Create(*v1alpha1.Acknowledgement) (*v1alpha1.Acknowledgement, error)
	Delete(name string, options *v1.DeleteOptions) error
	AcknowledgementExpansion
}

// acknowledgements implements AcknowledgementInterface
type acknowledgements struct {
	client rest.Interface
	ns     string
}

// newAcknowledgements returns a Acknowledgements
func newAcknowledgements(c *IncidentsV1alpha1Client, namespace string) *acknowledgements {
	return &acknowledgements{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Create takes the representation of a acknowledgement and creates it.  Returns the server's representation of the acknowledgement, and an error, if there is any.
func (c *acknowledgements) Create(acknowledgement *v1alpha1.Acknowledgement) (result *v1alpha1.Acknowledgement, err error) {
	result = &v1alpha1.Acknowledgement{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("acknowledgements").
		Body(acknowledgement).
		Do().
		Into(result)
	return
}

// Delete takes name of the acknowledgement and deletes it. Returns an error if one occurs.
func (c *acknowledgements) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("acknowledgements").
		Name(name).
		Body(options).
		Do().
		Error()
}
