package v1alpha1

import (
	"reflect"

	"github.com/appscode/go/log"
	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	meta_util "github.com/appscode/kutil/meta"
	"github.com/golang/glog"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

func (d DormantDatabase) OffshootName() string {
	return d.Name
}

func (d DormantDatabase) ResourceShortCode() string {
	return ResourceCodeDormantDatabase
}

func (d DormantDatabase) ResourceKind() string {
	return ResourceKindDormantDatabase
}

func (d DormantDatabase) ResourceSingular() string {
	return ResourceSingularDormantDatabase
}

func (d DormantDatabase) ResourcePlural() string {
	return ResourcePluralDormantDatabase
}

func (d DormantDatabase) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePluralDormantDatabase,
		Singular:      ResourceSingularDormantDatabase,
		Kind:          ResourceKindDormantDatabase,
		ShortNames:    []string{ResourceCodeDormantDatabase},
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
		SpecDefinitionName:      "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1.DormantDatabase",
		EnableValidation:        false,
		GetOpenAPIDefinitions:   GetOpenAPIDefinitions,
		EnableStatusSubresource: EnableStatusSubresource,
		AdditionalPrinterColumns: []apiextensions.CustomResourceColumnDefinition{
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

func (d *DormantDatabase) Migrate() {
	if d == nil {
		return
	}
	d.Spec.Origin.Spec.Elasticsearch.Migrate()
	d.Spec.Origin.Spec.Postgres.Migrate()
	d.Spec.Origin.Spec.MySQL.Migrate()
	d.Spec.Origin.Spec.MongoDB.Migrate()
	d.Spec.Origin.Spec.Redis.Migrate()
	d.Spec.Origin.Spec.Memcached.Migrate()
	d.Spec.Origin.Spec.Etcd.Migrate()
}

func (d *DormantDatabase) AlreadyObserved(other *DormantDatabase) bool {
	if d == nil {
		return other == nil
	}
	if other == nil { // && d != nil
		return false
	}
	if d == other {
		return true
	}

	var match bool

	if EnableStatusSubresource {
		match = d.Status.ObservedGeneration >= d.Generation
	} else {
		match = meta_util.Equal(d.Spec, other.Spec)
	}
	if match {
		match = reflect.DeepEqual(d.Labels, other.Labels)
	}
	if match {
		match = meta_util.EqualAnnotation(d.Annotations, other.Annotations)
	}

	if !match && bool(glog.V(log.LevelDebug)) {
		diff := meta_util.Diff(other, d)
		glog.V(log.LevelDebug).Infof("%s %s/%s has changed. Diff: %s", meta_util.GetKind(d), d.Namespace, d.Name, diff)
	}
	return match
}
