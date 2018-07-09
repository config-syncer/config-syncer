package v1alpha1

import (
	"fmt"
	"strings"

	"github.com/appscode/kube-mon/api"
	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	crd_api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

func (p Etcd) OffshootName() string {
	return p.Name
}

func (p Etcd) OffshootLabels() map[string]string {
	return map[string]string{
		LabelDatabaseName: p.Name,
		LabelDatabaseKind: ResourceKindEtcd,
	}
}

func (p Etcd) StatefulSetLabels() map[string]string {
	labels := p.OffshootLabels()
	for key, val := range p.Labels {
		if !strings.HasPrefix(key, GenericKey+"/") && !strings.HasPrefix(key, EtcdKey+"/") {
			labels[key] = val
		}
	}
	return labels
}

func (p Etcd) StatefulSetAnnotations() map[string]string {
	annotations := make(map[string]string)
	for key, val := range p.Annotations {
		if !strings.HasPrefix(key, GenericKey+"/") && !strings.HasPrefix(key, EtcdKey+"/") {
			annotations[key] = val
		}
	}
	return annotations
}

func (p Etcd) ResourceShortCode() string {
	return ResourceCodeEtcd
}

func (p Etcd) ResourceKind() string {
	return ResourceKindEtcd
}

func (p Etcd) ResourceSingular() string {
	return ResourceSingularEtcd
}

func (p Etcd) ResourcePlural() string {
	return ResourcePluralEtcd
}

func (p Etcd) ServiceName() string {
	return p.OffshootName()
}

func (p Etcd) ServiceMonitorName() string {
	return fmt.Sprintf("kubedb-%s-%s", p.Namespace, p.Name)
}

func (p Etcd) Path() string {
	return fmt.Sprintf("/kubedb.com/v1alpha1/namespaces/%s/%s/%s/metrics", p.Namespace, p.ResourcePlural(), p.Name)
}

func (p Etcd) Scheme() string {
	return ""
}

func (p *Etcd) StatsAccessor() api.StatsAccessor {
	return p
}

func (m *Etcd) GetMonitoringVendor() string {
	if m.Spec.Monitor != nil {
		return m.Spec.Monitor.Agent.Vendor()
	}
	return ""
}

func (p Etcd) CustomResourceDefinition() *crd_api.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Version:       SchemeGroupVersion.Version,
		Plural:        ResourcePluralEtcd,
		Singular:      ResourceSingularEtcd,
		Kind:          ResourceKindEtcd,
		ShortNames:    []string{ResourceCodeEtcd},
		ResourceScope: string(apiextensions.NamespaceScoped),
		Labels: crdutils.Labels{
			LabelsMap: map[string]string{"app": "kubedb"},
		},
		SpecDefinitionName:      "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1.Etcd",
		EnableValidation:        true,
		GetOpenAPIDefinitions:   GetOpenAPIDefinitions,
		EnableStatusSubresource: EnableStatusSubresource,
	}, setNameSchema)
}
