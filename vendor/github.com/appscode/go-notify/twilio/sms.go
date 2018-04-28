package twilio

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"errors"
	"github.com/appscode/envconfig"
	"github.com/appscode/go-notify"
)

const UID = "twilio"

type Options struct {
	AccountSid string   `envconfig:"ACCOUNT_SID" required:"true"`
	AuthToken  string   `envconfig:"AUTH_TOKEN" required:"true"`
	From       string   `envconfig:"FROM" required:"true"`
	To         []string `envconfig:"TO"`
}

type client struct {
	opt  Options
	body string
}

var _ notify.BySMS = &client{}

func New(opt Options) *client {
	return &client{
		opt: opt,
	}
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

func (c client) From(from string) notify.BySMS {
	c.opt.From = from
	return &c
}

func (c client) WithBody(body string) notify.BySMS {
	c.body = body
	return &c
}

func (c client) To(to string, cc ...string) notify.BySMS {
	c.opt.To = append([]string{to}, cc...)
	return &c
}

func (c *client) Send() error {
	if len(c.opt.To) == 0 {
		return errors.New("Missing to")
	}

	hc := &http.Client{Timeout: time.Second * 10}
	v := url.Values{}
	v.Set("From", c.opt.From)
	v.Set("Body", c.body)
	for _, receiver := range c.opt.To {
		v.Set("To", receiver)
		urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%v/Messages.json", c.opt.AccountSid)
		req, err := http.NewRequest("POST", urlStr, strings.NewReader(v.Encode()))
		if err != nil {
			return err
		}

		req.SetBasicAuth(c.opt.AccountSid, c.opt.AuthToken)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		_, err = hc.Do(req)
		if err != nil {
			return err
		}
	}
	return nil
}
