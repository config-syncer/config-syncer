package v1alpha1

type ResourceInfo interface {
	ResourceShortCode() string
	ResourceKind() string
	ResourceSingular() string
	ResourcePlural() string
}
