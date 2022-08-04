package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"

	l "github.com/Hatsker01/nt-staff-eval/pkg/logger"
	"github.com/Hatsker01/nt-staff-eval/pkg/structs"
	"github.com/Hatsker01/nt-staff-eval/pkg/utils"
	"github.com/Hatsker01/nt-staff-eval/storage/postgres"
)

// CreateRated ...
// @Summary CreateRated
// @Description This API for creating a new rated
// @Tags rated
// @Accept json
// @Produce json
// @Param rated request body structs.CreateRatedReq true "ratedCreateRequest"
// @Success 200 {object} structs.CreateRated
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/rateds/ [post]
func (h *handlerV1) CreateRated(c *gin.Context) {
	var body structs.CreateRated

	err := c.ShouldBindJSON(&body)
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

	err = postgres.NewRatedRepo(h.db).CreateRated(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create Rated", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Success",
	})
}

// GetListRateds ...
// @Summary ListRateds
// @Description This API for getting list of rateds
// @Tags rated
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query string false "Page"
// @Param limit query string false "Limit"
// @Success 200 {object} []structs.Rated
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/rateds [get]
func (h *handlerV1) GetListRateds(c *gin.Context) {
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
	queryParams := c.Request.URL.Query()

	params, errStr := utils.ParseQueryParams(queryParams)
	if errStr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": errStr[0],
		})
		h.log.Error("failed to parse query params json" + errStr[0])
		return
	}

	response, count, err := postgres.NewRatedRepo(h.db).GetListRateds(int(params.Page), int(params.Limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to list rateds", l.Error(err))
		return
	}

	for i := range response {
		user, err := postgres.NewUserRepo(h.db).GetUser(response[i].User.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to list rateds get user by id", l.Error(err))
			return
		}
		response[i].User, err = h.FullUserInform(user, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to list rateds full inform by user", l.Error(err))
			return
		}
		//for j := range response[i].Evaluations {
		//	fmt.Println(response[i].Evaluations[j].Section)
		//	response[i].Evaluations[j].Section, err = postgres.NewSectionRepo(h.db).GetSection(response[i].Evaluations[j].Section.Id)
		//	if err != nil {
		//		c.JSON(http.StatusInternalServerError, gin.H{
		//			"error": err.Error(),
		//		})
		//		h.log.Error("failed to list rateds evaluations section", l.Error(err))
		//		return
		//	}
		//}
	}

	c.JSON(http.StatusOK, gin.H{
		"rateds": response,
		"count":  count,
	})
}

// DeleteRated ...
// @Summary DeleteRated
// @Description This API for deleting the rated
// @Tags rated
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "RatedID"
// @Success 200
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/rateds/{id} [delete]
func (h *handlerV1) DeleteRated(c *gin.Context) {
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
	guid := c.Param("id")

	err = postgres.NewRatedRepo(h.db).DeleteRated(guid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete rated", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "Deleted",
	})
}
