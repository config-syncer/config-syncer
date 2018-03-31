package v1alpha1

const (
	DatabaseNamePrefix = "kubedb"

	GenericKey = "kubedb.com"

	LabelDatabaseKind = GenericKey + "/kind"
	LabelDatabaseName = GenericKey + "/name"

	PostgresKey         = ResourcePluralPostgres + "." + GenericKey
	ElasticsearchKey    = ResourcePluralElasticsearch + "." + GenericKey
	MySQLKey            = ResourcePluralMySQL + "." + GenericKey
	MongoDBKey          = ResourcePluralMongoDB + "." + GenericKey
	RedisKey            = ResourcePluralRedis + "." + GenericKey
	MemcachedKey        = ResourcePluralMemcached + "." + GenericKey
	SnapshotKey         = ResourcePluralSnapshot + "." + GenericKey
	LabelSnapshotStatus = SnapshotKey + "/status"

	AnnotationInitialized = GenericKey + "/initialized"
	AnnotationJobType     = GenericKey + "/job-type"

	PrometheusExporterPortNumber = 56790
	PrometheusExporterPortName   = "prom-http"

	JobTypeBackup  = "backup"
	JobTypeRestore = "restore"
)
