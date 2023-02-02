/*
Copyright 2023 The Crossplane Authors.

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

package fake

import (
	"context"

	v1beta1 "github.com/crossplane/crossplane/apis/pkg/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeConfigurations implements ConfigurationInterface
type FakeConfigurations struct {
	Fake *FakePkgV1beta1
}

var configurationsResource = schema.GroupVersionResource{Group: "pkg.ibm.crossplane.io", Version: "v1beta1", Resource: "configurations"}

var configurationsKind = schema.GroupVersionKind{Group: "pkg.ibm.crossplane.io", Version: "v1beta1", Kind: "Configuration"}

// Get takes name of the configuration, and returns the corresponding configuration object, and an error if there is any.
func (c *FakeConfigurations) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.Configuration, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(configurationsResource, name), &v1beta1.Configuration{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Configuration), err
}

// List takes label and field selectors, and returns the list of Configurations that match those selectors.
func (c *FakeConfigurations) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.ConfigurationList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(configurationsResource, configurationsKind, opts), &v1beta1.ConfigurationList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.ConfigurationList{ListMeta: obj.(*v1beta1.ConfigurationList).ListMeta}
	for _, item := range obj.(*v1beta1.ConfigurationList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested configurations.
func (c *FakeConfigurations) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(configurationsResource, opts))
}

// Create takes the representation of a configuration and creates it.  Returns the server's representation of the configuration, and an error, if there is any.
func (c *FakeConfigurations) Create(ctx context.Context, configuration *v1beta1.Configuration, opts v1.CreateOptions) (result *v1beta1.Configuration, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(configurationsResource, configuration), &v1beta1.Configuration{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Configuration), err
}

// Update takes the representation of a configuration and updates it. Returns the server's representation of the configuration, and an error, if there is any.
func (c *FakeConfigurations) Update(ctx context.Context, configuration *v1beta1.Configuration, opts v1.UpdateOptions) (result *v1beta1.Configuration, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(configurationsResource, configuration), &v1beta1.Configuration{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Configuration), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeConfigurations) UpdateStatus(ctx context.Context, configuration *v1beta1.Configuration, opts v1.UpdateOptions) (*v1beta1.Configuration, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(configurationsResource, "status", configuration), &v1beta1.Configuration{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Configuration), err
}

// Delete takes name of the configuration and deletes it. Returns an error if one occurs.
func (c *FakeConfigurations) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(configurationsResource, name), &v1beta1.Configuration{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeConfigurations) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(configurationsResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta1.ConfigurationList{})
	return err
}

// Patch applies the patch and returns the patched configuration.
func (c *FakeConfigurations) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.Configuration, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(configurationsResource, name, pt, data, subresources...), &v1beta1.Configuration{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Configuration), err
}
