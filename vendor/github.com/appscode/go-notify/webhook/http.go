package webhook

import (
	"net/http"

	"github.com/appscode/envconfig"
	"github.com/appscode/go-notify"
	"github.com/appscode/go/net/httpclient"
)

const UID = "webhook"

type Options struct {
	URL                string   `envconfig:"URL" required:"true"`
	To                 []string `envconfig:"TO" required:"true"`
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
	to   []string
	body string
}

var _ notify.ByChat = &client{}

func New(opt Options) *client {
	return &client{
		opt: opt,
		to:  opt.To,
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

func (c client) WithBody(body string) notify.ByChat {
	c.body = body
	return &c
}

func (c client) To(to string, cc ...string) notify.ByChat {
	c.to = append([]string{to}, cc...)
	return &c
}

func (c *client) Send() error {
	client := httpclient.Default().
		WithBaseURL(c.opt.URL).
		WithBasicAuth(c.opt.Username, c.opt.Password).
		WithBearerToken(c.opt.Token)
	if c.opt.CACertData != "" {
		if c.opt.ClientCertData != "" && c.opt.ClientKeyData != "" {
			client = client.WithTLSConfig([]byte(c.opt.CACertData), []byte(c.opt.ClientKeyData), []byte(c.opt.ClientKeyData))
		} else {
			client = client.WithTLSConfig([]byte(c.opt.CACertData))
		}
	}
	if c.opt.InsecureSkipVerify {
		client = client.WithInsecureSkipVerify()
	}

	msg := struct {
		To   []string `json:"to,omitempty"`
		Body string   `json:"body,omitempty"`
	}{
		c.opt.To,
		c.body,
	}
	_, err := client.Call(http.MethodPost, "", msg, nil, true)
	return err
}
