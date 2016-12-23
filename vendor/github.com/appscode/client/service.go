package client

import "google.golang.org/grpc"

type serviceInterface interface {
	loneClientInterface
	multiClientInterface
	Connection() *grpc.ClientConn
}

type services struct {
	loneClientInterface
	multiClientInterface

	connection *grpc.ClientConn
}

func (s *services) loneClientServices() loneClientInterface {
	return s.loneClientInterface
}

func (s *services) multiClientServices() multiClientInterface {
	return s.multiClientInterface
}

func (s *services) Connection() *grpc.ClientConn {
	return s.connection
}

func newServices(conn *grpc.ClientConn) serviceInterface {
	return &services{
		loneClientInterface:  newLoneClientService(conn),
		multiClientInterface: newMultiClientService(conn),
		connection:           conn,
	}
}
