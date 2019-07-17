package v1alpha1

import (
	"fmt"

	"github.com/appscode/go/types"
	"github.com/kubedb/apimachinery/apis"
	"github.com/kubedb/apimachinery/apis/kubedb"
	apps "k8s.io/api/apps/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	crdutils "kmodules.xyz/client-go/apiextensions/v1beta1"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
)

var _ apis.ResourceInfo = &Percona{}

func (p Percona) OffshootName() string {
	return p.Name
}

func (p Percona) OffshootSelectors() map[string]string {
	return map[string]string{
		LabelDatabaseName: p.Name,
		LabelDatabaseKind: ResourceKindPercona,
	}
}

func (p Percona) OffshootLabels() map[string]string {
	out := p.OffshootSelectors()
	out[meta_util.NameLabelKey] = ResourceSingularPercona
	out[meta_util.VersionLabelKey] = string(p.Spec.Version)
	out[meta_util.InstanceLabelKey] = p.Name
	out[meta_util.ComponentLabelKey] = "database"
	out[meta_util.ManagedByLabelKey] = GenericKey
	return meta_util.FilterKeys(GenericKey, out, p.Labels)
}

func (p Percona) ResourceShortCode() string {
	return ResourceCodePercona
}

func (p Percona) ResourceKind() string {
	return ResourceKindPercona
}

func (p Percona) ResourceSingular() string {
	return ResourceSingularPercona
}

func (p Percona) ResourcePlural() string {
	return ResourcePluralPercona
}

func (p Percona) ServiceName() string {
	return p.OffshootName()
}

func (p Percona) GoverningServiceName() string {
	return p.OffshootName() + "-gvr"
}

type perconaApp struct {
	*Percona
}

func (p perconaApp) Name() string {
	return p.Percona.Name
}

func (p perconaApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularPercona))
}

func (p Percona) AppBindingMeta() appcat.AppBindingMeta {
	return &perconaApp{&p}
}

type perconaStatsService struct {
	*Percona
}

func (p perconaStatsService) GetNamespace() string {
	return p.Percona.GetNamespace()
}

func (p perconaStatsService) ServiceName() string {
	return p.OffshootName() + "-stats"
}

func (p perconaStatsService) ServiceMonitorName() string {
	return fmt.Sprintf("kubedb-%s-%s", p.Namespace, p.Name)
}

func (p perconaStatsService) Path() string {
	return "/metrics"
}

func (p perconaStatsService) Scheme() string {
	return ""
}

func (p Percona) StatsService() mona.StatsAccessor {
	return &perconaStatsService{&p}
}

func (p Percona) StatsServiceLabels() map[string]string {
	lbl := meta_util.FilterKeys(GenericKey, p.OffshootSelectors(), p.Labels)
	lbl[LabelRole] = "stats"
	return lbl
}

func (p *Percona) GetMonitoringVendor() string {
	if p.Spec.Monitor != nil {
		return p.Spec.Monitor.Agent.Vendor()
	}
	return ""
}

func (p Percona) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePluralPercona,
		Singular:      ResourceSingularPercona,
		Kind:          ResourceKindPercona,
		ShortNames:    []string{ResourceCodePercona},
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
		SpecDefinitionName:      "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1.Percona",
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

func (p *Percona) SetDefaults() {
	if p == nil {
		return
	}
	p.Spec.SetDefaults()
}

func (p *PerconaSpec) SetDefaults() {
	if p == nil {
		return
	}

	if p.Replicas == nil {
		p.Replicas = types.Int32P(1)
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

func (p *PerconaSpec) GetSecrets() []string {
	if p == nil {
		return nil
	}

	var secrets []string
	if p.DatabaseSecret != nil {
		secrets = append(secrets, p.DatabaseSecret.SecretName)
	}
	return secrets
}
