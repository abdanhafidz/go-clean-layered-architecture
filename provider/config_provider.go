package provider

import "abdanhafidz.com/go-boilerplate/config"

type ConfigProvider interface {
	ProvideEnvConfig() config.EnvConfig
	ProvideDatabaseConfig() config.DatabaseConfig
	ProvideJWTConfig() config.JWTConfig
}

type configProvider struct {
	envConfig      config.EnvConfig
	databaseConfig config.DatabaseConfig
	jwtConfig      config.JWTConfig
}

func NewConfigProvider() ConfigProvider {
	envConfig := config.NewEnvConfig("Asia/Jakarta")
	return &configProvider{
		envConfig: envConfig,
		databaseConfig: config.NewDatabaseConfig(
			envConfig.GetDatabaseHost(),
			envConfig.GetDatabaseUser(),
			envConfig.GetDatabasePassword(),
			envConfig.GetDatabaseName(),
			envConfig.GetDatabasePort(),
		),
		jwtConfig: config.NewJWTConfig(
			envConfig.GetSalt(),
		),
	}
}

func (c *configProvider) ProvideEnvConfig() config.EnvConfig {
	return c.envConfig
}
func (c *configProvider) ProvideDatabaseConfig() config.DatabaseConfig {
	return c.databaseConfig
}

func (c *configProvider) ProvideJWTConfig() config.JWTConfig {
	return c.jwtConfig
}
