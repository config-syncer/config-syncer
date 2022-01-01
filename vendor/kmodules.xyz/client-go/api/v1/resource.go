/*
Copyright AppsCode Inc. and Contributors

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// ResourceID identifies a resource
type ResourceID struct {
	Group   string `json:"group" protobuf:"bytes,1,opt,name=group"`
	Version string `json:"version" protobuf:"bytes,2,opt,name=version"`
	// Name is the plural name of the resource to serve.  It must match the name of the CustomResourceDefinition-registration
	// too: plural.group and it must be all lowercase.
	Name string `json:"name" protobuf:"bytes,3,opt,name=name"`
	// Kind is the serialized kind of the resource.  It is normally CamelCase and singular.
	Kind  string        `json:"kind" protobuf:"bytes,4,opt,name=kind"`
	Scope ResourceScope `json:"scope" protobuf:"bytes,5,opt,name=scope,casttype=ResourceScope"`
}

// ResourceScope is an enum defining the different scopes available to a custom resource
type ResourceScope string

const (
	ClusterScoped   ResourceScope = "Cluster"
	NamespaceScoped ResourceScope = "Namespaced"
)

func (r ResourceID) GroupVersion() schema.GroupVersion {
	return schema.GroupVersion{Group: r.Group, Version: r.Version}
}

func (r ResourceID) GroupResource() schema.GroupResource {
	return schema.GroupResource{Group: r.Group, Resource: r.Name}
}

func (r ResourceID) TypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{APIVersion: r.GroupVersion().String(), Kind: r.Kind}
}

func (r ResourceID) GroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{Group: r.Group, Version: r.Version, Resource: r.Name}
}

func (r ResourceID) GroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{Group: r.Group, Version: r.Version, Kind: r.Kind}
}

func (r ResourceID) MetaGVR() metav1.GroupVersionResource {
	return metav1.GroupVersionResource{Group: r.Group, Version: r.Version, Resource: r.Name}
}

func (r ResourceID) MetaGVK() metav1.GroupVersionKind {
	return metav1.GroupVersionKind{Group: r.Group, Version: r.Version, Kind: r.Kind}
}

func FromMetaGVR(in metav1.GroupVersionResource) schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    in.Group,
		Version:  in.Version,
		Resource: in.Resource,
	}
}

func ToMetaGVR(in schema.GroupVersionResource) metav1.GroupVersionResource {
	return metav1.GroupVersionResource{
		Group:    in.Group,
		Version:  in.Version,
		Resource: in.Resource,
	}
}

func FromMetaGVK(in metav1.GroupVersionKind) schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   in.Group,
		Version: in.Version,
		Kind:    in.Kind,
	}
}

func ToMetaGVK(in schema.GroupVersionKind) metav1.GroupVersionKind {
	return metav1.GroupVersionKind{
		Group:   in.Group,
		Version: in.Version,
		Kind:    in.Kind,
	}
}

// FromAPIVersionAndKind returns a GVK representing the provided fields for types that
// do not use TypeMeta. This method exists to support test types and legacy serializations
// that have a distinct group and kind.
func FromAPIVersionAndKind(apiVersion, kind string) metav1.GroupVersionKind {
	if gv, err := schema.ParseGroupVersion(apiVersion); err == nil {
		return metav1.GroupVersionKind{Group: gv.Group, Version: gv.Version, Kind: kind}
	}
	return metav1.GroupVersionKind{Kind: kind}
}

func EqualsGVK(a schema.GroupVersionKind, b metav1.GroupVersionKind) bool {
	return a.Group == b.Group &&
		a.Version == b.Version &&
		a.Kind == b.Kind
}

func EqualsGVR(a schema.GroupVersionResource, b metav1.GroupVersionResource) bool {
	return a.Group == b.Group &&
		a.Version == b.Version &&
		a.Resource == b.Resource
}
