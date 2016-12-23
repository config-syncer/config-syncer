package credential

import (
	"google.golang.org/grpc/credentials"
)

type AuthType string

const (
	Basic  AuthType = "Basic"
	Bearer AuthType = "Bearer"
)

type Auth struct {
	Namespace string
	Username  string
	AuthType  AuthType
	Secret    string
}

func NewBasicAuth(ns, username, secret string) *Auth {
	return &Auth{
		Namespace: ns,
		Username:  username,
		AuthType:  Basic,
		Secret:    secret,
	}
}

func NewBearerAuth(ns, secret string) *Auth {
	return &Auth{
		Namespace: ns,
		AuthType:  Bearer,
		Secret:    secret,
	}
}

func (a *Auth) Credential() credentials.PerRPCCredentials {
	if a.AuthType == Basic {
		return NewBasicAuthCredential(a.Namespace, a.Username, a.Secret)
	}
	return NewBearerAuthCredential(a.Namespace, a.Secret)
}
