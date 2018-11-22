package v1alpha1

const DefaultPostgresDatabasePlugin = "postgresql-database-plugin"

func (p *PostgresConfiguration) SetDefaults() {
	if p == nil {
		return
	}

	if p.PluginName == "" {
		p.PluginName = DefaultPostgresDatabasePlugin
	}
}
