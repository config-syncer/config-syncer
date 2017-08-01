package smtp

import (
	"crypto/tls"

	"errors"
	"github.com/appscode/envconfig"
	"github.com/appscode/go-notify"
	gomail "gopkg.in/gomail.v2"
)

const UID = "smtp"

type Options struct {
	Host               string   `json:"host" envconfig:"HOST" required:"true" form:"smtp_host"`
	Port               int      `json:"port" envconfig:"PORT" required:"true" form:"smtp_port"`
	InsecureSkipVerify bool     `json:"insecure_skip_verify" envconfig:"INSECURE_SKIP_VERIFY" form:"smtp_insecure_skip_verify"`
	Username           string   `json:"username" envconfig:"USERNAME" required:"true" form:"smtp_username"`
	Password           string   `json:"password" envconfig:"PASSWORD" required:"true" form:"smtp_password"`
	From               string   `json:"from" envconfig:"FROM" required:"true" form:"smtp_from"`
	To                 []string `json:"to" envconfig:"TO" form:"smtp_to"`
}

type client struct {
	opt     Options
	subject string
	body    string
	html    bool
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
	return &c
}

func (c client) WithNoTracking() notify.ByEmail {
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

	mail := gomail.NewMessage()
	mail.SetHeader("From", c.opt.From)
	mail.SetHeader("To", c.opt.To...)
	mail.SetHeader("Subject", c.subject)
	if c.html {
		mail.SetBody("text/html", c.body)
	} else {
		mail.SetBody("text/plain", c.body)
	}

	var d *gomail.Dialer
	if c.opt.Username != "" && c.opt.Password != "" {
		d = gomail.NewDialer(c.opt.Host, c.opt.Port, c.opt.Username, c.opt.Password)
	} else {
		d = &gomail.Dialer{Host: c.opt.Host, Port: c.opt.Port}
	}
	if c.opt.InsecureSkipVerify {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return d.DialAndSend(mail)
}

func (c client) SendHtml() error {
	c.html = true
	return c.Send()
}
