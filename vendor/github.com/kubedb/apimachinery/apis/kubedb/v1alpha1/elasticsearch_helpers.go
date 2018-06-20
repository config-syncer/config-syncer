package v1alpha1

import (
	"fmt"
	"strings"

	"github.com/appscode/kube-mon/api"
	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	"github.com/appscode/kutil/meta"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	crd_api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

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
	return annotations
}

var _ ResourceInfo = &Elasticsearch{}

func (e Elasticsearch) ResourceShortCode() string {
	return ResourceCodeElasticsearch
}

func (e Elasticsearch) ResourceKind() string {
	return ResourceKindElasticsearch
}

func (e Elasticsearch) ResourceSingular() string {
	return ResourceSingularElasticsearch
}

func (e Elasticsearch) ResourcePlural() string {
	return ResourcePluralElasticsearch
}

func (r Elasticsearch) ServiceName() string {
	return r.OffshootName()
}

func (r *Elasticsearch) MasterServiceName() string {
	return fmt.Sprintf("%v-master", r.ServiceName())
}

func (r Elasticsearch) ServiceMonitorName() string {
	return fmt.Sprintf("kubedb-%s-%s", r.Namespace, r.Name)
}

func (r Elasticsearch) Path() string {
	return fmt.Sprintf("/kubedb.com/v1alpha1/namespaces/%s/%s/%s/metrics", r.Namespace, r.ResourcePlural(), r.Name)
}

func (r Elasticsearch) Scheme() string {
	return ""
}

func (r *Elasticsearch) StatsAccessor() api.StatsAccessor {
	return r
}

func (e *Elasticsearch) GetMonitoringVendor() string {
	if e.Spec.Monitor != nil {
		return e.Spec.Monitor.Agent.Vendor()
	}
	return ""
}

func (r Elasticsearch) CustomResourceDefinition() *crd_api.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Version:       SchemeGroupVersion.Version,
		Plural:        ResourcePluralElasticsearch,
		Singular:      ResourceSingularElasticsearch,
		Kind:          ResourceKindElasticsearch,
		ShortNames:    []string{ResourceCodeElasticsearch},
		ResourceScope: string(apiextensions.NamespaceScoped),
		Labels: crdutils.Labels{
			LabelsMap: map[string]string{"app": "kubedb"},
		},
		SpecDefinitionName:    "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1.Elasticsearch",
		EnableValidation:      true,
		GetOpenAPIDefinitions: GetOpenAPIDefinitions,
	}, setNameSchema)
}

const (
	ESSearchGuardDisabled = ElasticsearchKey + "/searchguard-disabled"
)

func (r Elasticsearch) SearchGuardDisabled() bool {
	v, _ := meta.GetBoolValue(r.Annotations, ESSearchGuardDisabled)
	return v
}
