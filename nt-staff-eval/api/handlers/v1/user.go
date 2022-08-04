package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	l "github.com/Hatsker01/nt-staff-eval/pkg/logger"
	"github.com/Hatsker01/nt-staff-eval/pkg/structs"
	"github.com/Hatsker01/nt-staff-eval/pkg/utils"
	"github.com/Hatsker01/nt-staff-eval/storage/postgres"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser ...
// @Summary CreateUser
// @Description This API for creating a new user
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user request body structs.CreateUserReq true "userCreateRequestt"
// @Success 200 {object} structs.CreateUser
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/users/ [post]
func (h *handlerV1) CreateUser(c *gin.Context) {
	claims, err := GetClaims(*h, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized access",
		})
		return
	}

	if claims.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "token not found",
		})
		return
	}
	var body structs.CreateUser

	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	id, err := uuid.NewV4()
	if err != nil {
		h.log.Error("failed while generating uuid", l.Error(err))
	}

	body.Id = id.String()

	response, err := postgres.NewUserRepo(h.db).CreateUser(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create user", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetUser ...
// @Summary GetUser
// @Description This API for getting user detail
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "UserId"
// @Success 200 {object} structs.UserStruct
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/users/{id} [get]
func (h *handlerV1) GetUser(c *gin.Context) {
	claims, err := GetClaims(*h, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized access",
		})
		return
	}

	if claims.Role != "admin" && claims.Role != "user" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "token not found",
		})
		return
	}

	guid := c.Param("id")

	response, err := postgres.NewUserRepo(h.db).GetUser(guid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user", l.Error(err))
		return
	}
	response, err = h.FullUserInform(response, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get users inform", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *handlerV1) GetMe(c *gin.Context) {
	var accessToken string
	prefix := "Bearer "
	auth := c.Request.Header["Authorization"][0]
	if auth != "" && strings.HasPrefix(auth, prefix) {
		accessToken = auth[len(prefix):]
	}

	response, err := postgres.NewUserRepo(h.db).HandleToken(accessToken)
	if err != nil {
		h.log.Error("token is invalid", l.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid token",
			"type":  "unauthorized",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *handlerV1) ChangePassword(c *gin.Context) {
	var updatePassword structs.UpdatePassword
	var accessToken string
	var prefix string = "Bearer "
	auth := c.Request.Header["Authorization"][0]
	if auth != "" && strings.HasPrefix(auth, prefix) {
		accessToken = auth[len(prefix):]
	}
	response, err := postgres.NewUserRepo(h.db).HandleToken(accessToken)
	if err != nil {
		h.log.Error("token is invalid", l.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid token",
			"type":  "unauthorized",
		})
		return
	}

	err = c.ShouldBindJSON(updatePassword)
	if err != nil {
		h.log.Error("error while ShouldBindJSON into LoginAdmin", l.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error while ShouldBindJSON into LoginAdmin" + err.Error(),
		})
		return
	}

	changed, err := postgres.NewUserRepo(h.db).ChangePassword(response, updatePassword)
	if changed == false && err == nil {
		h.log.Error("password is invalid", l.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid password",
			"type":  "invalid_password",
		})
		return
	} else if changed == false && err != nil {
		h.log.Error("server error", l.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error with change password",
			"type":  "internal_server_error",
		})
		return
	} else {
		h.log.Error("server error", l.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"message": "password successfully changed",
			"type":    "success",
		})
		return
	}
}

// UpdateUser ...
// @Summary UpdateUser
// @Description This API for updating user
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "UserId"
// @Param User request body structs.UpdateUserFromUser true "UserUpdateRequest"
// @Success 200 {object} structs.UserStruct
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/users/{id} [put]
func (h *handlerV1) UpdateUser(c *gin.Context) {
	claims, err := GetClaims(*h, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized access",
		})
		return
	}

	if claims.Role != "admin" && claims.Role != "user" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "token not found",
		})
		return
	}
	var body structs.UpdateUserFromUser

	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	body.ID = c.Param("id")

	response, err := postgres.NewUserRepo(h.db).UpdateUser(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetListUsers ...
// @Summary ListUsers
// @Description This API for getting list of users
// @Security BearerAuth
// @Tags user
// @Accept json
// @Produce json
// @Param section query string false "Section"
// @Param searchId query string false "SearchId"
// @Param branchId query string false "BranchId"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param placeTop query int false "Place of Top"
// @Success 200 {object} []structs.UserListResp
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/users [get]
func (h *handlerV1) GetListUsers(c *gin.Context) {
	claims, err := GetClaims(*h, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized access",
		})
		return
	}

	if claims.Role != "admin" && claims.Role != "user" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "token not found",
		})
		return
	}
	queryParams := c.Request.URL.Query()

	params, errStr := utils.ParseQueryParams(queryParams)
	if errStr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": errStr[0],
		})
		h.log.Error("failed to parse query params json" + errStr[0])
		return
	}

	response, count, err := postgres.NewUserRepo(h.db).GetListUsers(params.Filters, int(params.Page), int(params.Limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to list users", l.Error(err))
		return
	}
	fmt.Println("123\n\n")

	for i,_ := range response {
		respRole, err := postgres.NewRoleRepo(h.db).GetRole(response[i].Role.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to get role by id for user list", l.Error(err))
			return
		}
		fmt.Println("\nhatsker")
		response[i].Role = respRole
		
		respSection, err := postgres.NewSectionRepo(h.db).GetSection(response[i].Role.Section.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to get section by id for user list", l.Error(err))
			return
		}
		response[i].Role.Section = respSection
	}
	c.JSON(http.StatusOK, gin.H{
		"users": response,
		"count": count,
	})
}

// GetTopUsers ...
// @Summary ListTopUsers
// @Description This API for getting list of top users
// @Security BearerAuth
// @Tags user
// @Accept json
// @Produce json
// @Param page query string false "Page"
// @Param limit query string false "Limit"
// @Success 200 {object} []structs.UserStruct
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/top_users [get]
func (h *handlerV1) GetTopUsers(c *gin.Context) {
	claims, err := GetClaims(*h, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized access",
		})
		return
	}

	if claims.Role != "admin" && claims.Role != "user" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "token not found",
		})
		return
	}
	queryParams := c.Request.URL.Query()

	params, errStr := utils.ParseQueryParams(queryParams)
	if errStr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": errStr[0],
		})
		h.log.Error("failed to parse query params json at get top user" + errStr[0])
		return
	}

	response, count, err := postgres.NewUserRepo(h.db).GetTopUsers(int(params.Page), int(params.Limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to list top users", l.Error(err))
		return
	}

	for i, user := range response {
		response[i], err = h.FullUserInform(user, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to list top users get user inform", l.Error(err))
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"users": response,
		"count": count,
	})
}

// LoginUser ...
// @Summary LoginUser
// @Tags login
// @Accept json
// @Produce json
// @Param loginData body structs.LoginUser true "login data"
// @Success 200 {object} structs.UserLoginResp
// @Failure 400 {object} structs.StandardErrorModel
// @Router /v1/login [post]
func (h *handlerV1) LoginUser(c *gin.Context) {
	var (
		admin         structs.AdminLogin
		userToken     structs.UserLoginResp
		user          structs.LoginUser
		adminEntering = true
	)

	err := c.ShouldBindJSON(&admin)
	if err != nil {
		h.log.Error("error while ShouldBindJSON into LoginAdmin", l.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error while ShouldBindJSON into LoginAdmin" + err.Error(),
		})
		return
	}

	adminToken, err := postgres.NewAdminRepo(h.db).Login(admin)
	if err != nil {
		adminEntering = false
		user.Email = admin.Login
		user.Password = admin.Password
		userToken, err = postgres.NewUserRepo(h.db).LoginUser(user)
		if err != nil {
			h.log.Error("error while logging into LoginUser", l.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"type":    "bad_request",
				"message": err.Error(),
			})
			return
		}
	}
	if !adminEntering {
		err = bcrypt.CompareHashAndPassword([]byte(userToken.Password), []byte(user.Password))
		//err = bcrypt.CompareHashAndPassword([]byte("$2a$10$C0tnO6XwPJC3mmZYzL/BKuKZ/95jqPHYC6FsY5U2It7pSz8ozx3iC"), []byte(user.Password))
		if err != nil {
			h.log.Error("error while hashing password", l.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"type":    "bad_request",
				"message": "wrong username or password",
			})
			return
		}
	}

	if adminEntering {
		h.jwtHandler.Sub = strconv.Itoa(adminToken.ID)
		h.jwtHandler.Role = "admin"
		h.jwtHandler.Aud = []string{"admin_profile"}
	} else {
		h.jwtHandler.Sub = userToken.Id
		h.jwtHandler.Role = "user"
		h.jwtHandler.Aud = []string{"user_profile"}
	}

	access, refresh, err := h.jwtHandler.GenerateAuthJWT()
	if err != nil {
		h.log.Error("error while generating tokens LoginUser ", l.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error while generating tokens LoginUser " + err.Error(),
		})
		return
	}

	if adminEntering {
		adminToken.RefreshToken = refresh
		adminToken.AccessToken = access

		err = postgres.NewAdminRepo(h.db).Update(adminToken)
		if err != nil {
			h.log.Error("error while updating tokens at LoginAdmin", l.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error while updating tokens LoginAdmin" + err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"access_token": adminToken.AccessToken,
		})
	} else {
		userAuth := structs.UserAuth{
			Id:           userToken.Id,
			AccessToken:  access,
			RefreshToken: refresh,
		}

		err = postgres.NewUserRepo(h.db).UpdateToken(userAuth)
		if err != nil {
			h.log.Error("error while updating tokens LoginUser ", l.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "error while updating tokens LoginUser " + err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"access_token": userAuth.AccessToken,
		})
	}
}

func (h *handlerV1) FullUserInform(user structs.UserStruct, c *gin.Context) (structs.UserStruct, error) {
	if user.Branch.Id != 0 {
		respBranch, err := postgres.NewBranchRepo(h.db).GetBranch(user.Branch.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to get branch by id for user", l.Error(err))
			return structs.UserStruct{}, err
		}
		user.Branch = respBranch
	}
	respRole, err := postgres.NewRoleRepo(h.db).GetRole(user.Role.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get role by id for user", l.Error(err))
		return structs.UserStruct{}, err
	}
	user.Role = respRole
	respSection, err := postgres.NewSectionRepo(h.db).GetSection(user.Role.Section.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get section by id for user", l.Error(err))
		return structs.UserStruct{}, err
	}
	user.Role.Section = respSection

	return user, nil
}
