package v1alpha1

import (
	"fmt"
	"strings"

	"github.com/appscode/kube-mon/api"
	core "k8s.io/api/core/v1"
	crd_api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	annotations[ElasticsearchDatabaseVersion] = string(e.Spec.Version)
	return annotations
}

var _ ResourceInfo = &Elasticsearch{}

func (e Elasticsearch) ResourceCode() string {
	return ResourceCodeElasticsearch
}

func (e Elasticsearch) ResourceKind() string {
	return ResourceKindElasticsearch
}

func (e Elasticsearch) ResourceName() string {
	return ResourceNameElasticsearch
}

func (e Elasticsearch) ResourceType() string {
	return ResourceTypeElasticsearch
}

func (e Elasticsearch) ObjectReference() *core.ObjectReference {
	return &core.ObjectReference{
		APIVersion:      SchemeGroupVersion.String(),
		Kind:            e.ResourceKind(),
		Namespace:       e.Namespace,
		Name:            e.Name,
		UID:             e.UID,
		ResourceVersion: e.ResourceVersion,
	}
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
	return fmt.Sprintf("/kubedb.com/v1alpha1/namespaces/%s/%s/%s/metrics", r.Namespace, r.ResourceType(), r.Name)
}

func (r Elasticsearch) Scheme() string {
	return ""
}

func (r *Elasticsearch) StatsAccessor() api.StatsAccessor {
	return r
}

func (r Elasticsearch) CustomResourceDefinition() *crd_api.CustomResourceDefinition {
	resourceName := ResourceTypeElasticsearch + "." + SchemeGroupVersion.Group

	return &crd_api.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: resourceName,
			Labels: map[string]string{
				"app": "kubedb",
			},
		},
		Spec: crd_api.CustomResourceDefinitionSpec{
			Group:   SchemeGroupVersion.Group,
			Version: SchemeGroupVersion.Version,
			Scope:   crd_api.NamespaceScoped,
			Names: crd_api.CustomResourceDefinitionNames{
				Plural:     ResourceTypeElasticsearch,
				Kind:       ResourceKindElasticsearch,
				ShortNames: []string{ResourceCodeElasticsearch},
			},
		},
	}
}
