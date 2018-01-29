package v1alpha1

const (
	DatabaseNamePrefix = "kubedb"

	GenericKey = "kubedb.com"

	LabelDatabaseKind = GenericKey + "/kind"
	LabelDatabaseName = GenericKey + "/name"

	PostgresKey         = ResourceTypePostgres + "." + GenericKey
	ElasticsearchKey    = ResourceTypeElasticsearch + "." + GenericKey
	MySQLKey            = ResourceTypeMySQL + "." + GenericKey
	MongoDBKey          = ResourceTypeMongoDB + "." + GenericKey
	RedisKey            = ResourceTypeRedis + "." + GenericKey
	MemcachedKey        = ResourceTypeMemcached + "." + GenericKey
	SnapshotKey         = ResourceTypeSnapshot + "." + GenericKey
	LabelSnapshotStatus = SnapshotKey + "/status"

	AnnotationInitialized = GenericKey + "/initialized"
	AnnotationJobType     = GenericKey + "/job-type"

	PrometheusExporterPortNumber = 56790
	PrometheusExporterPortName   = "prom-http"

	JobTypeBackup  = "backup"
	JobTypeRestore = "restore"
)
