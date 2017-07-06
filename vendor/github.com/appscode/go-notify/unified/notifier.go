package unified

import (
	"errors"
	"fmt"
	"os"

	"github.com/appscode/go-notify/hipchat"
	"github.com/appscode/go-notify/mailgun"
	"github.com/appscode/go-notify/plivo"
	"github.com/appscode/go-notify/slack"
	"github.com/appscode/go-notify/smtp"
	"github.com/appscode/go-notify/twilio"
)

func Default() (interface{}, error) {
	via, ok := os.LookupEnv("NOTIFY_VIA")
	if !ok {
		return nil, errors.New(`"NOTIFY_VIA" is not set.`)
	}
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
	}
	return nil, fmt.Errorf("Unknown notifier %s", via)
}

func Load(loader func(string) (string, bool)) (interface{}, error) {
	via, ok := loader("NOTIFY_VIA")
	if !ok {
		return nil, errors.New(`"NOTIFY_VIA" is not set.`)
	}
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
	}
	return nil, fmt.Errorf("Unknown notifier %s", via)
}
