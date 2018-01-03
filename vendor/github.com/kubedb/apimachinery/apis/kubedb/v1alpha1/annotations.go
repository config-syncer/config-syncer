package v1alpha1

const (
	DatabaseNamePrefix = "kubedb"

	GenericKey = "kubedb.com"

	LabelDatabaseKind = GenericKey + "/kind"
	LabelDatabaseName = GenericKey + "/name"
	LabelJobType      = GenericKey + "/job-type"

	PostgresKey         = ResourceTypePostgres + "." + GenericKey
	ElasticsearchKey    = ResourceTypeElasticsearch + "." + GenericKey
	MySQLKey            = ResourceTypeMySQL + "." + GenericKey
	MongoDBKey          = ResourceTypeMongoDB + "." + GenericKey
	RedisKey            = ResourceTypeRedis + "." + GenericKey
	MemcachedKey        = ResourceTypeMemcached + "." + GenericKey
	SnapshotKey         = ResourceTypeSnapshot + "." + GenericKey
	LabelSnapshotStatus = SnapshotKey + "/status"

	GenericInitSpec = GenericKey + "/init"

	AgentCoreosPrometheus        = "coreos-prometheus-operator"
	PrometheusExporterPortNumber = 56790
	PrometheusExporterPortName   = "http"
)
