package api

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupName is the group name use in this package
const GroupName = "kubedb.com"

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

// Adds the list of known types to metav1.Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		// Snapshot
		&Snapshot{},
		&SnapshotList{},
		// DormantDatabase
		&DormantDatabase{},
		&DormantDatabaseList{},
		// kubedb Elastic
		&Elastic{},
		&ElasticList{},
		// kubedb Postgres
		&Postgres{},
		&PostgresList{},
	)
	return nil
}

func (s *Snapshot) GetObjectKind() schema.ObjectKind       { return &s.TypeMeta }
func (obj *SnapshotList) GetObjectKind() schema.ObjectKind { return &obj.TypeMeta }

func (d *DormantDatabase) GetObjectKind() schema.ObjectKind       { return &d.TypeMeta }
func (obj *DormantDatabaseList) GetObjectKind() schema.ObjectKind { return &obj.TypeMeta }

func (e *Elastic) GetObjectKind() schema.ObjectKind       { return &e.TypeMeta }
func (obj *ElasticList) GetObjectKind() schema.ObjectKind { return &obj.TypeMeta }

func (p *Postgres) GetObjectKind() schema.ObjectKind       { return &p.TypeMeta }
func (obj *PostgresList) GetObjectKind() schema.ObjectKind { return &obj.TypeMeta }
