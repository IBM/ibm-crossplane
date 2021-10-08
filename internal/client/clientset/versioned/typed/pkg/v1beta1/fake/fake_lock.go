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

	v1beta1 "github.com/crossplane/crossplane/apis/pkg/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeLocks implements LockInterface
type FakeLocks struct {
	Fake *FakePkgV1beta1
}

var locksResource = schema.GroupVersionResource{Group: "pkg.ibm.crossplane.io", Version: "v1beta1", Resource: "locks"}

var locksKind = schema.GroupVersionKind{Group: "pkg.ibm.crossplane.io", Version: "v1beta1", Kind: "Lock"}

// Get takes name of the lock, and returns the corresponding lock object, and an error if there is any.
func (c *FakeLocks) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.Lock, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(locksResource, name), &v1beta1.Lock{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Lock), err
}

// List takes label and field selectors, and returns the list of Locks that match those selectors.
func (c *FakeLocks) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.LockList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(locksResource, locksKind, opts), &v1beta1.LockList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.LockList{ListMeta: obj.(*v1beta1.LockList).ListMeta}
	for _, item := range obj.(*v1beta1.LockList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested locks.
func (c *FakeLocks) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(locksResource, opts))
}

// Create takes the representation of a lock and creates it.  Returns the server's representation of the lock, and an error, if there is any.
func (c *FakeLocks) Create(ctx context.Context, lock *v1beta1.Lock, opts v1.CreateOptions) (result *v1beta1.Lock, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(locksResource, lock), &v1beta1.Lock{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Lock), err
}

// Update takes the representation of a lock and updates it. Returns the server's representation of the lock, and an error, if there is any.
func (c *FakeLocks) Update(ctx context.Context, lock *v1beta1.Lock, opts v1.UpdateOptions) (result *v1beta1.Lock, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(locksResource, lock), &v1beta1.Lock{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Lock), err
}

// Delete takes name of the lock and deletes it. Returns an error if one occurs.
func (c *FakeLocks) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(locksResource, name), &v1beta1.Lock{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeLocks) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(locksResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta1.LockList{})
	return err
}

// Patch applies the patch and returns the patched lock.
func (c *FakeLocks) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.Lock, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(locksResource, name, pt, data, subresources...), &v1beta1.Lock{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Lock), err
}