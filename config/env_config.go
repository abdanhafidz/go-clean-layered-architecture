package config

import (
	"os"
	"strconv"
	"strings"

	"abdanhafidz.com/go-boilerplate/utils"
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
	GetSupabaseURL() string
	GetSupabaseKey() string
	GetSupabaseBucket() string
	GetXenditAPIKey() string
	GetXenditCallbackToken() string
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
	return utils.GetEnv("HOST_ADDRESS") + ":" + utils.GetEnv("HOST_PORT")
}

func (e *envConfig) GetLogPath() string {
	return utils.GetEnv("LOG_PATH")
}

func (e *envConfig) GetHostAddress() string {
	return utils.GetEnv("HOST_ADDRESS")
}

func (e *envConfig) GetHostPort() string {
	return utils.GetEnv("HOST_PORT")
}

func (e *envConfig) GetEmailVerificationDuration() int {
	duration, err := strconv.Atoi(utils.GetEnv("EMAIL_VERIFICATION_DURATION"))
	if err != nil {
		return 0 // Default value if parsing fails
	}
	return duration
}

func (e *envConfig) GetDatabaseHost() string {
	return utils.GetEnv("DB_HOST")
}

func (e *envConfig) GetDatabasePort() string {
	return utils.GetEnv("DB_PORT")
}

func (e *envConfig) GetDatabaseUser() string {
	return utils.GetEnv("DB_USER")
}

func (e *envConfig) GetDatabasePassword() string {
	return utils.GetEnv("DB_PASSWORD")
}

func (e *envConfig) GetDatabaseName() string {
	return utils.GetEnv("DB_NAME")
}

func (e *envConfig) GetSalt() string {
	salt := utils.GetEnv("SALT")
	if salt == "" {
		return "Def4u|7" // Default salt value
	}
	return salt
}

func (e *envConfig) GetSupabaseURL() string {
	return strings.TrimSpace(utils.GetEnv("SUPABASE_URL"))
}

func (e *envConfig) GetSupabaseKey() string {
	return strings.TrimSpace(utils.GetEnv("SUPABASE_SERVICE_KEY"))
}

func (e *envConfig) GetSupabaseBucket() string {
	return strings.TrimSpace(utils.GetEnv("SUPABASE_BUCKET_NAME"))
}

func (e *envConfig) GetXenditAPIKey() string {
	return strings.TrimSpace(utils.GetEnv("XENDIT_API_KEY"))
}

func (e *envConfig) GetXenditCallbackToken() string {
	return strings.TrimSpace(utils.GetEnv("XENDIT_CALLBACK_TOKEN"))
}
