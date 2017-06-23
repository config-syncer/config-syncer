package plivo

import (
	"errors"
	"strings"

	"github.com/appscode/go-notify/plivo"
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
	if _, found := opts["plivo_auth_id"]; !found {
		return errors.New("plivo_auth_id not found")
	}
	if _, found := opts["plivo_auth_token"]; !found {
		return errors.New("plivo_auth_token not found")
	}
	if _, found := opts["plivo_from"]; !found {
		return errors.New("plivo_from not found")
	}
	if _, found := opts["cluster_admin_phone"]; !found {
		return errors.New("cluster_admin_phone  not found")
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
