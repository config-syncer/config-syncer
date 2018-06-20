package v1alpha1

import (
	"fmt"
	"path/filepath"

	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	"github.com/pkg/errors"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	crd_api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

func (s Snapshot) OffshootName() string {
	return s.Name
}

func (s Snapshot) Location() (string, error) {
	spec := s.Spec.SnapshotStorageSpec
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

func (s SnapshotStorageSpec) Container() (string, error) {
	if s.S3 != nil {
		return s.S3.Bucket, nil
	} else if s.GCS != nil {
		return s.GCS.Bucket, nil
	} else if s.Azure != nil {
		return s.Azure.Container, nil
	} else if s.Local != nil {
		return s.Local.MountPath, nil
	} else if s.Swift != nil {
		return s.Swift.Container, nil
	}
	return "", errors.New("no storage provider is configured")
}

func (s SnapshotStorageSpec) Location() (string, error) {
	if s.S3 != nil {
		return "s3:" + s.S3.Bucket, nil
	} else if s.GCS != nil {
		return "gs:" + s.GCS.Bucket, nil
	} else if s.Azure != nil {
		return "azure:" + s.Azure.Container, nil
	} else if s.Local != nil {
		return "local:" + s.Local.MountPath, nil
	} else if s.Swift != nil {
		return "swift:" + s.Swift.Container, nil
	}
	return "", errors.New("no storage provider is configured")
}

func (s Snapshot) OSMSecretName() string {
	return fmt.Sprintf("osm-%v", s.Name)
}

func (s Snapshot) CustomResourceDefinition() *crd_api.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Version:       SchemeGroupVersion.Version,
		Plural:        ResourcePluralSnapshot,
		Singular:      ResourceSingularSnapshot,
		Kind:          ResourceKindSnapshot,
		ShortNames:    []string{ResourceCodeSnapshot},
		ResourceScope: string(apiextensions.NamespaceScoped),
		Labels: crdutils.Labels{
			LabelsMap: map[string]string{"app": "kubedb"},
		},
		SpecDefinitionName:    "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1.Snapshot",
		EnableValidation:      true,
		GetOpenAPIDefinitions: GetOpenAPIDefinitions,
	}, setNameSchema)
}
