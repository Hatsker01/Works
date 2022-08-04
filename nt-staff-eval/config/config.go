package config

import (
	"os"

	"github.com/spf13/cast"
)

type Config struct {
	Environment       string
	PostgresHost      string
	PostgresPort      int
	PostgresDatabase  string
	PostgresUser      string
	PostgresPassword  string
	LogLevel          string
	ReviewServiceHost string
	ReviewServicePort int
	HTTPPort          string
	SignKey           string
	PathUserImage     string
	PathSectionImage  string
}

func Load() Config {
	c := Config{}

	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "release"))

	c.HTTPPort = cast.ToString(getOrReturnDefault("HTTP_PORT", ":8080"))

	c.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "134.122.116.0"))
	c.PostgresPort = cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5432))
	c.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DATABASE", "ntstaff"))
	c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "postgres"))
	c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "pgpwd"))

	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "release"))
	c.SignKey = cast.ToString(getOrReturnDefault("SIGN_KEY", "secret"))

	c.PathUserImage = cast.ToString(getOrReturnDefault("PATH_USER_IMAGE", "/server/image/path"))
	c.PathSectionImage = cast.ToString(getOrReturnDefault("PATH_SECTION_IMAGE", "/server/image/path"))
	
	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}
