package api

import (
	"github.com/Hatsker01/nt-staff-eval/api/token"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/Hatsker01/nt-staff-eval/api/docs" // swag
	v1 "github.com/Hatsker01/nt-staff-eval/api/handlers/v1"
	"github.com/Hatsker01/nt-staff-eval/config"
	"github.com/Hatsker01/nt-staff-eval/pkg/logger"
)

type Option struct {
	Db     *sqlx.DB
	Conf   config.Config
	Logger logger.Logger
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func Routers(option Option) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())

	jwt := token.JWTHandler{
		SigninKey: option.Conf.SignKey,
		Log:       option.Logger,
	}

	handlerV1 := v1.New(&v1.HandlerV1Config{
		Db:         option.Db,
		Logger:     option.Logger,
		JwtHandler: jwt,
		Cfg:        option.Conf,
	})

	api := router.Group("/v1")
	api.Static("/media", "./media")

	// Evaluation
	api.POST("/evaluations", handlerV1.CreateEvaluation)
	api.GET("/evaluations/:id", handlerV1.GetEvaluation)
	api.GET("/evaluations", handlerV1.GetListEvaluations)
	api.PUT("/evaluations/:id", handlerV1.UpdateEvaluation)
	api.DELETE("/evaluations/:id", handlerV1.DeleteEvaluation)
	//
	//// Behavior
	//api.POST("/behaviors", handlerV1.CreateBehavior)
	//api.GET("/behaviors/:id", handlerV1.GetBehavior)
	//api.GET("/behaviors", handlerV1.GetListBehaviors)
	//api.PUT("/behaviors/:id", handlerV1.UpdateBehavior)
	//api.DELETE("/behaviors/:id", handlerV1.DeleteBehavior)

	// Section
	api.POST("/sections", handlerV1.CreateSection)
	api.GET("/sections/:id", handlerV1.GetSection)
	api.GET("/sections", handlerV1.GetListSections)
	api.GET("/sections/", handlerV1.GetListSections)
	api.PUT("/sections/:photo", handlerV1.UpdateSection)
	api.DELETE("/sections/:id", handlerV1.DeleteSection)

	// Role
	api.POST("/roles", handlerV1.CreateRole)
	api.GET("/roles/:id", handlerV1.GetRole)
	api.GET("/roles", handlerV1.GetListRoles)
	api.PUT("/roles/:id", handlerV1.UpdateRole)
	api.DELETE("/roles/:id", handlerV1.DeleteRole)

	// Rated
	api.POST("/rateds", handlerV1.CreateRated)
	api.GET("/rateds", handlerV1.GetListRateds)
	api.DELETE("/rateds/:id", handlerV1.DeleteRated)

	// Users
	api.POST("/login", handlerV1.LoginUser)
	api.POST("/users", handlerV1.CreateUser)
	api.GET("/users", handlerV1.GetListUsers)
	api.GET("/top_users", handlerV1.GetTopUsers)
	api.GET("/users/:id", handlerV1.GetUser)
	api.PUT("/change/password", handlerV1.ChangePassword)
	api.GET("/me", handlerV1.GetMe)
	api.PUT("/users/:id", handlerV1.UpdateUser)

	// Branch
	api.POST("/branches", handlerV1.CreateBranch)
	api.GET("/branches/:id", handlerV1.GetBranch)
	api.GET("/branches", handlerV1.GetListBranches)
	api.PUT("/branches/:id", handlerV1.UpdateBranch)
	api.DELETE("/branches/:id", handlerV1.DeleteBranch)

	// Suggest
	api.POST("/suggests", handlerV1.CreateSuggest)
	api.GET("/suggests/:id", handlerV1.GetSuggest)
	api.GET("/suggests", handlerV1.GetSuggests)
	api.PUT("/suggests/:id", handlerV1.UpdateSuggestStatus)
	api.DELETE("/suggests/:id", handlerV1.DeleteSuggest)

	// News
	api.POST("/news", handlerV1.CreateNews)
	api.GET("/news/:id", handlerV1.GetNews)
	api.GET("/news", handlerV1.GetAllNews)
	api.PUT("/news/:id", handlerV1.UpdateNews)
	api.DELETE("/news/:id", handlerV1.DeleteNews)

	// Category
	api.POST("/categories", handlerV1.CreateNewsCategory)
	api.GET("/categories", handlerV1.GetAllNewsCategory)
	api.DELETE("/categories/:id", handlerV1.DeleteNews)

	// Image
	api.POST("images/users/upload/:id", handlerV1.UploadUserImage)
	api.POST("images/sections/upload/:id", handlerV1.UploadSectionImage)

	//// Admin
	//api.POST("/admin", handlerV1.LoginAdmin)

	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
