package unified

import (
	"errors"
	"fmt"
	"os"

	"github.com/appscode/envconfig"
	"github.com/appscode/go-notify/hipchat"
	"github.com/appscode/go-notify/log"
	"github.com/appscode/go-notify/mailgun"
	"github.com/appscode/go-notify/plivo"
	"github.com/appscode/go-notify/slack"
	"github.com/appscode/go-notify/smtp"
	"github.com/appscode/go-notify/twilio"
	"github.com/appscode/go-notify/webhook"
)

const (
	NotifyVia = "NOTIFY_VIA"
)

func Default() (interface{}, error) {
	via, ok := os.LookupEnv(NotifyVia)
	if !ok {
		return nil, errors.New(`"NOTIFY_VIA" is not set.`)
	}
	return DefaultVia(via)
}

func DefaultVia(via string) (interface{}, error) {
	switch via {
	case plivo.UID:
		return plivo.Default()
	case twilio.UID:
		return twilio.Default()
	case smtp.UID:
		return smtp.Default()
	case mailgun.UID:
		return mailgun.Default()
	case hipchat.UID:
		return hipchat.Default()
	case slack.UID:
		return slack.Default()
	case log.UID:
		return log.Default()
	case webhook.UID:
		return webhook.Default()
	}
	return nil, fmt.Errorf("Unknown notifier %s", via)
}

func Load(loader envconfig.LoaderFunc) (interface{}, error) {
	via, ok := loader(NotifyVia)
	if !ok {
		return nil, errors.New(`"NOTIFY_VIA" is not set.`)
	}
	return LoadVia(via, loader)
}

func LoadVia(via string, loader envconfig.LoaderFunc) (interface{}, error) {
	switch via {
	case plivo.UID:
		return plivo.Load(loader)
	case twilio.UID:
		return twilio.Load(loader)
	case smtp.UID:
		return smtp.Load(loader)
	case mailgun.UID:
		return mailgun.Load(loader)
	case hipchat.UID:
		return hipchat.Load(loader)
	case slack.UID:
		return slack.Load(loader)
	case log.UID:
		return log.Load(loader)
	case webhook.UID:
		return webhook.Load(loader)
	}
	return nil, fmt.Errorf("Unknown notifier %s", via)
}
