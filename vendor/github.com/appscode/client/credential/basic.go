package credential

import (
	"github.com/appscode/client/util"
	"golang.org/x/net/context"
	"google.golang.org/grpc/credentials"
)

type basicAuth struct {
	Namespace string
	UserName  string
	Password  string
}

func NewBasicAuthCredential(namespace, username, password string) credentials.PerRPCCredentials {
	return &basicAuth{
		Namespace: namespace,
		UserName:  username,
		Password:  password,
	}
}

func (b *basicAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Basic " + util.BasicEncode(b.Namespace, b.UserName, b.Password),
	}, nil
}

func (b *basicAuth) RequireTransportSecurity() bool {
	return false
}
