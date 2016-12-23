package credential

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/credentials"
)

type bearerAuth struct {
	Namespace   string
	BearerToken string
}

func NewBearerAuthCredential(namespace, token string) credentials.PerRPCCredentials {
	return &bearerAuth{
		Namespace:   namespace,
		BearerToken: token,
	}
}

func (b *bearerAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + b.Namespace + ":" + b.BearerToken,
	}, nil
}

func (b *bearerAuth) RequireTransportSecurity() bool {
	return false
}
