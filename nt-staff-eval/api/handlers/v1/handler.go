package v1

import (
	"fmt"

	"github.com/Hatsker01/nt-staff-eval/api/token"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"

	"github.com/Hatsker01/nt-staff-eval/config"
	"github.com/Hatsker01/nt-staff-eval/pkg/logger"
)

type handlerV1 struct {
	db         *sqlx.DB
	log        logger.Logger
	jwtHandler token.JWTHandler
	cfg        config.Config
}

// HandlerV1Config ...
type HandlerV1Config struct {
	Db         *sqlx.DB
	Logger     logger.Logger
	JwtHandler token.JWTHandler
	Cfg        config.Config
}

// New ...
func New(c *HandlerV1Config) *handlerV1 {
	return &handlerV1{
		db:         c.Db,
		log:        c.Logger,
		jwtHandler: c.JwtHandler,
		cfg:        c.Cfg,
	}
}

func GetClaims(h handlerV1, c *gin.Context) (*token.CustomClaims, error) {

	var claims token.CustomClaims

	strToken := c.GetHeader("Authorization")

	token, err := jwt.Parse(strToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(h.jwtHandler.SigninKey), nil
	})

	if err != nil {
		h.log.Error("invalid access token")
		return nil, err
	}
	rawClaims := token.Claims.(jwt.MapClaims)

	claims.Sub = rawClaims["sub"].(string)
	claims.Role = rawClaims["role"].(string)
	claims.Exp = rawClaims["exp"].(float64)
	fmt.Printf("%T type of value in map %v\n", rawClaims["exp"], rawClaims["exp"])
	fmt.Printf("%T type of value in map %v\n", rawClaims["iat"], rawClaims["iat"])

	claims.Iat = rawClaims["iat"].(float64)

	var aud = make([]string, len(rawClaims["aud"].([]interface{})))

	for i, v := range rawClaims["aud"].([]interface{}) {
		aud[i] = v.(string)
	}

	claims.Aud = aud
	claims.Iss = rawClaims["iss"].(string)

	return &claims, nil

}
