package config

import (
	"os"

	"github.com/spf13/cast"
)

type Config struct{
	Environment string
	
	UserServiceHost string
	UserServicePort int
	PostSericeHost string
	PostSericePort int
	RedisHost string
	RedisPort int

	//context time in second
	CtxTimeout int

	LogLevel string
	HTTPPort string
}

func Load() Config{
	c:=Config{}

	c.Environment=cast.ToString(getOrReturnDefault("ENVIRONMENT","develop"))


	c.LogLevel=cast.ToString(getOrReturnDefault("LOG_LEVEL","debug"))
	c.HTTPPort=cast.ToString(getOrReturnDefault("HTTP_PORT",":8090"))
	c.UserServiceHost=cast.ToString(getOrReturnDefault("USER_SERVICE_HOST","127.0.0.1"))
	c.UserServicePort=cast.ToInt(getOrReturnDefault("USER_SERVICE_PORT",9000))
	c.RedisHost=cast.ToString(getOrReturnDefault("REDIS_HOST","localhost"))
	c.RedisPort=cast.ToInt(getOrReturnDefault("REDIS_PORT",6379))
	

	c.CtxTimeout=cast.ToInt(getOrReturnDefault("CTX_TIMEOUT",7))
	return c

}

func getOrReturnDefault(key string, defaultValue interface{}) interface{}{
	_,exists:=os.LookupEnv(key)
	if exists{
		return os.Getenv(key)
	}
	return defaultValue
}