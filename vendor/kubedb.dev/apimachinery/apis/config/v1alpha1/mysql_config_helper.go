package v1alpha1

const DefaultMySQLDatabasePlugin = "mysql-database-plugin"

func (m *MySQLConfiguration) SetDefaults() {
	if m == nil {
		return
	}

	if m.PluginName == "" {
		m.PluginName = DefaultMySQLDatabasePlugin
	}
}
