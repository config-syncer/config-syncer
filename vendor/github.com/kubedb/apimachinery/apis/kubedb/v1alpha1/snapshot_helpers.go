package v1alpha1

import (
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/appscode/go/log"
	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	meta_util "github.com/appscode/kutil/meta"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

func (s Snapshot) OffshootName() string {
	return s.Name
}

func (s Snapshot) Location() (string, error) {
	spec := s.Spec.Backend
	if spec.S3 != nil {
		return filepath.Join(spec.S3.Prefix, DatabaseNamePrefix, s.Namespace, s.Spec.DatabaseName), nil
	} else if spec.GCS != nil {
		return filepath.Join(spec.GCS.Prefix, DatabaseNamePrefix, s.Namespace, s.Spec.DatabaseName), nil
	} else if spec.Azure != nil {
		return filepath.Join(spec.Azure.Prefix, DatabaseNamePrefix, s.Namespace, s.Spec.DatabaseName), nil
	} else if spec.Local != nil {
		return filepath.Join(DatabaseNamePrefix, s.Namespace, s.Spec.DatabaseName), nil
	} else if spec.Swift != nil {
		return filepath.Join(spec.Swift.Prefix, DatabaseNamePrefix, s.Namespace, s.Spec.DatabaseName), nil
	}
	return "", errors.New("no storage provider is configured")
}

func (s Snapshot) ResourceShortCode() string {
	return ResourceCodeSnapshot
}

func (s Snapshot) ResourceKind() string {
	return ResourceKindSnapshot
}

func (s Snapshot) ResourceSingular() string {
	return ResourceSingularSnapshot
}

func (s Snapshot) ResourcePlural() string {
	return ResourcePluralSnapshot
}

func (s Snapshot) OSMSecretName() string {
	return fmt.Sprintf("osm-%v", s.Name)
}

func (s Snapshot) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePluralSnapshot,
		Singular:      ResourceSingularSnapshot,
		Kind:          ResourceKindSnapshot,
		ShortNames:    []string{ResourceCodeSnapshot},
		Categories:    []string{"datastore", "kubedb", "appscode", "all"},
		ResourceScope: string(apiextensions.NamespaceScoped),
		Versions: []apiextensions.CustomResourceDefinitionVersion{
			{
				Name:    SchemeGroupVersion.Version,
				Served:  true,
				Storage: true,
			},
		},
		Labels: crdutils.Labels{
			LabelsMap: map[string]string{"app": "kubedb"},
		},
		SpecDefinitionName:      "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1.Snapshot",
		EnableValidation:        true,
		GetOpenAPIDefinitions:   GetOpenAPIDefinitions,
		EnableStatusSubresource: EnableStatusSubresource,
		AdditionalPrinterColumns: []apiextensions.CustomResourceColumnDefinition{
			{
				Name:     "DatabaseName",
				Type:     "string",
				JSONPath: ".spec.databaseName",
			},
			{
				Name:     "Status",
				Type:     "string",
				JSONPath: ".status.phase",
			},
			{
				Name:     "Age",
				Type:     "date",
				JSONPath: ".metadata.creationTimestamp",
			},
		},
	}, setNameSchema)
}

func (s *Snapshot) Migrate() {
	if s == nil {
		return
	}
	if s.Spec.Resources != nil {
		s.Spec.PodTemplate.Spec.Resources = *s.Spec.Resources
		s.Spec.Resources = nil
	}
}

func (s *Snapshot) AlreadyObserved(other *Snapshot) bool {
	if s == nil {
		return other == nil
	}
	if other == nil { // && d != nil
		return false
	}
	if s == other {
		return true
	}

	var match bool

	if EnableStatusSubresource {
		match = s.Status.ObservedGeneration >= s.Generation
	} else {
		match = meta_util.Equal(s.Spec, other.Spec)
	}
	if match {
		match = reflect.DeepEqual(s.Labels, other.Labels)
	}
	if match {
		match = meta_util.EqualAnnotation(s.Annotations, other.Annotations)
	}

	if !match && bool(glog.V(log.LevelDebug)) {
		diff := meta_util.Diff(other, s)
		glog.V(log.LevelDebug).Infof("%s %s/%s has changed. Diff: %s", meta_util.GetKind(s), s.Namespace, s.Name, diff)
	}
	return match
}
