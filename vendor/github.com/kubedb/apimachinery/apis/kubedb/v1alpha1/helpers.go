package v1alpha1

import (
	"github.com/go-openapi/spec"
	"k8s.io/kube-openapi/pkg/common"
)

var (
	EnableStatusSubresource bool
)

type ResourceInfo interface {
	ResourceShortCode() string
	ResourceKind() string
	ResourceSingular() string
	ResourcePlural() string
}

func setNameSchema(openapiSpec map[string]common.OpenAPIDefinition) {
	// ref: https://github.com/kubedb/project/issues/166
	// https://github.com/kubernetes/apimachinery/blob/94ebb086c69b9fec4ddbfb6a1433d28ecca9292b/pkg/util/validation/validation.go#L153
	var maxLength int64 = 63
	openapiSpec["k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"].Schema.SchemaProps.Properties["name"] = spec.Schema{
		SchemaProps: spec.SchemaProps{
			Description: "Name must be unique within a namespace. Is required when creating resources, although some resources may allow a client to request the generation of an appropriate name automatically. Name is primarily intended for creation idempotence and configuration definition. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/identifiers#names",
			Type:        []string{"string"},
			Format:      "",
			Pattern:     `^[a-z]([-a-z0-9]*[a-z0-9])?$`,
			MaxLength:   &maxLength,
		},
	}
}

func (e *BackupScheduleSpec) SetDefaults() {
	if e == nil {
		return
	}
	if e.Resources != nil {
		e.PodTemplate.Spec.Resources = *e.Resources
		e.Resources = nil
	}
}
