package telegram

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/appscode/envconfig"
	"github.com/appscode/go-notify"
	"github.com/golang/glog"
)

const UID = "telegram"

type Options struct {
	Token   string   `envconfig:"TOKEN" required:"true"`
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
		data.Set("chat_id", channel)
		data.Set("text", c.body)
		resp, err := http.PostForm(u, data)
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
