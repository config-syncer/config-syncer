package pushover

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/appscode/envconfig"
	"github.com/appscode/go-notify"
)

const (
	PushoverApiURL = "https://api.pushover.net/1/messages.json"
	UID            = "pushover"
)

// Options allows full configuration of the message sent to the Pushover API
type Options struct {
	Token string `envconfig:"TOKEN" required:"true"`
	// User may be either a user key or a group key.
	User    string `envconfig:"USER_KEY"`
	Message string `envconfig:"MESSAGE"`

	// Optional params
	Device    []string `envconfig:"DEVICE"`
	Title     string   `envconfig:"TITLE"`
	URL       string   `envconfig:"URL"`
	URLTitle  string   `envconfig:"URL_TITLE"`
	Priority  string   `envconfig:"PRIORITY"`
	Timestamp string   `envconfig:"TIMESTAMP"`
	Sound     string   `envconfig:"SOUND"`
}

type client struct {
	opt Options
}

var _ notify.ByPush = &client{}

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

func (c client) WithBody(body string) notify.ByPush {
	c.opt.Message = body
	return &c
}

func (c client) To(to ...string) notify.ByPush {
	c.opt.Device = append([]string{}, to...)
	return &c
}

func (c *client) Send() error {
	if c.opt.Token == "" {
		return errors.New("Missing token")
	}
	if c.opt.User =="" {
		return errors.New("Missing user")
	}
	if c.opt.Message == "" {
		return errors.New("Missing message")
	}

	msg := makeFormPayload(c.opt)
	buf := bytes.NewBufferString(msg.Encode())

	hc := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, PushoverApiURL, buf)

	resp, err := hc.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	bodyBuf := bytes.NewBuffer([]byte{})
	if _, err := bodyBuf.ReadFrom(resp.Body); err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("%s: %s", resp.Status, bodyBuf.String())
	}

	return nil
}

func makeFormPayload(opt Options) url.Values {
	data := url.Values{}

	if opt.Token != "" {
		data.Set("token", opt.Token)
	}

	if opt.User != "" {
		data.Set("user", opt.User)
	}

	if opt.Message != "" {
		data.Set("message", opt.Message)
	}

	if len(opt.Device) > 0 {
		data.Set("device", strings.Join(opt.Device, ","))
	}

	if opt.Title != "" {
		data.Set("title", opt.Title)
	}

	if opt.URL != "" {
		data.Set("url", opt.URL)
	}

	if opt.URLTitle != "" {
		data.Set("url_title", opt.URLTitle)
	}

	if opt.Priority != "" {
		data.Set("priority", opt.Priority)
	}

	if opt.Timestamp != "" {
		data.Set("timestamp", opt.Timestamp)
	}

	if opt.Sound != "" {
		data.Set("sound", opt.Sound)
	}

	return data
}
