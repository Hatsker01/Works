package v1

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	l "github.com/Hatsker01/nt-staff-eval/pkg/logger"
	"github.com/Hatsker01/nt-staff-eval/pkg/structs"
	"github.com/Hatsker01/nt-staff-eval/storage/postgres"
)

func (h *handlerV1) TokenValid(c *gin.Context) {
	tokenAuth, err := postgres.NewAuthRepo(h.db).ExtractTokenMetadata(c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Please login first"})
		return
	}

	userID, err := postgres.NewAuthRepo(h.db).FetchAuth(tokenAuth)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Please login first"})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	//To be called from GetUserID()
	c.Set("UserID", userID)
}

func (h *handlerV1) Refresh(c *gin.Context) {
	var tokenForm structs.Token

	if c.ShouldBindJSON(&tokenForm) != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Invalid form",
			"form":    tokenForm,
		})
		h.log.Error("failed to bind json")
		c.Abort()
		return
	}

	// verify the token
	token, err := jwt.Parse(tokenForm.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid authorization, please login again",
			"error":   err.Error(),
		})
		h.log.Error("failed to v1/auth.go/Refresh", l.Error(err))
		return
	}

	// is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid authorization, please login again",
			"error":   err.Error(),
		})
		h.log.Error("failed to v1/auth.go/Refresh", l.Error(err))
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		refreshUUID, ok := claims["refresh_uuid"].(string) //convert interface to string
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid authorization, please login again",
				"error":   err.Error(),
			})
			h.log.Error("failed to v1/auth.go/Refresh", l.Error(err))
			return
		}
		//Delete the precious Refresh Token
		deleted, delErr := postgres.NewAuthRepo(h.db).DeleteAuth(refreshUUID)
		if delErr != nil || deleted == 0 { // if any goes wrong
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid authorization, please login again",
				"error":   err.Error(),
			})
			h.log.Error("failed to v1/auth.go/Refresh", l.Error(err))
			return
		}

	}
}
