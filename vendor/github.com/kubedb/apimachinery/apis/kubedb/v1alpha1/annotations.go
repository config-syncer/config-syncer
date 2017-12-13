package v1alpha1

const (
	DatabaseNamePrefix = "kubedb"

	GenericKey = "kubedb.com"

	LabelDatabaseKind = GenericKey + "/kind"
	LabelDatabaseName = GenericKey + "/name"
	LabelJobType      = GenericKey + "/job-type"

	PostgresKey             = ResourceTypePostgres + "." + GenericKey
	PostgresDatabaseVersion = PostgresKey + "/version"

	ElasticsearchKey             = ResourceTypeElasticsearch + "." + GenericKey
	ElasticsearchDatabaseVersion = ElasticsearchKey + "/version"

	MySQLKey             = ResourceTypeMySQL + "." + GenericKey
	MySQLDatabaseVersion = MySQLKey + "/version"

	MongoDBKey             = ResourceTypeMongoDB + "." + GenericKey
	MongoDBDatabaseVersion = MongoDBKey + "/version"

	RedisKey             = ResourceTypeRedis + "." + GenericKey
	RedisDatabaseVersion = RedisKey + "/version"

	MemcachedKey             = ResourceTypeMemcached + "." + GenericKey
	MemcachedDatabaseVersion = MemcachedKey + "/version"

	SnapshotKey         = ResourceTypeSnapshot + "." + GenericKey
	LabelSnapshotStatus = SnapshotKey + "/status"

	PostgresInitSpec      = PostgresKey + "/init"
	ElasticsearchInitSpec = ElasticsearchKey + "/init"
	MySQLInitSpec         = MySQLKey + "/init"
	MongoDBInitSpec       = MongoDBKey + "/init"
	RedisInitSpec         = RedisKey + "/init"
	MemcachedInitSpec     = MemcachedKey + "/init"

	PostgresIgnore      = PostgresKey + "/ignore"
	ElasticsearchIgnore = ElasticsearchKey + "/ignore"
	MySQLIgnore         = MySQLKey + "/ignore"
	MongoDBIgnore       = MongoDBKey + "/ignore"
	RedisIgnore         = RedisKey + "/ignore"
	MemcachedIgnore     = MemcachedKey + "/ignore"

	AgentCoreosPrometheus        = "coreos-prometheus-operator"
	PrometheusExporterPortNumber = 56790
	PrometheusExporterPortName   = "http"
)
