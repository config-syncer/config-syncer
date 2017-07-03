package smtp

import (
	"strconv"
	"strings"

	"github.com/appscode/go-notify/smtp"
	"github.com/appscode/kubed/pkg/notifier"
	"github.com/appscode/kubed/pkg/notifier/extpoints"
)

type biblio struct {
	opts smtp.Options
}

func init() {
	extpoints.Drivers.Register(new(biblio), smtp.UID)
}

func (b *biblio) Notify(body string) error {
	return smtp.New(b.opts).
		WithSubject("Notification").
		WithBody(body).
		SendHtml()
}

func (b *biblio) SetOptions(opts map[string]string) error {
	reqKeys := []string{
		"smtp_host",
		"smtp_port",
		"smtp_username",
		"smtp_password",
		"smtp_from",
		"cluster_admin_email",
	}
	if err := notifier.EnsureRequiredKeys(opts, reqKeys); err != nil {
		return err
	}
	port, err := strconv.ParseInt(opts["smtp_port"], 10, 64)
	if err != nil {
		return err
	}
	b.opts = smtp.Options{
		Host:     opts["smtp_host"],
		Port:     int(port),
		Username: opts["smtp_username"],
		Password: opts["smtp_password"],
		From:     opts["smtp_from"],
		To:       strings.Split(opts["cluster_admin_email"], ","),
	}
	return nil
}

func (b *biblio) Uid() string {
	return smtp.UID
}
