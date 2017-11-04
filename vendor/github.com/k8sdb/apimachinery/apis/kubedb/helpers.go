package kubedb

import (
	"errors"
	"path/filepath"
	"strings"
)

const (
	DatabaseNamePrefix = "kubedb"

	GenericKey = "kubedb.com"

	LabelDatabaseKind = GenericKey + "/kind"
	LabelDatabaseName = GenericKey + "/name"
	LabelJobType      = GenericKey + "/job-type"

	PostgresKey             = ResourceTypePostgres + "." + GenericKey
	PostgresDatabaseVersion = PostgresKey + "/version"

	ElasticsearchKey             = ResourceTypeElasticsearch + "." + GenericKey
	ElasticsearchDatabaseVersion = ElasticsearchKey + "/version"

	MySQLKey             = ResourceTypeMySQL + "." + GenericKey
	MySQLDatabaseVersion = MySQLKey + "/version"

	SnapshotKey         = ResourceTypeSnapshot + "." + GenericKey
	LabelSnapshotStatus = SnapshotKey + "/status"

	PostgresInitSpec      = PostgresKey + "/init"
	ElasticsearchInitSpec = ElasticsearchKey + "/init"
	MySQLInitSpec         = MySQLKey + "/init"

	PostgresIgnore      = PostgresKey + "/ignore"
	ElasticsearchIgnore = ElasticsearchKey + "/ignore"
	MySQLIgnore         = MySQLKey + "/ignore"
)

type RuntimeObject interface {
	ResourceCode() string
	ResourceKind() string
	ResourceName() string
	ResourceType() string
}

func (p Postgres) OffshootName() string {
	return p.Name
}

func (p Postgres) OffshootLabels() map[string]string {
	return map[string]string{
		LabelDatabaseName: p.Name,
		LabelDatabaseKind: ResourceKindPostgres,
	}
}

func (p Postgres) StatefulSetLabels() map[string]string {
	labels := p.OffshootLabels()
	for key, val := range p.Labels {
		if !strings.HasPrefix(key, GenericKey+"/") && !strings.HasPrefix(key, PostgresKey+"/") {
			labels[key] = val
		}
	}
	return labels
}

func (p Postgres) StatefulSetAnnotations() map[string]string {
	annotations := make(map[string]string)
	for key, val := range p.Annotations {
		if !strings.HasPrefix(key, GenericKey+"/") && !strings.HasPrefix(key, PostgresKey+"/") {
			annotations[key] = val
		}
	}
	annotations[PostgresDatabaseVersion] = string(p.Spec.Version)
	return annotations
}

func (p Postgres) ResourceCode() string {
	return ResourceCodePostgres
}

func (p Postgres) ResourceKind() string {
	return ResourceKindPostgres
}

func (p Postgres) ResourceName() string {
	return ResourceNamePostgres
}

func (p Postgres) ResourceType() string {
	return ResourceTypePostgres
}

func (m MySQL) OffshootName() string {
	return m.Name
}

func (m MySQL) OffshootLabels() map[string]string {
	return map[string]string{
		LabelDatabaseName: m.Name,
		LabelDatabaseKind: ResourceKindMySQL,
	}
}

func (m MySQL) StatefulSetLabels() map[string]string {
	labels := m.OffshootLabels()
	for key, val := range m.Labels {
		if !strings.HasPrefix(key, GenericKey+"/") && !strings.HasPrefix(key, MySQLKey+"/") {
			labels[key] = val
		}
	}
	return labels
}

func (m MySQL) StatefulSetAnnotations() map[string]string {
	annotations := make(map[string]string)
	for key, val := range m.Annotations {
		if !strings.HasPrefix(key, GenericKey+"/") && !strings.HasPrefix(key, MySQLKey+"/") {
			annotations[key] = val
		}
	}
	annotations[MySQLDatabaseVersion] = string(m.Spec.Version)
	return annotations
}

func (m MySQL) ResourceCode() string {
	return ResourceCodeMySQL
}

func (m MySQL) ResourceKind() string {
	return ResourceKindMySQL
}

func (m MySQL) ResourceName() string {
	return ResourceNameMySQL
}

func (m MySQL) ResourceType() string {
	return ResourceTypeMySQL
}

func (e Elasticsearch) OffshootName() string {
	return e.Name
}

func (e Elasticsearch) OffshootLabels() map[string]string {
	return map[string]string{
		LabelDatabaseKind: ResourceKindElasticsearch,
		LabelDatabaseName: e.Name,
	}
}

func (e Elasticsearch) StatefulSetLabels() map[string]string {
	labels := e.OffshootLabels()
	for key, val := range e.Labels {
		if !strings.HasPrefix(key, GenericKey+"/") && !strings.HasPrefix(key, ElasticsearchKey+"/") {
			labels[key] = val
		}
	}
	return labels
}

func (e Elasticsearch) StatefulSetAnnotations() map[string]string {
	annotations := make(map[string]string)
	for key, val := range e.Annotations {
		if !strings.HasPrefix(key, GenericKey+"/") && !strings.HasPrefix(key, ElasticsearchKey+"/") {
			annotations[key] = val
		}
	}
	annotations[ElasticsearchDatabaseVersion] = string(e.Spec.Version)
	return annotations
}

func (p Elasticsearch) ResourceCode() string {
	return ResourceCodeElasticsearch
}

func (p Elasticsearch) ResourceKind() string {
	return ResourceKindElasticsearch
}

func (p Elasticsearch) ResourceName() string {
	return ResourceNameElasticsearch
}

func (p Elasticsearch) ResourceType() string {
	return ResourceTypeElasticsearch
}

func (d DormantDatabase) OffshootName() string {
	return d.Name
}

func (p DormantDatabase) ResourceCode() string {
	return ResourceCodeDormantDatabase
}

func (p DormantDatabase) ResourceKind() string {
	return ResourceKindDormantDatabase
}

func (p DormantDatabase) ResourceName() string {
	return ResourceNameDormantDatabase
}

func (p DormantDatabase) ResourceType() string {
	return ResourceTypeDormantDatabase
}

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
	return "", errors.New("No storage provider is configured.")
}

func (p Snapshot) ResourceCode() string {
	return ResourceCodeSnapshot
}

func (p Snapshot) ResourceKind() string {
	return ResourceKindSnapshot
}

func (p Snapshot) ResourceName() string {
	return ResourceNameSnapshot
}

func (p Snapshot) ResourceType() string {
	return ResourceTypeSnapshot
}

func (s SnapshotStorageSpec) Container() (string, error) {
	if s.S3 != nil {
		return s.S3.Bucket, nil
	} else if s.GCS != nil {
		return s.GCS.Bucket, nil
	} else if s.Azure != nil {
		return s.Azure.Container, nil
	} else if s.Local != nil {
		return s.Local.Path, nil
	} else if s.Swift != nil {
		return s.Swift.Container, nil
	}
	return "", errors.New("No storage provider is configured.")
}

func (s SnapshotStorageSpec) Location() (string, error) {
	if s.S3 != nil {
		return "s3:" + s.S3.Bucket, nil
	} else if s.GCS != nil {
		return "gs:" + s.GCS.Bucket, nil
	} else if s.Azure != nil {
		return "azure:" + s.Azure.Container, nil
	} else if s.Local != nil {
		return "local:" + s.Local.Path, nil
	} else if s.Swift != nil {
		return "swift:" + s.Swift.Container, nil
	}
	return "", errors.New("No storage provider is configured.")
}
