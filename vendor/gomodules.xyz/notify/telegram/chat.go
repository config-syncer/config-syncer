package telegram

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/golang/glog"
	"gomodules.xyz/envconfig"
	"gomodules.xyz/notify"
)

const UID = "telegram"

type Options struct {
	Token   string   `envconfig:"TOKEN" required:"true"`
	Channel []string `envconfig:"CHANNEL"`
}

type client struct {
	opt       Options
	body      string
	parseMode string
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

func (c *client) WithChannel(name string) {
	c.opt.Channel = []string{name}
}

func (c client) UID() string {
	return UID
}

func (c client) WithBody(body string) notify.ByChat {
	c.body = body
	return &c
}

func (c *client) WithParseMode() {
	c.parseMode = "HTML"
}

func (c client) To(to string, cc ...string) notify.ByChat {
	c.opt.Channel = append([]string{to}, cc...)
	return &c
}

func (c *client) Send() error {
	if len(c.opt.Channel) == 0 {
		return errors.New("missing to")
	}

	type ErrorResponse struct {
		Ok          bool   `json:"ok"`
		ErrorCode   int    `json:"error_code"`
		Description string `json:"description"`
	}

	u := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.opt.Token)

	for _, channel := range c.opt.Channel {
		data := url.Values{}

		if c.parseMode != "" {
			data.Set("parse_mode", c.parseMode)
		}
		data.Set("text", c.body)
		data.Set("chat_id", channel)

		resp, err := http.PostForm(u, data)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			var r ErrorResponse
			err := json.NewDecoder(resp.Body).Decode(&r)
			if err != nil {
				glog.Warningf("failed to send message to channel %s. Reason: %s", channel, err)
			} else if !r.Ok {
				glog.Warningf("failed to send message to channel %s. Reason: %d - %s", channel, r.ErrorCode, r.Description)
			}
		}
		resp.Body.Close()
	}
	return nil
}
