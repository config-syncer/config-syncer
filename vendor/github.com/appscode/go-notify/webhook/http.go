package webhook

import (
	"net/http"

	"errors"
	"github.com/appscode/envconfig"
	"github.com/appscode/go-notify"
	"github.com/appscode/go/net/httpclient"
)

const UID = "webhook"

type Options struct {
	URL                string   `envconfig:"URL" required:"true"`
	To                 []string `envconfig:"TO"`
	Username           string   `envconfig:"USERNAME"`
	Password           string   `envconfig:"PASSWORD"`
	Token              string   `envconfig:"TOKEN"`
	CACertData         string   `envconfig:"CA_CERT_DATA"`
	ClientCertData     string   `envconfig:"CLIENT_CERT_DATA"`
	ClientKeyData      string   `envconfig:"CLIENT_KEY_DATA"`
	InsecureSkipVerify bool     `envconfig:"INSECURE_SKIP_VERIFY"`
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

	hc := httpclient.Default().
		WithBaseURL(c.opt.URL).
		WithBasicAuth(c.opt.Username, c.opt.Password).
		WithBearerToken(c.opt.Token)
	if c.opt.CACertData != "" {
		if c.opt.ClientCertData != "" && c.opt.ClientKeyData != "" {
			hc = hc.WithTLSConfig([]byte(c.opt.CACertData), []byte(c.opt.ClientKeyData), []byte(c.opt.ClientKeyData))
		} else {
			hc = hc.WithTLSConfig([]byte(c.opt.CACertData))
		}
	}
	if c.opt.InsecureSkipVerify {
		hc = hc.WithInsecureSkipVerify()
	}

	msg := struct {
		To   []string `json:"to,omitempty"`
		Body string   `json:"body,omitempty"`
	}{
		c.opt.To,
		c.body,
	}
	_, err := hc.Call(http.MethodPost, "", msg, nil, true)
	return err
}
