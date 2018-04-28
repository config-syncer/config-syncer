package mailgun

import (
	"errors"
	"github.com/appscode/envconfig"
	notify "github.com/appscode/go-notify"
	h2t "github.com/jaytaylor/html2text"
	mailgun "github.com/mailgun/mailgun-go"
)

const UID = "mailgun"

type Options struct {
	Domain          string   `json:"domain" envconfig:"DOMAIN" required:"true" form:"mailgun_domain"`
	ApiKey          string   `json:"api_key" envconfig:"API_KEY" required:"true" form:"mailgun_api_key"`
	PublicApiKey    string   `json:"public_api_key" envconfig:"PUBLIC_API_KEY" form:"mailgun_public_api_key"`
	From            string   `json:"from" envconfig:"FROM" required:"true" form:"mailgun_from"`
	To              []string `json:"to" envconfig:"TO" form:"mailgun_to"`
	DisableTracking bool     `json:"disable_tracking" envconfig:"DISABLE_TRACKING" from:"disable_tracking"`
}

type client struct {
	opt     Options
	subject string
	body    string
	html    bool
	tag     string
}

var _ notify.ByEmail = &client{}

func New(opt Options) *client {
	return &client{opt: opt}
}

func Default() (*client, error) {
	var opt Options
	err := envconfig.Process(UID, &opt)
	if err != nil {
		return nil, err
	}
	return New(opt), nil
}

func Load(loader envconfig.LoaderFunc) (*client, error) {
	var opt Options
	err := envconfig.Load(UID, &opt, loader)
	if err != nil {
		return nil, err
	}
	return New(opt), nil
}

func (c client) UID() string {
	return UID
}

func (c client) From(from string) notify.ByEmail {
	c.opt.From = from
	return &c
}

func (c client) WithSubject(subject string) notify.ByEmail {
	c.subject = subject
	return &c
}

func (c client) WithBody(body string) notify.ByEmail {
	c.body = body
	return &c
}

func (c client) WithTag(tag string) notify.ByEmail {
	c.tag = tag
	return &c
}

func (c client) WithNoTracking() notify.ByEmail {
	c.opt.DisableTracking = true
	return &c
}

func (c client) To(to string, cc ...string) notify.ByEmail {
	c.opt.To = append([]string{to}, cc...)
	return &c
}

func (c *client) Send() error {
	if len(c.opt.To) == 0 {
		return errors.New("Missing to")
	}

	mg := mailgun.NewMailgun(c.opt.Domain, c.opt.ApiKey, c.opt.PublicApiKey)
	text := c.body
	if c.html {
		if t, err := h2t.FromString(c.body); err == nil {
			text = t
		}
	}
	msg := mg.NewMessage(c.opt.From, c.subject, text, c.opt.To...)
	if c.html {
		msg.SetHtml(c.body)
	}
	if c.tag != "" {
		msg.AddTag(c.tag)
	}
	msg.SetTracking(!c.opt.DisableTracking)
	msg.SetTrackingClicks(!c.opt.DisableTracking)
	msg.SetTrackingOpens(!c.opt.DisableTracking)
	_, _, err := mg.Send(msg)
	return err
}

func (c client) SendHtml() error {
	c.html = true
	return c.Send()
}
