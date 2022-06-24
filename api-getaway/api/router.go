package api

import (
	v1 "github.com/Hatsker01/Works/api-getaway/api/handlers/v1"
	"github.com/Hatsker01/Works/api-getaway/config"
	"github.com/Hatsker01/Works/api-getaway/pkg/logger"
	"github.com/Hatsker01/Works/api-getaway/services"
	"github.com/gin-gonic/gin"
)

// Option ...
type Option struct {
	Conf           config.Config
	Logger         logger.Logger
	ServiceManager services.IServiceManager
}

// New ...
func New(option Option) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	handlerV1 := v1.New(&v1.HandlerV1Config{
		Logger:         option.Logger,
		ServiceManager: option.ServiceManager,
		Cfg:            option.Conf,
	})

	api := router.Group("/v1")
	api.POST("/posts", handlerV1.Create)
	api.GET("/posts/update", handlerV1.Update)
	api.DELETE("/posts/delete/:id",handlerV1.Delete)
	api.GET("/posts/getbyid/:id",handlerV1.GetPostbyId)
	api.GET("/posts/all",handlerV1.GetAll)
	
	
	// api.GET("/users", handlerV1.ListUsers)
	// api.PUT("/users/:id", handlerV1.UpdateUser)
	// api.DELETE("/users/:id", handlerV1.DeleteUser)

	return router
}
