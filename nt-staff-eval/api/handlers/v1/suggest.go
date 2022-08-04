package v1

import (
	"net/http"

	l "github.com/Hatsker01/nt-staff-eval/pkg/logger"
	"github.com/Hatsker01/nt-staff-eval/pkg/structs"
	"github.com/Hatsker01/nt-staff-eval/pkg/utils"
	"github.com/Hatsker01/nt-staff-eval/storage/postgres"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// CreateSuggest ...
// @Summary Create suggest
// @Description Create suggest
// @Tags suggest
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param suggest body structs.CreateSuggestReq true "suggestCreateRequest"
// @Success 200 {object} structs.Suggest
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/suggests/ [post]
func (h *handlerV1) CreateSuggest(c *gin.Context) {
	claims, err := GetClaims(*h, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized access",
		})
		return
	}

	if claims.Role != "user" && claims.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "token not found",
		})
		return
	}
	var body structs.CreateSuggest
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

	response, err := postgres.NewSuggestRepo(h.db).CreateSuggest(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create suggest", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetSuggest ...
// @Summary Get suggest
// @Description Get suggest
// @Tags suggest
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "suggestId"
// @Success 200 {object} structs.Suggest
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/suggests/{id} [get]
func (h *handlerV1) GetSuggest(c *gin.Context) {
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
	id := c.Param("id")

	response, err := postgres.NewSuggestRepo(h.db).GetSuggest(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get suggest", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetSuggests ...
// @Summary Get suggests
// @Description Get suggests
// @Tags suggest
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Param status query string false "status type: new, active, inactive"
// @Param user_id query string false "User Id"
// @Success 200 {object} []structs.Suggest
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/suggests/ [get]
func (h *handlerV1) GetSuggests(c *gin.Context) {
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
		h.log.Error("failed to parse query params json in suggest get list" + errStr[0])
		return
	}

	response, count, err := postgres.NewSuggestRepo(h.db).GetListSuggests(params.Filters, int(params.Page), int(params.Limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get suggests", l.Error(err))
		return
	}

	for i := range response {
		response[i].User, err = postgres.NewUserRepo(h.db).GetUser(response[i].User.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to get user in suggest get list", l.Error(err))
			return
		}
		response[i].User, err = h.FullUserInform(response[i].User, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to get users inform in suggests get list", l.Error(err))
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"suggests": response,
		"count":    count,
	})
}

// UpdateSuggestStatus ...
// @Summary UpdateSuggestStatus
// @Description This API for updating suggest status
// @Tags suggest
// @Accept json
// @Produce json
// @Param id path string true "suggestId"
// @Security BearerAuth
// @Param User request body structs.UpdateStatusSuggest true "SectionUpdateRequest"
// @Success 200 {object} structs.Suggest
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/suggests/{id} [put]
func (h *handlerV1) UpdateSuggestStatus(c *gin.Context) {
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

	var body structs.UpdateStatusSuggestReq
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	body.Id = c.Param("id")

	response, err := postgres.NewSuggestRepo(h.db).UpdateStatusSuggest(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update suggest", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteSuggest ...
// @Summary DeleteSuggest
// @Description This API for deleting suggest
// @Tags suggest
// @Accept json
// @Produce json
// @Param id path string true "SuggestId"
// @Security BearerAuth
// @Success 200
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/suggests/{id} [delete]
func (h *handlerV1) DeleteSuggest(c *gin.Context) {
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

	err = postgres.NewSuggestRepo(h.db).DeleteSuggest(guid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete suggest", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "Deleted",
	})
}
