/*
Copyright 2021 The Crossplane Authors.

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

	v1alpha1 "github.com/crossplane/crossplane/apis/pkg/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeControllerConfigs implements ControllerConfigInterface
type FakeControllerConfigs struct {
	Fake *FakePkgV1alpha1
}

var controllerconfigsResource = schema.GroupVersionResource{Group: "pkg.crossplane.io", Version: "v1alpha1", Resource: "controllerconfigs"}

var controllerconfigsKind = schema.GroupVersionKind{Group: "pkg.crossplane.io", Version: "v1alpha1", Kind: "ControllerConfig"}

// Get takes name of the controllerConfig, and returns the corresponding controllerConfig object, and an error if there is any.
func (c *FakeControllerConfigs) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ControllerConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(controllerconfigsResource, name), &v1alpha1.ControllerConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ControllerConfig), err
}

// List takes label and field selectors, and returns the list of ControllerConfigs that match those selectors.
func (c *FakeControllerConfigs) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ControllerConfigList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(controllerconfigsResource, controllerconfigsKind, opts), &v1alpha1.ControllerConfigList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ControllerConfigList{ListMeta: obj.(*v1alpha1.ControllerConfigList).ListMeta}
	for _, item := range obj.(*v1alpha1.ControllerConfigList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested controllerConfigs.
func (c *FakeControllerConfigs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(controllerconfigsResource, opts))
}

// Create takes the representation of a controllerConfig and creates it.  Returns the server's representation of the controllerConfig, and an error, if there is any.
func (c *FakeControllerConfigs) Create(ctx context.Context, controllerConfig *v1alpha1.ControllerConfig, opts v1.CreateOptions) (result *v1alpha1.ControllerConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(controllerconfigsResource, controllerConfig), &v1alpha1.ControllerConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ControllerConfig), err
}

// Update takes the representation of a controllerConfig and updates it. Returns the server's representation of the controllerConfig, and an error, if there is any.
func (c *FakeControllerConfigs) Update(ctx context.Context, controllerConfig *v1alpha1.ControllerConfig, opts v1.UpdateOptions) (result *v1alpha1.ControllerConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(controllerconfigsResource, controllerConfig), &v1alpha1.ControllerConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ControllerConfig), err
}

// Delete takes name of the controllerConfig and deletes it. Returns an error if one occurs.
func (c *FakeControllerConfigs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(controllerconfigsResource, name), &v1alpha1.ControllerConfig{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeControllerConfigs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(controllerconfigsResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ControllerConfigList{})
	return err
}

// Patch applies the patch and returns the patched controllerConfig.
func (c *FakeControllerConfigs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ControllerConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(controllerconfigsResource, name, pt, data, subresources...), &v1alpha1.ControllerConfig{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ControllerConfig), err
}
