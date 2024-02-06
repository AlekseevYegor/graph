package constant

import "github.com/spf13/viper"

const DBDriverName = "postgres"

func init() {
	viper.AutomaticEnv()

	viper.SetDefault("DB_HOST", "127.0.0.1")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_NAME", "graph")
	viper.SetDefault("DB_USER", "graph_db_user")
	viper.SetDefault("DB_PASSWORD", "graph_db_user")
	viper.SetDefault("DB_SCHEMA", "graph")
	viper.SetDefault("SSL_MODE", false)
}
