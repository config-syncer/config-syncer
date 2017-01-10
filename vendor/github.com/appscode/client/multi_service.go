package client

import (
	artifactory "github.com/appscode/api/artifactory/v1beta1"
	auth "github.com/appscode/api/auth/v1beta1"
	backup "github.com/appscode/api/backup/v1beta1"
	ca "github.com/appscode/api/certificate/v1beta1"
	db "github.com/appscode/api/db/v1beta1"
	glusterfs "github.com/appscode/api/glusterfs/v1beta1"
	kubernetesV1beta1 "github.com/appscode/api/kubernetes/v1beta1"
	kubernetesV1beta2 "github.com/appscode/api/kubernetes/v1beta2"
	namespace "github.com/appscode/api/namespace/v1beta1"
	pv "github.com/appscode/api/pv/v1beta1"
	"google.golang.org/grpc"
)

// multi client services are grouped by there main client. the api service
// clients are wrapped around with sub-service.
type multiClientInterface interface {
	Artifactory() *artifactoryService
	Authentication() *authenticationService
	Backup() *backupService
	Namespace() *nsService
	CA() *caService
	DB() *dbService
	GlusterFS() *glusterFSService
	Kubernetes() *versionedKubernetesService
	PV() *pvService
}

type multiClientServices struct {
	artifactoryClient         *artifactoryService
	authenticationClient      *authenticationService
	backupClient              *backupService
	nsClient                  *nsService
	caClient                  *caService
	glusterfsClient           *glusterFSService
	versionedKubernetesClient *versionedKubernetesService
	pvClient                  *pvService
	dbClient                  *dbService
}

func newMultiClientService(conn *grpc.ClientConn) multiClientInterface {
	return &multiClientServices{
		artifactoryClient: &artifactoryService{
			artifactClient: artifactory.NewArtifactsClient(conn),
			versionClient:  artifactory.NewVersionsClient(conn),
		},
		authenticationClient: &authenticationService{
			authenticationClient: auth.NewAuthenticationClient(conn),
			conduitClient:        auth.NewConduitClient(conn),
			projectClient:        auth.NewProjectsClient(conn),
		},
		backupClient: &backupService{
			backupServerClient: backup.NewServersClient(conn),
			backupClientClient: backup.NewClientsClient(conn),
		},
		nsClient: &nsService{
			teamClient:    namespace.NewTeamsClient(conn),
			billingClient: namespace.NewBillingClient(conn),
		},
		caClient: &caService{
			certificateClient: ca.NewCertificatesClient(conn),
		},
		glusterfsClient: &glusterFSService{
			clusterClient: glusterfs.NewClustersClient(conn),
			volumeClient:  glusterfs.NewVolumesClient(conn),
		},
		versionedKubernetesClient: &versionedKubernetesService{
			v1beta1Service: &kubernetesV1beta1Service{
				clientsClient:      kubernetesV1beta1.NewClientsClient(conn),
				clusterClient:      kubernetesV1beta1.NewClustersClient(conn),
				eventsClient:       kubernetesV1beta1.NewEventsClient(conn),
				incidentClient:     kubernetesV1beta1.NewIncidentsClient(conn),
				loadBalancerClient: kubernetesV1beta1.NewLoadBalancersClient(conn),
				metdataClient:      kubernetesV1beta1.NewMetadataClient(conn),
			},
			v1beta2Service: &kubernetesV1beta2Service{
				clientsClient: kubernetesV1beta2.NewClientsClient(conn),
				diskClient:    kubernetesV1beta2.NewDisksClient(conn),
			},
		},
		pvClient: &pvService{
			diskClient: pv.NewDisksClient(conn),
			pvClient:   pv.NewPersistentVolumesClient(conn),
			pvcClient:  pv.NewPersistentVolumeClaimsClient(conn),
		},
		dbClient: &dbService{
			database: db.NewDatabasesClient(conn),
			snapshot: db.NewSnapshotsClient(conn),
		},
	}
}

func (s *multiClientServices) Artifactory() *artifactoryService {
	return s.artifactoryClient
}

func (s *multiClientServices) Authentication() *authenticationService {
	return s.authenticationClient
}

func (s *multiClientServices) Backup() *backupService {
	return s.backupClient
}

func (s *multiClientServices) Namespace() *nsService {
	return s.nsClient
}

func (s *multiClientServices) CA() *caService {
	return s.caClient
}

func (s *multiClientServices) GlusterFS() *glusterFSService {
	return s.glusterfsClient
}

func (s *multiClientServices) Kubernetes() *versionedKubernetesService {
	return s.versionedKubernetesClient
}

func (s *multiClientServices) PV() *pvService {
	return s.pvClient
}

func (s *multiClientServices) DB() *dbService {
	return s.dbClient
}

// original service clients that needs to exposed under grouped wrapper
// services.
type artifactoryService struct {
	artifactClient artifactory.ArtifactsClient
	versionClient  artifactory.VersionsClient
}

func (a *artifactoryService) Artifacts() artifactory.ArtifactsClient {
	return a.artifactClient
}

func (a *artifactoryService) Versions() artifactory.VersionsClient {
	return a.versionClient
}

type authenticationService struct {
	authenticationClient auth.AuthenticationClient
	conduitClient        auth.ConduitClient
	projectClient        auth.ProjectsClient
}

func (a *authenticationService) Authentication() auth.AuthenticationClient {
	return a.authenticationClient
}

func (a *authenticationService) Conduit() auth.ConduitClient {
	return a.conduitClient
}

func (a *authenticationService) Project() auth.ProjectsClient {
	return a.projectClient
}

type backupService struct {
	backupClientClient backup.ClientsClient
	backupServerClient backup.ServersClient
}

func (b *backupService) Server() backup.ServersClient {
	return b.backupServerClient
}

func (b *backupService) Client() backup.ClientsClient {
	return b.backupClientClient
}

type nsService struct {
	teamClient    namespace.TeamsClient
	billingClient namespace.BillingClient
}

func (b *nsService) Team() namespace.TeamsClient {
	return b.teamClient
}

func (b *nsService) Billing() namespace.BillingClient {
	return b.billingClient
}

type caService struct {
	certificateClient ca.CertificatesClient
}

func (c *caService) CertificatesClient() ca.CertificatesClient {
	return c.certificateClient
}

type glusterFSService struct {
	clusterClient glusterfs.ClustersClient
	volumeClient  glusterfs.VolumesClient
}

func (g *glusterFSService) Cluster() glusterfs.ClustersClient {
	return g.clusterClient
}

func (g *glusterFSService) Volume() glusterfs.VolumesClient {
	return g.volumeClient
}

type versionedKubernetesService struct {
	v1beta1Service *kubernetesV1beta1Service
	v1beta2Service *kubernetesV1beta2Service
}

func (v *versionedKubernetesService) V1beta1() *kubernetesV1beta1Service {
	return v.v1beta1Service
}

func (v *versionedKubernetesService) V1beta2() *kubernetesV1beta2Service {
	return v.v1beta2Service
}

type kubernetesV1beta1Service struct {
	clientsClient      kubernetesV1beta1.ClientsClient
	clusterClient      kubernetesV1beta1.ClustersClient
	eventsClient       kubernetesV1beta1.EventsClient
	incidentClient     kubernetesV1beta1.IncidentsClient
	loadBalancerClient kubernetesV1beta1.LoadBalancersClient
	metdataClient      kubernetesV1beta1.MetadataClient
}

func (k *kubernetesV1beta1Service) Client() kubernetesV1beta1.ClientsClient {
	return k.clientsClient
}

func (k *kubernetesV1beta1Service) Cluster() kubernetesV1beta1.ClustersClient {
	return k.clusterClient
}

func (k *kubernetesV1beta1Service) Event() kubernetesV1beta1.EventsClient {
	return k.eventsClient
}

func (a *kubernetesV1beta1Service) Incident() kubernetesV1beta1.IncidentsClient {
	return a.incidentClient
}

func (a *kubernetesV1beta1Service) LoadBalancer() kubernetesV1beta1.LoadBalancersClient {
	return a.loadBalancerClient
}

func (k *kubernetesV1beta1Service) Metadata() kubernetesV1beta1.MetadataClient {
	return k.metdataClient
}

type kubernetesV1beta2Service struct {
	clientsClient kubernetesV1beta2.ClientsClient
	diskClient    kubernetesV1beta2.DisksClient
}

func (k *kubernetesV1beta2Service) Client() kubernetesV1beta2.ClientsClient {
	return k.clientsClient
}

func (k *kubernetesV1beta2Service) Disk() kubernetesV1beta2.DisksClient {
	return k.diskClient
}

type pvService struct {
	diskClient pv.DisksClient
	pvClient   pv.PersistentVolumesClient
	pvcClient  pv.PersistentVolumeClaimsClient
}

func (p *pvService) Disk() pv.DisksClient {
	return p.diskClient
}

func (p *pvService) PersistentVolume() pv.PersistentVolumesClient {
	return p.pvClient
}

func (p *pvService) PersistentVolumeClaim() pv.PersistentVolumeClaimsClient {
	return p.pvcClient
}

type dbService struct {
	database db.DatabasesClient
	snapshot db.SnapshotsClient
}

func (p *dbService) Database() db.DatabasesClient {
	return p.database
}

func (p *dbService) Snapshot() db.SnapshotsClient {
	return p.snapshot
}
