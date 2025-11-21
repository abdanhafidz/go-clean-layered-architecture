package config

type JWTConfig interface {
	SetSecretKey(key string)
	GetSecretKey() string
}

type jwtConfig struct {
	secretKey string
}

func NewJWTConfig(secretKey string) JWTConfig {
	return &jwtConfig{
		secretKey: secretKey,
	}
}

func (cfg *jwtConfig) SetSecretKey(key string) {
	cfg.secretKey = key
}

func (cfg *jwtConfig) GetSecretKey() string {
	return cfg.secretKey
}
