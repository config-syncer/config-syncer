package v1alpha1

import (
	"fmt"

	"github.com/appscode/go/types"
	apps "k8s.io/api/apps/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	crdutils "kmodules.xyz/client-go/apiextensions/v1beta1"
	v1 "kmodules.xyz/client-go/core/v1"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/apis/kubedb"
)

var _ apis.ResourceInfo = &PerconaXtraDB{}

func (p PerconaXtraDB) OffshootName() string {
	return p.Name
}

func (p PerconaXtraDB) OffshootSelectors() map[string]string {
	return map[string]string{
		LabelDatabaseName: p.Name,
		LabelDatabaseKind: ResourceKindPerconaXtraDB,
	}
}

func (p PerconaXtraDB) OffshootLabels() map[string]string {
	out := p.OffshootSelectors()
	out[meta_util.NameLabelKey] = ResourceSingularPerconaXtraDB
	out[meta_util.VersionLabelKey] = string(p.Spec.Version)
	out[meta_util.InstanceLabelKey] = p.Name
	out[meta_util.ComponentLabelKey] = "database"
	out[meta_util.ManagedByLabelKey] = GenericKey
	return meta_util.FilterKeys(GenericKey, out, p.Labels)
}

func (p PerconaXtraDB) ResourceShortCode() string {
	return ResourceCodePerconaXtraDB
}

func (p PerconaXtraDB) ResourceKind() string {
	return ResourceKindPerconaXtraDB
}

func (p PerconaXtraDB) ResourceSingular() string {
	return ResourceSingularPerconaXtraDB
}

func (p PerconaXtraDB) ResourcePlural() string {
	return ResourcePluralPerconaXtraDB
}

func (p PerconaXtraDB) ServiceName() string {
	return p.OffshootName()
}

func (p PerconaXtraDB) GoverningServiceName() string {
	return p.OffshootName() + "-gvr"
}

func (p PerconaXtraDB) PeerName(idx int) string {
	return fmt.Sprintf("%s-%d.%s.%s", p.OffshootName(), idx, p.GoverningServiceName(), p.Namespace)
}

func (p PerconaXtraDB) ClusterName() string {
	return p.Spec.PXC.ClusterName
}

func (p PerconaXtraDB) ClusterLabels() map[string]string {
	return v1.UpsertMap(p.OffshootLabels(), map[string]string{
		PerconaXtraDBClusterLabelKey: p.ClusterName(),
	})
}

func (p PerconaXtraDB) ClusterSelectors() map[string]string {
	return v1.UpsertMap(p.OffshootSelectors(), map[string]string{
		PerconaXtraDBClusterLabelKey: p.ClusterName(),
	})
}

func (p PerconaXtraDB) XtraDBLabels() map[string]string {
	if p.Spec.PXC != nil {
		return p.ClusterLabels()
	}
	return p.OffshootLabels()
}

func (p PerconaXtraDB) XtraDBSelectors() map[string]string {
	if p.Spec.PXC != nil {
		return p.ClusterSelectors()
	}
	return p.OffshootSelectors()
}

func (p PerconaXtraDB) ProxysqlName() string {
	return fmt.Sprintf("%s-proxysql", p.OffshootName())
}

func (p PerconaXtraDB) ProxysqlServiceName() string {
	return p.ProxysqlName()
}

func (p PerconaXtraDB) ProxysqlLabels() map[string]string {
	return v1.UpsertMap(p.OffshootLabels(), map[string]string{
		PerconaXtraDBProxysqlLabelKey: p.ProxysqlName(),
	})
}

func (p PerconaXtraDB) ProxysqlSelectors() map[string]string {
	return v1.UpsertMap(p.OffshootSelectors(), map[string]string{
		PerconaXtraDBProxysqlLabelKey: p.ProxysqlName(),
	})
}

type perconaXtraDBApp struct {
	*PerconaXtraDB
}

func (p perconaXtraDBApp) Name() string {
	return p.PerconaXtraDB.Name
}

func (p perconaXtraDBApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularPerconaXtraDB))
}

func (p PerconaXtraDB) AppBindingMeta() appcat.AppBindingMeta {
	return &perconaXtraDBApp{&p}
}

type perconaXtraDBStatsService struct {
	*PerconaXtraDB
}

func (p perconaXtraDBStatsService) GetNamespace() string {
	return p.PerconaXtraDB.GetNamespace()
}

func (p perconaXtraDBStatsService) ServiceName() string {
	return p.OffshootName() + "-stats"
}

func (p perconaXtraDBStatsService) ServiceMonitorName() string {
	return fmt.Sprintf("kubedb-%s-%s", p.Namespace, p.Name)
}

func (p perconaXtraDBStatsService) Path() string {
	return "/metrics"
}

func (p perconaXtraDBStatsService) Scheme() string {
	return ""
}

func (p PerconaXtraDB) StatsService() mona.StatsAccessor {
	return &perconaXtraDBStatsService{&p}
}

func (p PerconaXtraDB) StatsServiceLabels() map[string]string {
	lbl := meta_util.FilterKeys(GenericKey, p.OffshootSelectors(), p.Labels)
	lbl[LabelRole] = "stats"
	return lbl
}

func (p *PerconaXtraDB) GetMonitoringVendor() string {
	if p.Spec.Monitor != nil {
		return p.Spec.Monitor.Agent.Vendor()
	}
	return ""
}

func (p PerconaXtraDB) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePluralPerconaXtraDB,
		Singular:      ResourceSingularPerconaXtraDB,
		Kind:          ResourceKindPerconaXtraDB,
		ShortNames:    []string{ResourceCodePerconaXtraDB},
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
		SpecDefinitionName:      "kubedb.dev/apimachinery/apis/kubedb/v1alpha1.PerconaXtraDB",
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

func (p *PerconaXtraDB) SetDefaults() {
	if p == nil {
		return
	}
	p.Spec.SetDefaults()
}

func (p *PerconaXtraDBSpec) SetDefaults() {
	if p == nil {
		return
	}

	if p.Replicas == nil {
		p.Replicas = types.Int32P(1)
	}

	if p.PXC != nil {
		if *p.Replicas < 3 {
			p.Replicas = types.Int32P(PerconaXtraDBDefaultClusterSize)
		}

		if p.PXC.Proxysql.Replicas == nil {
			p.PXC.Proxysql.Replicas = types.Int32P(1)
		}
	}

	if p.StorageType == "" {
		p.StorageType = StorageTypeDurable
	}
	if p.UpdateStrategy.Type == "" {
		p.UpdateStrategy.Type = apps.RollingUpdateStatefulSetStrategyType
	}
	if p.TerminationPolicy == "" {
		if p.StorageType == StorageTypeEphemeral {
			p.TerminationPolicy = TerminationPolicyDelete
		} else {
			p.TerminationPolicy = TerminationPolicyPause
		}
	}
}

func (p *PerconaXtraDBSpec) GetSecrets() []string {
	if p == nil {
		return nil
	}

	var secrets []string
	if p.DatabaseSecret != nil {
		secrets = append(secrets, p.DatabaseSecret.SecretName)
	}
	return secrets
}
