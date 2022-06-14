package main

import (
	//"fmt"

	"fmt"

	"github.com/Hatsker01/Works/APi-connection/api/api"
	"github.com/Hatsker01/Works/APi-connection/api/config"
	"github.com/Hatsker01/Works/APi-connection/api/pkg/logger"
	"github.com/Hatsker01/Works/APi-connection/api/services"
	"github.com/gomodule/redigo/redis"

	//	"github.com/Hatsker01/Works/APi-connection/api/storage/redis"
	rds "github.com/Hatsker01/Works/APi-connection/api/storage/redis"
	//"github.com/Hatsker01/Works/APi-connection/api/storage/repo"
	//"github.com/gomodule/redigo/redis"
)

func main() {
	//var redisRepo repo.RepositoryStorage
	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "api")

	pool := redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort))
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}

	redisRepo := rds.NewRedisRepo(&pool)
	serviceManager, err := services.NewServiceManager(&cfg)
	if err != nil {
		log.Error("gRPC dial error", logger.Error(err))
	}

	server := api.New(api.Option{
		Conf:           cfg,
		Logger:         log,
		ServiceManager: serviceManager,
		RedisRepo:      redisRepo,
	})

	if err := server.Run(cfg.HTTPPort); err != nil {
		log.Fatal("failed to run http server", logger.Error(err))
		panic(err)
	}

}
