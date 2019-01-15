/*
Copyright 2019 The Stash Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha2

import (
	v1alpha2 "github.com/appscode/stash/apis/stash/v1alpha2"
	scheme "github.com/appscode/stash/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// AgentTemplatesGetter has a method to return a AgentTemplateInterface.
// A group's client should implement this interface.
type AgentTemplatesGetter interface {
	AgentTemplates() AgentTemplateInterface
}

// AgentTemplateInterface has methods to work with AgentTemplate resources.
type AgentTemplateInterface interface {
	Create(*v1alpha2.AgentTemplate) (*v1alpha2.AgentTemplate, error)
	Update(*v1alpha2.AgentTemplate) (*v1alpha2.AgentTemplate, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha2.AgentTemplate, error)
	List(opts v1.ListOptions) (*v1alpha2.AgentTemplateList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha2.AgentTemplate, err error)
	AgentTemplateExpansion
}

// agentTemplates implements AgentTemplateInterface
type agentTemplates struct {
	client rest.Interface
}

// newAgentTemplates returns a AgentTemplates
func newAgentTemplates(c *StashV1alpha2Client) *agentTemplates {
	return &agentTemplates{
		client: c.RESTClient(),
	}
}

// Get takes name of the agentTemplate, and returns the corresponding agentTemplate object, and an error if there is any.
func (c *agentTemplates) Get(name string, options v1.GetOptions) (result *v1alpha2.AgentTemplate, err error) {
	result = &v1alpha2.AgentTemplate{}
	err = c.client.Get().
		Resource("agenttemplates").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of AgentTemplates that match those selectors.
func (c *agentTemplates) List(opts v1.ListOptions) (result *v1alpha2.AgentTemplateList, err error) {
	result = &v1alpha2.AgentTemplateList{}
	err = c.client.Get().
		Resource("agenttemplates").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested agentTemplates.
func (c *agentTemplates) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Resource("agenttemplates").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a agentTemplate and creates it.  Returns the server's representation of the agentTemplate, and an error, if there is any.
func (c *agentTemplates) Create(agentTemplate *v1alpha2.AgentTemplate) (result *v1alpha2.AgentTemplate, err error) {
	result = &v1alpha2.AgentTemplate{}
	err = c.client.Post().
		Resource("agenttemplates").
		Body(agentTemplate).
		Do().
		Into(result)
	return
}

// Update takes the representation of a agentTemplate and updates it. Returns the server's representation of the agentTemplate, and an error, if there is any.
func (c *agentTemplates) Update(agentTemplate *v1alpha2.AgentTemplate) (result *v1alpha2.AgentTemplate, err error) {
	result = &v1alpha2.AgentTemplate{}
	err = c.client.Put().
		Resource("agenttemplates").
		Name(agentTemplate.Name).
		Body(agentTemplate).
		Do().
		Into(result)
	return
}

// Delete takes name of the agentTemplate and deletes it. Returns an error if one occurs.
func (c *agentTemplates) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("agenttemplates").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *agentTemplates) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Resource("agenttemplates").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched agentTemplate.
func (c *agentTemplates) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha2.AgentTemplate, err error) {
	result = &v1alpha2.AgentTemplate{}
	err = c.client.Patch(pt).
		Resource("agenttemplates").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
