package client

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	defaultApiEndpoint = "127.0.0.1:50077"
)

type Client struct {
	connection *grpc.ClientConn

	// serviceInterface implements all the underlying services
	// for the api server. single service clients are direct
	// implementations of grpc service. multi service clients are
	// wrapper around services to group the underlying clients
	// with parent service.
	serviceInterface
}

// this constructs and returns an new api client to access the
// underlying services. ClientOption defines the grpc dial
// options and authentication mechanisms for api server.
func New(option *ClientOption) (*Client, error) {
	c := &Client{}
	var err error

	dialOpts := option.parse()
	c.connection, err = grpc.Dial(option.target(), dialOpts...)
	if err != nil {
		return nil, err
	}
	c.serviceInterface = newServices(c.connection)
	return c, nil
}

// closes the grpc connection with the api server.
func (c *Client) Close() error {
	return c.connection.Close()
}

func (c *Client) Context() context.Context {
	return context.Background()
}
