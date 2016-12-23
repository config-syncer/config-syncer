package errorhandlers

import (
	"appscode/pkg/system"
	"sync"

	"github.com/appscode/errors"
	"github.com/appscode/errors/h/mailgun"
	_env "github.com/appscode/go/env"
)

const (
	MailFrom     = "postmaster@"
	MailToSuffix = "@appscode.com"
)

var cOnce sync.Once

func Init() {
	cOnce.Do(func() {
		// initialize the error handlers in sequence
		if !_env.FromHost().DevMode() {
			system.Init()
			if h := NewEmailHandler(); h != nil {
				errors.Handlers.Add(h)
			}
		}
	})
}

func NewEmailHandler() *mailgun.EmailHandler {
	if system.Secrets.Mailgun.PublicDomain != "" && system.Secrets.Mailgun.ApiKey != "" {
		emailOptions := mailgun.NewDefaultEmailOptions(
			system.Secrets.Mailgun.PublicDomain,
			system.Secrets.Mailgun.ApiKey,
			"ERROR - ",
			MailFrom+system.Config.Network.PublicUrls.BaseDomain,
			[]string{"oplog" + "+" + "api" + "-" + _env.FromHost().String() + MailToSuffix},
		)

		return mailgun.NewEmailhandler([]string{errors.Internal}, emailOptions)
	}
	return nil
}

func SendMailAndIgnore(err error) {
	if err != nil {
		errors.New().WithCause(err).Internal()
	}
}

func SendMailWithContextAndIgnore(ctx errors.Context, err error) {
	if err != nil {
		errors.New().WithCause(err).WithContext(ctx).Internal()
	}
}
