package provider

import (
	"abdanhafidz.com/go-boilerplate/config"
)

type ConfigProvider interface {
	ProvideDatabaseConfig() config.DatabaseConfig
	ProvideEnvConfig() config.EnvConfig
	ProvideUploadConfig() config.UploadConfig
	ProvideSupabaseConfig() config.SupabaseConfig
	ProvideJWTConfig() config.JWTConfig
	ProvideXenditConfig() config.XenditConfig
}

type configProvider struct {
	databaseConfig config.DatabaseConfig
	envConfig      config.EnvConfig
	uploadConfig   config.UploadConfig
	supabaseConfig config.SupabaseConfig
	jWTConfig      config.JWTConfig
	xenditConfig   config.XenditConfig
}

func NewConfigProvider() ConfigProvider {
	envConfig := config.NewEnvConfig("Asia / Jakarta")
	databaseConfig := config.NewDatabaseConfig(envConfig.GetDatabaseHost(), envConfig.GetDatabaseUser(), envConfig.GetDatabasePassword(), envConfig.GetDatabaseName(), envConfig.GetDatabasePort())
	uploadConfig := config.NewUploadConfig()
	supabaseConfig := config.NewSupabaseConfig(envConfig.GetSupabaseURL(), envConfig.GetSupabaseKey(), envConfig.GetSupabaseBucket())
	jWTConfig := config.NewJWTConfig(envConfig.GetSalt())
	xenditConfig := config.NewXenditConfig(envConfig)
	return &configProvider{
		databaseConfig: databaseConfig,
		envConfig:      envConfig,
		uploadConfig:   uploadConfig,
		supabaseConfig: supabaseConfig,
		jWTConfig:      jWTConfig,
		xenditConfig:   xenditConfig,
	}
}

func (c *configProvider) ProvideDatabaseConfig() config.DatabaseConfig {
	return c.databaseConfig
}

func (c *configProvider) ProvideEnvConfig() config.EnvConfig {
	return c.envConfig
}

func (c *configProvider) ProvideUploadConfig() config.UploadConfig {
	return c.uploadConfig
}

func (c *configProvider) ProvideSupabaseConfig() config.SupabaseConfig {
	return c.supabaseConfig
}

func (c *configProvider) ProvideJWTConfig() config.JWTConfig {
	return c.jWTConfig
}

func (c *configProvider) ProvideXenditConfig() config.XenditConfig {
	return c.xenditConfig
}
