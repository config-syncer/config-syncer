package hipchat

import (
	"errors"
	"github.com/appscode/envconfig"
	"github.com/appscode/go-notify"
	"github.com/tbruyelle/hipchat-go/hipchat"
)

const UID = "hipchat"

type Options struct {
	AuthToken string   `envconfig:"AUTH_TOKEN" required:"true"`
	To        []string `envconfig:"TO"`
}

type client struct {
	opt  Options
	body string
}

var _ notify.ByChat = &client{}

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

func (c client) WithBody(body string) notify.ByChat {
	c.body = body
	return &c
}

func (c client) To(to string, cc ...string) notify.ByChat {
	c.opt.To = append([]string{to}, cc...)
	return &c
}

func (c *client) Send() error {
	if len(c.opt.To) == 0 {
		return errors.New("Missing to")
	}

	h := hipchat.NewClient(c.opt.AuthToken)
	for _, room := range c.opt.To {
		_, err := h.Room.Notification(room, &hipchat.NotificationRequest{Message: c.body})
		if err != nil {
			return err
		}
	}
	return nil
}
