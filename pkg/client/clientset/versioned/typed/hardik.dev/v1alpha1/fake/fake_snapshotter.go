/*
Copyright The Kubernetes Authors.

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

	v1alpha1 "github.com/hdkshingala/snapshotter/pkg/apis/hardik.dev/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeSnapshotters implements SnapshotterInterface
type FakeSnapshotters struct {
	Fake *FakeHardikV1alpha1
	ns   string
}

var snapshottersResource = schema.GroupVersionResource{Group: "hardik.dev", Version: "v1alpha1", Resource: "snapshotters"}

var snapshottersKind = schema.GroupVersionKind{Group: "hardik.dev", Version: "v1alpha1", Kind: "Snapshotter"}

// Get takes name of the snapshotter, and returns the corresponding snapshotter object, and an error if there is any.
func (c *FakeSnapshotters) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Snapshotter, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(snapshottersResource, c.ns, name), &v1alpha1.Snapshotter{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Snapshotter), err
}

// List takes label and field selectors, and returns the list of Snapshotters that match those selectors.
func (c *FakeSnapshotters) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.SnapshotterList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(snapshottersResource, snapshottersKind, c.ns, opts), &v1alpha1.SnapshotterList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.SnapshotterList{ListMeta: obj.(*v1alpha1.SnapshotterList).ListMeta}
	for _, item := range obj.(*v1alpha1.SnapshotterList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested snapshotters.
func (c *FakeSnapshotters) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(snapshottersResource, c.ns, opts))

}

// Create takes the representation of a snapshotter and creates it.  Returns the server's representation of the snapshotter, and an error, if there is any.
func (c *FakeSnapshotters) Create(ctx context.Context, snapshotter *v1alpha1.Snapshotter, opts v1.CreateOptions) (result *v1alpha1.Snapshotter, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(snapshottersResource, c.ns, snapshotter), &v1alpha1.Snapshotter{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Snapshotter), err
}

// Update takes the representation of a snapshotter and updates it. Returns the server's representation of the snapshotter, and an error, if there is any.
func (c *FakeSnapshotters) Update(ctx context.Context, snapshotter *v1alpha1.Snapshotter, opts v1.UpdateOptions) (result *v1alpha1.Snapshotter, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(snapshottersResource, c.ns, snapshotter), &v1alpha1.Snapshotter{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Snapshotter), err
}

// Delete takes name of the snapshotter and deletes it. Returns an error if one occurs.
func (c *FakeSnapshotters) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(snapshottersResource, c.ns, name, opts), &v1alpha1.Snapshotter{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSnapshotters) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(snapshottersResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.SnapshotterList{})
	return err
}

// Patch applies the patch and returns the patched snapshotter.
func (c *FakeSnapshotters) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Snapshotter, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(snapshottersResource, c.ns, name, pt, data, subresources...), &v1alpha1.Snapshotter{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Snapshotter), err
}
