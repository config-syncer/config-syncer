package plivo

import (
	"strings"

	"github.com/appscode/go-notify/plivo"
	"github.com/appscode/kubed/pkg/notifier"
	"github.com/appscode/kubed/pkg/notifier/extpoints"
)

type biblio struct {
	opts plivo.Options
}

func init() {
	extpoints.Drivers.Register(new(biblio), plivo.UID)
}

func (b *biblio) Notify(body string) error {
	return plivo.New(b.opts).WithBody(body).Send()
}

func (b *biblio) SetOptions(opts map[string]string) error {
	reqKeys := []string{
		"plivo_auth_id",
		"plivo_auth_token",
		"plivo_from",
		"cluster_admin_phone",
	}
	if err := notifier.EnsureRequiredKeys(opts, reqKeys); err != nil {
		return err
	}
	b.opts = plivo.Options{
		AuthID:    opts["plivo_auth_id"],
		AuthToken: opts["plivo_auth_token"],
		To:        strings.Split(opts["cluster_admin_phone"], ","),
		From:      opts["plivo_from"],
	}
	return nil
}

func (b *biblio) Uid() string {
	return plivo.UID
}
