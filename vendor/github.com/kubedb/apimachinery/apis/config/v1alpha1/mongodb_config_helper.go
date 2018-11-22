package v1alpha1

const DefaultMongoDBDatabasePlugin = "mongodb-database-plugin"

func (m *MongoDBConfiguration) SetDefaults() {
	if m == nil {
		return
	}

	if m.PluginName == "" {
		m.PluginName = DefaultMongoDBDatabasePlugin
	}
}
