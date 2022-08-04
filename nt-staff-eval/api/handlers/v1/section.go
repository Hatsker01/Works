package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	l "github.com/Hatsker01/nt-staff-eval/pkg/logger"
	"github.com/Hatsker01/nt-staff-eval/pkg/structs"
	"github.com/Hatsker01/nt-staff-eval/storage/postgres"
)

// CreateSection ...
// @Summary CreateSection
// @Description This API for creating a new section
// @Tags section
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param section request body structs.CreateSectionStruct true "sectionCreateRequest"
// @Success 200 {object} structs.SectionStruct
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/sections/ [post]
func (h *handlerV1) CreateSection(c *gin.Context) {
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
	var body structs.SectionStruct

	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	response, err := postgres.NewSectionRepo(h.db).CreateSection(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create section", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetSection ...
// @Summary GetSection
// @Description This API for getting section detail
// @Tags section
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "SectionId"
// @Success 200 {object} structs.SectionStruct
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/sections/{id} [get]
func (h *handlerV1) GetSection(c *gin.Context) {
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
	guid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.log.Error("Failed to parse string to int", l.Error(err))
	}

	response, err := postgres.NewSectionRepo(h.db).GetSection(guid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get section", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetListSections ...
// @Summary ListSections
// @Description This API for getting list of sections
// @Tags section
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} []structs.SectionStruct
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/sections [get]
func (h *handlerV1) GetListSections(c *gin.Context) {
	//queryParams := c.Request.URL.Query()

	//params, errStr := utils.ParseQueryParams(queryParams)
	//if errStr != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"err": errStr[0],
	//	})
	//	h.log.Error("failed to parse query params json" + errStr[0])
	//	return
	//}
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

	response, count, err := postgres.NewSectionRepo(h.db).GetListSections()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to list sections", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sections": response,
		"count":    count,
	})
}

// UpdateSection ...
// @Summary UpdateSection
// @Description This API for updating section
// @Tags section
// @Accept json
// @Produce json
// @Param id path string true "SectionId"
// @Security BearerAuth
// @Param User request body structs.CreateSectionStruct true "SectionUpdateRequest"
// @Success 200 {object} structs.SectionStruct
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/sections/{id} [put]
func (h *handlerV1) UpdateSection(c *gin.Context) {
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

	var body structs.SectionStruct

	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	body.Id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		h.log.Error("failed to parse int to string", l.Error(err))
	}

	response, err := postgres.NewSectionRepo(h.db).UpdateSection(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update section", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteSection ...
// @Summary DeleteSection
// @Description This API for deleting the section
// @Tags section
// @Accept json
// @Produce json
// @Param id path string true "SectionId"
// @Security BearerAuth
// @Success 200
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/sections/{id} [delete]
func (h *handlerV1) DeleteSection(c *gin.Context) {
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

	guid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.log.Error("failed to parse int to string", l.Error(err))
	}

	err = postgres.NewSectionRepo(h.db).DeleteSection(guid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete section", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "Deleted",
	})
}
