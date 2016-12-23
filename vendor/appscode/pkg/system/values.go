package system

const (
	CertTrusted = iota - 1
	CertRoot
	CertNSRoot
	CertIntermediate
	CertLeaf

	RoleJenkinsMaster = "jenkins-master"
	RoleJenkinsAgent  = "jenkins-agent"
	RoleJenkinsShared = "jenkins-shared"

	RoleKubernetesMaster = "kubernetes-master"
	RoleKubernetesPool   = "kubernetes-pool"

	CIBotUser      = "ci-bot"
	ClusterBotUser = "k8s-bot"

	AllowHttpAuth                  = "diffusion.allow-http-auth"
	AppsCodePrivateAPIHttpEndpoint = "appscode.private-api-http-endpoint"
	AppsCodePublicAPIHttpEndpoint  = "appscode.public-api-http-endpoint"
	BaseUri                        = "phabricator.base-uri"
	BraintreeCustomerID            = "braintree.customer-id"
	BraintreeMerchantID            = "braintree.merchant-id"
	BraintreePrivateKey            = "braintree.private-key"
	BraintreePublicKey             = "braintree.public-key"
	CIBucket                       = "ci.data-bucket-name"
	CIDefaultBot                   = "ci.default-bot"
	CIMasterEndpoint               = "ci.master-endpoint"
	CIMasterService                = "ci.master-service"
	CIServiceAccount               = "ci.service-account"
	CSRF_Key                       = "phabricator.csrf-key"
	CSRF_Value                     = "0b7ec0592e0a2829d8b71df2fa269b2c6172eca3"
	DigitalOceanCredential         = "digitalocean.credential"
	DNSCredential                  = "dns.credential"
	ElasticSearchHost              = "search.elastic.host"
	ElasticSearchNamespace         = "search.elastic.namespace"
	MailgunApiKey                  = "mailgun.api-key"
	MailgunPublicDomain            = "mailgun.public-domain"
	MailgunTeamDomain              = "mailgun.domain"
	MetamtaDefaultAddress          = "metamta.default-address"
	MetamtaDomain                  = "metamta.domain"
	MetamtaMailAdapter             = "metamta.mail-adapter"
	MetamtaReplyHandlerDomain      = "metamta.reply-handler-domain"
	NSDeactivationPeriod           = "ns.deactivation-period"
	NSMinRollingPayment            = "ns.min-rolling-payment"
	NSRollingPayment               = "ns.rolling-payment"
	PhabricatorBucket              = "phabricator.data-bucket-name"
	PhabricatorS3Bucket            = "storage.s3.bucket"
	PhabricatorLocalDiskStorage    = "storage.local-disk.path"
	PygmentsEnabled                = "pygments.enabled"
	RepositoryPath                 = "repository.default-local-path"
	SecurityAlternateFileDomain    = "security.alternate-file-domain"
	ShortUri                       = "phurl.short-uri"
	TwilioAccountSid               = "twilio.account-sid"
	TwilioAuthToken                = "twilio.auth-token"
	TwilioPhoneNumber              = "twilio.phone-number"
	VCSHost                        = "diffusion.ssh-host"
	VCSUser                        = "diffusion.ssh-user"

	GCS   = "gcs"
	S3    = "s3"
	Local = "local"
)
