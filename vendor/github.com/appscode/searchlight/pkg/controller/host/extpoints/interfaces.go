package extpoints

import "github.com/appscode/searchlight/pkg/controller/types"

type IcingaHostType interface {
	CreateAlert(ctx *types.Context, specificObject string) error
	UpdateAlert(ctx *types.Context) error
	DeleteAlert(ctx *types.Context, specificObject string) error
}
