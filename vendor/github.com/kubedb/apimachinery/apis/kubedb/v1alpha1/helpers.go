package v1alpha1

import (
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type ResourceInfo interface {
	ResourceCode() string
	ResourceKind() string
	ResourceName() string
	ResourceType() string
}

func ObjectReferenceFor(obj runtime.Object) *core.ObjectReference {
	switch u := obj.(type) {
	case *DormantDatabase:
		return u.ObjectReference()
	case *Postgres:
		return u.ObjectReference()
	case *Elasticsearch:
		return u.ObjectReference()
	case *MySQL:
		return u.ObjectReference()
	case *MongoDB:
		return u.ObjectReference()
	case *Redis:
		return u.ObjectReference()
	case *Memcached:
		return u.ObjectReference()
	case *Snapshot:
		return u.ObjectReference()
	}
	return &core.ObjectReference{}
}
