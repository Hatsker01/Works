package v1

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"unicode"

	//pl "github.com/Hatsker01/Works/Api-service_user/api/api/model"
	pb "github.com/Hatsker01/Works/Api-service_user/api/genproto"
	emai "github.com/Hatsker01/Works/Api-service_user/api/mail"
	l "github.com/Hatsker01/Works/Api-service_user/api/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"google.golang.org/protobuf/encoding/protojson"
)

var code string

//Register register user
//@Summary Register user summary
//Description This api for registering user
//@Tags user
//@Accept json
//@Produce json
//@Param user body User true "user body"
//@Success 200 {string} Success!
//@Router /v1/users/registeruser [post]
func (h *handlerV1) RegisterUser(c *gin.Context) {
	var (
		body        User
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to blind json", l.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	check, err := h.serviceManager.UserService().CheckField(ctx, &pb.CheckFieldRequest{
		Field: `username`,
		Value: body.Username,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			`error`: err.Error(),
		})
		h.log.Error("failed to check username", l.Error(err))
		return

	}
	if !check.Check {
		check1, err := h.serviceManager.UserService().CheckField(ctx, &pb.CheckFieldRequest{
			Field: `email`,
			Value: body.Email,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to check email", l.Error(err))
			return

		}
		if check1.Check {

			return
		}

	} else {
		return
	}

	code, _ = genCaptchaCode()
	if err != nil {
		fmt.Println(err)
		return
	}

	eigthMore, number, upper, special, moredigits := verifyPassword(body.Password)
	if !eigthMore {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed in password not much characters",
		})
		h.log.Error("failed in password not much characters", l.Error(err))
		return
	}
	if !moredigits {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed in password not much characters",
		})
		h.log.Error("failed in password not much characters", l.Error(err))
		return
	}
	if !number {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed in password don't have numbers in password",
		})
		h.log.Error("failed in password don't have numbers in password", l.Error(err))
		return
	}
	if !upper {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed in password don't have upper symbole",
		})
		h.log.Error("failed in password don't have upper symbole", l.Error(err))
		return
	}
	if !special {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed in password don't have special sybole",
		})
		h.log.Error("failed in password don't have special sybole", l.Error(err))
		return
	}

	body.Code = code
	fmt.Println(code)
	//src := []byte("Hello Gopher!")

	dst := make([]byte, hex.EncodedLen(len(body.Password)))
	hex.Encode(dst, []byte(body.Password))
	body.Password = string(dst)
	//body.Password=
	bodyByte, err := json.Marshal(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to convert json", l.Error(err))
		return

	}

	//users := User{}
	err = h.redisStorage.SetWithTTL(body.Email, string(bodyByte), int64(time.Minute)*2)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to convert json2", l.Error(err))
		return

	}

	//value, err := h.redisStorage.Get(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to convert json2", l.Error(err))
		return

	}
	err = emai.SendMail(code, body.Email)
	if err != nil {
		fmt.Println(err)
	}
	genCaptchaCode()

	// }

}

var coded string
var email string

//Post user by code
//@Summary Post user summary
//Description This api for post user by code
//@Tags user
//@Accept json
//@Produce json
//@Param email path string true "Email"
//@Param coded path string true "Code"
//@Success 200 {string} User!
//@Router /v1/users/register/user/{email}/{codedw} [post]
func (h *handlerV1) Verify(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true
	email = c.Param("email")
	coded = c.Param("coded")
	fmt.Println(email)
	fmt.Println(coded, "   ", code)
	var (
		userm pb.User
	)

	vali, _ := redis.String(h.redisStorage.Get(email))
	err := json.Unmarshal([]byte(vali), &userm)
	if err != nil {
		return
	}
	//fmt.Println(val)
	ctxr, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	if userm.Code == coded {
		usersss, err := h.serviceManager.UserService().CreateUser(ctxr, &userm)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to create user", l.Error(err))
			return
		}
		// bodyByte, err := json.Marshal(err)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"error": err.Error(),
		// 	})
		// 	h.log.Error("failed to convert json", l.Error(err))
		// 	return
		// }

		// bodyByte, err := json.Marshal(value)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"error": err.Error(),
		// 	})
		// 	h.log.Error("failed to convert json", l.Error(err))
		// 	return

		// }
		c.JSON(http.StatusCreated, usersss)
	} else {
		err := "code erroe"
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		//h.log.Error("failed in code", nil)
		return

	}

}

type Register struct {
	User pb.User
	Code string `json:"code"`
}

func genCaptchaCode() (string, error) {
	codes := make([]byte, 6)
	if _, err := rand.Read(codes); err != nil {
		return "", err
	}

	for i := 0; i < 6; i++ {
		codes[i] = uint8(48 + (codes[i] % 10))
	}

	return string(codes), nil
}

func verifyPassword(s string) (eigthMore, number, upper, special, moredigits bool) {
	letters := 0
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
			letters++
		default:
			//return false, false, false, false
		}
	}
	eigthMore = letters >= 8
	moredigits = letters <= 32
	return
}
