package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvConfig interface {
	GetTCPAddress() string
	GetLogPath() string
	GetHostAddress() string
	GetHostPort() string
	GetEmailVerificationDuration() int
	GetDatabaseHost() string
	GetDatabasePort() string
	GetDatabaseUser() string
	GetDatabasePassword() string
	GetDatabaseName() string
	GetSalt() string
}

type envConfig struct {
	timezone string
}

func NewEnvConfig(timezone string) EnvConfig {
	godotenv.Load()
	os.Setenv("TZ", timezone)
	return &envConfig{
		timezone: timezone,
	}
}

func (e *envConfig) GetTCPAddress() string {
	return os.Getenv("HOST_ADDRESS") + ":" + os.Getenv("HOST_PORT")
}

func (e *envConfig) GetLogPath() string {
	return os.Getenv("LOG_PATH")
}

func (e *envConfig) GetHostAddress() string {
	return os.Getenv("HOST_ADDRESS")
}

func (e *envConfig) GetHostPort() string {
	return os.Getenv("HOST_PORT")
}

func (e *envConfig) GetEmailVerificationDuration() int {
	duration, err := strconv.Atoi(os.Getenv("EMAIL_VERIFICATION_DURATION"))
	if err != nil {
		return 0 // Default value if parsing fails
	}
	return duration
}

func (e *envConfig) GetDatabaseHost() string {
	return os.Getenv("DB_HOST")
}

func (e *envConfig) GetDatabasePort() string {
	return os.Getenv("DB_PORT")
}

func (e *envConfig) GetDatabaseUser() string {
	return os.Getenv("DB_USER")
}

func (e *envConfig) GetDatabasePassword() string {
	return os.Getenv("DB_PASSWORD")
}

func (e *envConfig) GetDatabaseName() string {
	return os.Getenv("DB_NAME")
}

func (e *envConfig) GetSalt() string {
	salt := os.Getenv("SALT")
	if salt == "" {
		return "Def4u|7" // Default salt value
	}
	return salt
}
