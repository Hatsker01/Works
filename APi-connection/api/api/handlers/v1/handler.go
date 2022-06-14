package v1

import (
	"github.com/Hatsker01/Works/APi-connection/api/config"
	"github.com/Hatsker01/Works/APi-connection/api/pkg/logger"
	"github.com/Hatsker01/Works/APi-connection/api/services"
	"github.com/Hatsker01/Works/APi-connection/api/storage/repo"
)

type handlerV1 struct {
	log            logger.Logger
	serviceManager services.IServiceManager
	cfg            config.Config
	redisStorage   repo.RepositoryStorage
}

type HandlerV1Config struct {
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Cfg            config.Config
	Redis          repo.RepositoryStorage
}

func New(c *HandlerV1Config) *handlerV1 {
	return &handlerV1{
		log:            c.Logger,
		serviceManager: c.ServiceManager,
		cfg:            c.Cfg,
		redisStorage:   c.Redis,
	}
}

// func New1(c *HandlerV1Config)*handlerV1{
// 	return &handlerV1{
// 		log: c.Logger,
// 		serviceManager: c.ServiceManager,
// 		cfg: c.Cfg,
// 	}
// }
