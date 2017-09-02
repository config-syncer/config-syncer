package clientset

import (
	aci "github.com/k8sdb/apimachinery/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

type SnapshotNamespacer interface {
	Snapshots(namespace string) SnapshotInterface
}

type SnapshotInterface interface {
	List(opts metav1.ListOptions) (*aci.SnapshotList, error)
	Get(name string) (*aci.Snapshot, error)
	Create(snapshot *aci.Snapshot) (*aci.Snapshot, error)
	Update(snapshot *aci.Snapshot) (*aci.Snapshot, error)
	Delete(name string) error
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	UpdateStatus(snapshot *aci.Snapshot) (*aci.Snapshot, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (*aci.Snapshot, error)
}

type SnapshotImpl struct {
	r  rest.Interface
	ns string
}

var _ SnapshotInterface = &SnapshotImpl{}

func newSnapshot(c *ExtensionClient, namespace string) *SnapshotImpl {
	return &SnapshotImpl{c.restClient, namespace}
}

func (c *SnapshotImpl) List(opts metav1.ListOptions) (result *aci.SnapshotList, err error) {
	result = &aci.SnapshotList{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource(aci.ResourceTypeSnapshot).
		VersionedParams(&opts, ExtendedCodec).
		Do().
		Into(result)
	return
}

func (c *SnapshotImpl) Get(name string) (result *aci.Snapshot, err error) {
	result = &aci.Snapshot{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource(aci.ResourceTypeSnapshot).
		Name(name).
		Do().
		Into(result)
	return
}

func (c *SnapshotImpl) Create(snapshot *aci.Snapshot) (result *aci.Snapshot, err error) {
	result = &aci.Snapshot{}
	err = c.r.Post().
		Namespace(c.ns).
		Resource(aci.ResourceTypeSnapshot).
		Body(snapshot).
		Do().
		Into(result)
	return
}

func (c *SnapshotImpl) Update(snapshot *aci.Snapshot) (result *aci.Snapshot, err error) {
	result = &aci.Snapshot{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource(aci.ResourceTypeSnapshot).
		Name(snapshot.Name).
		Body(snapshot).
		Do().
		Into(result)
	return
}

func (c *SnapshotImpl) Delete(name string) (err error) {
	return c.r.Delete().
		Namespace(c.ns).
		Resource(aci.ResourceTypeSnapshot).
		Name(name).
		Do().
		Error()
}

func (c *SnapshotImpl) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return c.r.Get().
		Prefix("watch").
		Namespace(c.ns).
		Resource(aci.ResourceTypeSnapshot).
		VersionedParams(&opts, ExtendedCodec).
		Watch()
}

func (c *SnapshotImpl) UpdateStatus(snapshot *aci.Snapshot) (result *aci.Snapshot, err error) {
	result = &aci.Snapshot{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource(aci.ResourceTypeSnapshot).
		Name(snapshot.Name).
		SubResource("status").
		Body(snapshot).
		Do().
		Into(result)
	return
}

func (c *SnapshotImpl) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *aci.Snapshot, err error) {
	result = &aci.Snapshot{}
	err = c.r.Patch(pt).
		Namespace(c.ns).
		Resource(aci.ResourceTypeSnapshot).
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
