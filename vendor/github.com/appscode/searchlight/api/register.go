package api

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupName is the group name use in this package
const GroupName = "monitoring.appscode.com"

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: runtime.APIVersionInternal}

// Kind takes an unqualified kind and returns back a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns back a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

// Adds the list of known types to apiv1.Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&PodAlert{},
		&PodAlertList{},

		&NodeAlert{},
		&NodeAlertList{},

		&ClusterAlert{},
		&ClusterAlertList{},
	)
	return nil
}

func (a *PodAlert) GetObjectKind() schema.ObjectKind       { return &a.TypeMeta }
func (obj *PodAlertList) GetObjectKind() schema.ObjectKind { return &obj.TypeMeta }

func (a *NodeAlert) GetObjectKind() schema.ObjectKind       { return &a.TypeMeta }
func (obj *NodeAlertList) GetObjectKind() schema.ObjectKind { return &obj.TypeMeta }

func (a *ClusterAlert) GetObjectKind() schema.ObjectKind       { return &a.TypeMeta }
func (obj *ClusterAlertList) GetObjectKind() schema.ObjectKind { return &obj.TypeMeta }
