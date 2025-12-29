package config

import (
	"sync"

	xendit "github.com/xendit/xendit-go/v7"
)

var (
	xenditOnce sync.Once
)

type XenditConfig interface {
	GetClient() *xendit.APIClient
}

type xenditConfig struct {
	envConfig EnvConfig
	client    *xendit.APIClient
}

func NewXenditConfig(envConfig EnvConfig) XenditConfig {
	return &xenditConfig{
		envConfig: envConfig,
		client:    xendit.NewClient(envConfig.GetXenditAPIKey()),
	}
}

func (c *xenditConfig) GetClient() *xendit.APIClient {
	return c.client
}
