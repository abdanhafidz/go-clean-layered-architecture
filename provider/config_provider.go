package provider

import "abdanhafidz.com/go-clean-layered-architecture/config"

type ConfigProvider interface {
	ProvideJWTConfig() config.JWTConfig
	ProvideEnvConfig() config.EnvConfig
	ProvideDatabaseConfig() config.DatabaseConfig
}

type configProvider struct {
	jWTConfig      config.JWTConfig
	envConfig      config.EnvConfig
	databaseConfig config.DatabaseConfig
}

func NewConfigProvider() ConfigProvider {
	envConfig := config.NewEnvConfig("Asia/Jakarta")
	jWTConfig := config.NewJWTConfig(envConfig.GetSalt())
	databaseConfig := config.NewDatabaseConfig(
		envConfig.GetDatabaseHost(),
		envConfig.GetDatabaseUser(),
		envConfig.GetDatabasePassword(),
		envConfig.GetDatabaseName(),
		envConfig.GetDatabasePort())
	return &configProvider{
		jWTConfig:      jWTConfig,
		envConfig:      envConfig,
		databaseConfig: databaseConfig,
	}
}

func (c *configProvider) ProvideJWTConfig() config.JWTConfig {
	return c.jWTConfig
}

func (c *configProvider) ProvideEnvConfig() config.EnvConfig {
	return c.envConfig
}

func (c *configProvider) ProvideDatabaseConfig() config.DatabaseConfig {
	return c.databaseConfig
}
