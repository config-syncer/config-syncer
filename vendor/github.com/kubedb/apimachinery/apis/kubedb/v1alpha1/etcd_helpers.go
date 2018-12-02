package v1alpha1

import (
	"fmt"

	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	meta_util "github.com/appscode/kutil/meta"
	"github.com/kubedb/apimachinery/apis"
	"github.com/kubedb/apimachinery/apis/kubedb"
	apps "k8s.io/api/apps/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
)

var _ apis.ResourceInfo = &Etcd{}

func (e Etcd) OffshootName() string {
	return e.Name
}

func (e Etcd) OffshootSelectors() map[string]string {
	return map[string]string{
		LabelDatabaseName: e.Name,
		LabelDatabaseKind: ResourceKindEtcd,
	}
}

func (e Etcd) OffshootLabels() map[string]string {
	return meta_util.FilterKeys(GenericKey, e.OffshootSelectors(), e.Labels)
}

func (e Etcd) ResourceShortCode() string {
	return ResourceCodeEtcd
}

func (e Etcd) ResourceKind() string {
	return ResourceKindEtcd
}

func (e Etcd) ResourceSingular() string {
	return ResourceSingularEtcd
}

func (e Etcd) ResourcePlural() string {
	return ResourcePluralEtcd
}

func (e Etcd) ClientServiceName() string {
	return e.OffshootName() + "-client"
}

func (e Etcd) PeerServiceName() string {
	return e.OffshootName()
}

type etcdApp struct {
	*Etcd
}

func (r etcdApp) Name() string {
	return r.Etcd.Name
}

func (r etcdApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularEtcd))
}

func (r Etcd) AppBindingMeta() appcat.AppBindingMeta {
	return &etcdApp{&r}
}

type etcdStatsService struct {
	*Etcd
}

func (e etcdStatsService) GetNamespace() string {
	return e.Etcd.GetNamespace()
}

func (e etcdStatsService) ServiceName() string {
	return e.OffshootName() + "-stats"
}

func (e etcdStatsService) ServiceMonitorName() string {
	return fmt.Sprintf("kubedb-%s-%s", e.Namespace, e.Name)
}

func (e etcdStatsService) Path() string {
	return fmt.Sprintf("/metrics")
}

func (e etcdStatsService) Scheme() string {
	return ""
}

func (e Etcd) StatsService() mona.StatsAccessor {
	return &etcdStatsService{&e}
}

func (e *Etcd) GetMonitoringVendor() string {
	if e.Spec.Monitor != nil {
		return e.Spec.Monitor.Agent.Vendor()
	}
	return ""
}

func (e Etcd) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePluralEtcd,
		Singular:      ResourceSingularEtcd,
		Kind:          ResourceKindEtcd,
		ShortNames:    []string{ResourceCodeEtcd},
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
		SpecDefinitionName:      "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1.Etcd",
		EnableValidation:        true,
		GetOpenAPIDefinitions:   GetOpenAPIDefinitions,
		EnableStatusSubresource: apis.EnableStatusSubresource,
		AdditionalPrinterColumns: []apiextensions.CustomResourceColumnDefinition{
			{
				Name:     "Version",
				Type:     "string",
				JSONPath: ".spec.version",
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
	}, apis.SetNameSchema)
}

func (e *Etcd) SetDefaults() {
	if e == nil {
		return
	}
	e.Spec.SetDefaults()
}

func (e *EtcdSpec) SetDefaults() {
	if e == nil {
		return
	}

	// perform defaulting
	if e.StorageType == "" {
		e.StorageType = StorageTypeDurable
	}
	if e.UpdateStrategy.Type == "" {
		e.UpdateStrategy.Type = apps.RollingUpdateStatefulSetStrategyType
	}
	if e.TerminationPolicy == "" {
		if e.StorageType == StorageTypeEphemeral {
			e.TerminationPolicy = TerminationPolicyDelete
		} else {
			e.TerminationPolicy = TerminationPolicyPause
		}
	}
}

func (e *EtcdSpec) GetSecrets() []string {
	return nil
}
