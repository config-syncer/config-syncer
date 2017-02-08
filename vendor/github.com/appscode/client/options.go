package client

import (
	"net"
	"strings"
	"time"

	"github.com/appscode/client/credential"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	DefaultOption = &ClientOption{
		Endpoint: defaultApiEndpoint,
	}
)

type ClientOption struct {
	Endpoint   string
	timeout    time.Duration
	auth       *credential.Auth
	tlsEnabled bool
	userAgent  string

	tracing bool
}

func NewDefaultOption() *ClientOption {
	return DefaultOption
}

func NewOption(endpoint string) *ClientOption {
	return &ClientOption{
		Endpoint: endpoint,
	}
}

func NewOptionWithAuth(endpoint string, a *credential.Auth) *ClientOption {
	return &ClientOption{
		Endpoint: endpoint,
		auth:     a,
	}
}

func (s *ClientOption) BearerAuth(namespace, token string) *ClientOption {
	s.auth = credential.NewBearerAuth(namespace, token)
	return s
}

func (s *ClientOption) BasicAuth(namespace, username, password string) *ClientOption {
	s.auth = credential.NewBasicAuth(namespace, username, password)
	return s
}

func (s *ClientOption) UserAgent(a string) *ClientOption {
	s.userAgent = a
	return s
}

func (s *ClientOption) Timeout(d time.Duration) *ClientOption {
	s.timeout = d
	return s
}

func (s *ClientOption) Trace() *ClientOption {
	s.tracing = true
	return s
}

func (o *ClientOption) target() string {
	if strings.HasPrefix(o.Endpoint, "http://") {
		o.Endpoint = o.Endpoint[7:]
	} else if strings.HasPrefix(o.Endpoint, "https://") {
		o.Endpoint = o.Endpoint[8:]
		o.tlsEnabled = true
	}
	return o.Endpoint
}

func (o *ClientOption) parse() []grpc.DialOption {
	dialOptions := make([]grpc.DialOption, 0)

	if o.auth != nil {
		cred := o.auth.Credential()
		dialOptions = append(dialOptions, grpc.WithPerRPCCredentials(cred))
	}

	host, _, err := net.SplitHostPort(o.target())
	if o.tlsEnabled && err == nil {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(
			credentials.NewClientTLSFromCert(nil, host),
		))
	} else {
		dialOptions = append(dialOptions, grpc.WithInsecure())
	}

	if o.timeout != time.Duration(0) {
		dialOptions = append(dialOptions, grpc.WithTimeout(o.timeout))
	}

	if o.userAgent != "" {
		dialOptions = append(dialOptions, grpc.WithUserAgent(o.userAgent))
	}

	if o.tracing {
		grpc.EnableTracing = true
	}

	return dialOptions
}
