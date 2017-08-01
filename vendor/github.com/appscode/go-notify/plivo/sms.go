package plivo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/appscode/envconfig"
	"github.com/appscode/go-notify"
)

const (
	UID         = "plivo"
	urlTemplate = "https://api.plivo.com/v1/Account/%v/Message/"
)

type Options struct {
	AuthID    string   `envconfig:"AUTH_ID" required:"true"`
	AuthToken string   `envconfig:"AUTH_TOKEN" required:"true"`
	From      string   `envconfig:"FROM" required:"true"`
	To        []string `envconfig:"TO"`
}

type client struct {
	opt  Options
	body string
}

var _ notify.BySMS = &client{}

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

	url := fmt.Sprintf(urlTemplate, c.opt.AuthID)
	params := struct {
		Src  string `json:"src,omitempty"`
		Dst  string `json:"dst,omitempty"`
		Text string `json:"text,omitempty"`
	}{
		c.opt.From,
		"",
		c.body,
	}
	for _, dst := range c.opt.To {
		params.Dst = dst
		buf := new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(params); err != nil {
			return err
		}

		req, err := http.NewRequest("POST", url, buf)
		if err != nil {
			return err
		}

		req.SetBasicAuth(c.opt.AuthID, c.opt.AuthToken)
		req.Header.Add("Content-Type", "application/json")

		resp, err := hc.Do(req)
		if err != nil {
			return err
		}

		respBody := struct {
			Error string `json:"error"`
		}{}
		if err = json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
			return err
		}
		if respBody.Error != "" {
			return errors.New(respBody.Error)
		}

		resp.Body.Close()
	}
	return nil
}
