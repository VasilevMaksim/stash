/*
Copyright 2017 The Stash Authors.

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
	v1alpha1 "github.com/appscode/stash/apis/stash/v1alpha1"
	scheme "github.com/appscode/stash/client/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ResticsGetter has a method to return a ResticInterface.
// A group's client should implement this interface.
type ResticsGetter interface {
	Restics(namespace string) ResticInterface
}

// ResticInterface has methods to work with Restic resources.
type ResticInterface interface {
	Create(*v1alpha1.Restic) (*v1alpha1.Restic, error)
	Update(*v1alpha1.Restic) (*v1alpha1.Restic, error)
	UpdateStatus(*v1alpha1.Restic) (*v1alpha1.Restic, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Restic, error)
	List(opts v1.ListOptions) (*v1alpha1.ResticList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Restic, err error)
	ResticExpansion
}

// restics implements ResticInterface
type restics struct {
	client rest.Interface
	ns     string
}

// newRestics returns a Restics
func newRestics(c *StashV1alpha1Client, namespace string) *restics {
	return &restics{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the restic, and returns the corresponding restic object, and an error if there is any.
func (c *restics) Get(name string, options v1.GetOptions) (result *v1alpha1.Restic, err error) {
	result = &v1alpha1.Restic{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("restics").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Restics that match those selectors.
func (c *restics) List(opts v1.ListOptions) (result *v1alpha1.ResticList, err error) {
	result = &v1alpha1.ResticList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("restics").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested restics.
func (c *restics) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("restics").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a restic and creates it.  Returns the server's representation of the restic, and an error, if there is any.
func (c *restics) Create(restic *v1alpha1.Restic) (result *v1alpha1.Restic, err error) {
	result = &v1alpha1.Restic{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("restics").
		Body(restic).
		Do().
		Into(result)
	return
}

// Update takes the representation of a restic and updates it. Returns the server's representation of the restic, and an error, if there is any.
func (c *restics) Update(restic *v1alpha1.Restic) (result *v1alpha1.Restic, err error) {
	result = &v1alpha1.Restic{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("restics").
		Name(restic.Name).
		Body(restic).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *restics) UpdateStatus(restic *v1alpha1.Restic) (result *v1alpha1.Restic, err error) {
	result = &v1alpha1.Restic{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("restics").
		Name(restic.Name).
		SubResource("status").
		Body(restic).
		Do().
		Into(result)
	return
}

// Delete takes name of the restic and deletes it. Returns an error if one occurs.
func (c *restics) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("restics").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *restics) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("restics").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched restic.
func (c *restics) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Restic, err error) {
	result = &v1alpha1.Restic{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("restics").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
