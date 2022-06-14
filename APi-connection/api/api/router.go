package api

import (
	v1 "github.com/Hatsker01/Works/APi-connection/api/api/handlers/v1"
	"github.com/Hatsker01/Works/APi-connection/api/config"
	"github.com/Hatsker01/Works/APi-connection/api/pkg/logger"
	"github.com/Hatsker01/Works/APi-connection/api/services"
	"github.com/Hatsker01/Works/APi-connection/api/storage/repo"
	"github.com/gin-gonic/gin"

	_ "github.com/Hatsker01/Works/APi-connection/api/api/docs"
	"github.com/swaggo/gin-swagger"
)

type Option struct {
	Conf           config.Config
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	RedisRepo      repo.RepositoryStorage
}

func New(option Option) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	handlerV1 := v1.New(&v1.HandlerV1Config{
		Logger:         option.Logger,
		ServiceManager: option.ServiceManager,
		Cfg:            option.Conf,
		Redis:          option.RedisRepo,
	})
	api := router.Group("/v1")
	api.POST("/users", handlerV1.CreateUser)
	api.GET("/users/:id", handlerV1.GetUser)
	api.DELETE("/users/delete/:id", handlerV1.DeleteUser)
	api.PUT("/users/update/:id", handlerV1.UpdateUser)
	api.GET("/users/alluser", handlerV1.GetAllUser)
	api.GET("/users/users", handlerV1.GetListUsers)
	//	api.GET("/users/lala/:first_name",handlerV1.SearchUser)
	url := ginSwagger.URL("swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	//api.POST("/post",handlerV1.CreatePost)

	return router
}
