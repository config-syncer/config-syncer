package client

import (
	credential "github.com/appscode/api/credential/v1beta1"
	"github.com/appscode/api/health"
	mailinglist "github.com/appscode/api/mailinglist/v1beta1"
	namespace "github.com/appscode/api/namespace/v1beta1"
	operation "github.com/appscode/api/operation/v1beta1"
	ssh "github.com/appscode/api/ssh/v1beta1"
	"google.golang.org/grpc"
)

// single client service in api. returned directly the api client.
type loneClientInterface interface {
	CloudCredential() credential.CredentialsClient
	Health() health.HealthClient
	MailingList() mailinglist.MailingListClient
	Team() namespace.TeamsClient
	Operation() operation.OperationsClient
	SSH() ssh.SSHClient
}

type loneClientServices struct {
	credClient        credential.CredentialsClient
	healthClient      health.HealthClient
	mailingListClient mailinglist.MailingListClient
	teamClient        namespace.TeamsClient
	operationClient   operation.OperationsClient
	sshClient         ssh.SSHClient
}

func newLoneClientService(conn *grpc.ClientConn) loneClientInterface {
	return &loneClientServices{
		credClient:        credential.NewCredentialsClient(conn),
		healthClient:      health.NewHealthClient(conn),
		mailingListClient: mailinglist.NewMailingListClient(conn),
		teamClient:        namespace.NewTeamsClient(conn),
		operationClient:   operation.NewOperationsClient(conn),
		sshClient:         ssh.NewSSHClient(conn),
	}
}

func (s *loneClientServices) CloudCredential() credential.CredentialsClient {
	return s.credClient
}

func (s *loneClientServices) Health() health.HealthClient {
	return s.healthClient
}

func (s *loneClientServices) Team() namespace.TeamsClient {
	return s.teamClient
}

func (s *loneClientServices) SSH() ssh.SSHClient {
	return s.sshClient
}

func (s *loneClientServices) MailingList() mailinglist.MailingListClient {
	return s.mailingListClient
}

func (s *loneClientServices) Operation() operation.OperationsClient {
	return s.operationClient
}
