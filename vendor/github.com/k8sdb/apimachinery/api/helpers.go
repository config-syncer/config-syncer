package api

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

	PostgresKey             = ResourceNamePostgres + "." + GenericKey
	PostgresDatabaseVersion = PostgresKey + "/version"

	ElasticsearchKey             = ResourceNameElasticsearch + ".kubedb.com"
	ElasticsearchDatabaseVersion = ElasticsearchKey + "/version"

	SnapshotKey         = ResourceNameSnapshot + "s.kubedb.com"
	LabelSnapshotStatus = SnapshotKey + "/status"
)

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

func (s Snapshot) OffshootName() string {
	return s.Name
}

func (d DormantDatabase) OffshootName() string {
	return d.Name
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
