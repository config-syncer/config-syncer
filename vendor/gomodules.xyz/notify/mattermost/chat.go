package mattermost

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"gomodules.xyz/envconfig"
	"gomodules.xyz/notify"
)

const UID = "mattermost"

type Options struct {
	Url     string   `envconfig:"URL" required:"true"`
	HookId  string   `envconfig:"HOOK_ID" required:"true"`
	IconUrl string   `envconfig:"ICON_URL"`
	BotName string   `envconfig:"BOT_NAME"`
	Channel []string `envconfig:"CHANNEL"`
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
	c.opt.Channel = append([]string{to}, cc...)
	return &c
}

type message struct {
	Channel  string `json:"channel,omitempty"`
	Username string `json:"username,omitempty"`
	IconUrl  string `json:"icon_url,omitempty"`
	Text     string `json:"text"`
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

	u := fmt.Sprintf("%s/hooks/%s", c.opt.Url, c.opt.HookId)

	for _, channel := range c.opt.Channel {

		m := message{
			Channel:  channel,
			Username: c.opt.BotName,
			IconUrl:  c.opt.IconUrl,
			Text:     c.body,
		}

		msg, err := json.Marshal(m)
		if err != nil {
			return err
		}

		resp, err := http.Post(u, "application/json", bytes.NewBuffer(msg))

		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			var r ErrorResponse
			err := json.NewDecoder(resp.Body).Decode(&r)
			if err == nil && !r.Ok {
				glog.Warningf("failed to send message to channel %s. Reason: %d - %s", channel, r.ErrorCode, r.Description)
			}
		}
		resp.Body.Close()
	}
	return nil
}
