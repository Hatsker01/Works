package api

import (
	//v1 "github.com/Hatsker01/Works/api-token/handlers/v1"
	"github.com/Hatsker01/Works/api-token/config"
	"github.com/Hatsker01/Works/api-token/pkg/logger"
	"github.com/Hatsker01/Works/api-token/services"
	"github.com/Hatsker01/Works/api-token/storage/repo"
	casbinN "github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"

	"github.com/Hatsker01/Works/api-token/auth"
	"github.com/Hatsker01/Works/api-token/casbin"
	_ "github.com/Hatsker01/Works/api-token/docs"
	v1 "github.com/Hatsker01/Works/api-token/handlers/v1"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Option struct {
	Conf           config.Config
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	RedisRepo      repo.RepositoryStorage
	Casbin         *casbinN.Enforcer
}

// New @BasePath /v1
// New ...
// @SecurityDefinitions.apikey BearerAuth
// @Description GetMyProfile
// @in header
// @name Authorization
func New(option Option) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	jwtHandler := auth.JwtHandler{
		SigninKey: option.Conf.SigninKey,
		Log:       option.Logger,
	}
	router.Use(casbin.NewJwtRoleStruct(option.Casbin, option.Conf, jwtHandler))

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
	api.POST("/users/registeruser", handlerV1.RegisterUser)
	api.POST("users/register/user/:email/:coded", handlerV1.Verify)
	api.GET("/users/login/user", handlerV1.Login)
	api.GET("/user/user/user", handlerV1.Useruse)
	//	api.GET("/users/lala/:first_name",handlerV1.SearchUser)
	url := ginSwagger.URL("swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	//api.POST("/post",handlerV1.CreatePost)

	return router
}
