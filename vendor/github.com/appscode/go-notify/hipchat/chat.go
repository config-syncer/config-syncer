package hipchat

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/appscode/envconfig"
	"github.com/appscode/go-notify"
	"github.com/tbruyelle/hipchat-go/hipchat"
)

const UID = "hipchat"

type Options struct {
	AuthToken          string        `envconfig:"AUTH_TOKEN" required:"true"`
	To                 []string      `envconfig:"TO"`
	BaseURL            string        `envconfig:"BASE_URL"`
	CACertData         string        `envconfig:"CA_CERT_DATA"`
	InsecureSkipVerify bool          `envconfig:"INSECURE_SKIP_VERIFY"`
	Color              hipchat.Color `envconfig:"COLOR"`
	Notify             bool          `envconfig:"NOTIFY"`
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
		return errors.New("missing to")
	}

	h := hipchat.NewClient(c.opt.AuthToken)
	if c.opt.BaseURL != "" {
		u, err := url.Parse(c.opt.BaseURL)
		if err != nil {
			return err
		}
		h.BaseURL = u
	}
	if c.opt.CACertData != "" {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(c.opt.CACertData))

		tlsConfig := &tls.Config{
			RootCAs: caCertPool,
		}
		tlsConfig.BuildNameToCertificate()

		transport := newTransport()
		transport.TLSClientConfig = tlsConfig
		h.SetHTTPClient(&http.Client{Transport: transport})
	}
	if c.opt.InsecureSkipVerify {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: c.opt.InsecureSkipVerify,
		}
		tlsConfig.BuildNameToCertificate()

		transport := newTransport()
		transport.TLSClientConfig = tlsConfig
		h.SetHTTPClient(&http.Client{Transport: transport})
	}

	for _, room := range c.opt.To {
		_, err := h.Room.Notification(room, &hipchat.NotificationRequest{
			Message: c.body,
			Color:   c.opt.Color,
			Notify:  c.opt.Notify,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Copied from http.DefaultTransport. It establishes network connections as needed
// and caches them for reuse by subsequent calls. It uses HTTP proxies
// as directed by the $HTTP_PROXY and $NO_PROXY (or $http_proxy and
// $no_proxy) environment variables.
func newTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}
