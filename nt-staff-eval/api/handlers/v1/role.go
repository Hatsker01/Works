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

// CreateRole ...
// @Summary CreateRole
// @Description This API for creating a new role
// @Tags role
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param role request body structs.CreateRoleStruct true "roleCreateRequest"
// @Success 200 {object} structs.RoleStruct
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/roles/ [post]
func (h *handlerV1) CreateRole(c *gin.Context) {
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
	var (
		req  structs.CreateRoleStruct
		body structs.RoleStruct
	)

	err = c.ShouldBindJSON(&req)
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
	body.Section.Id = req.SectionId
	body.Id = id.String()

	response, err := postgres.NewRoleRepo(h.db).CreateRole(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create role", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetRole ...
// @Summary GetRole
// @Description This API for getting role detail
// @Tags role
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "RoleId"
// @Success 200 {object} structs.RoleStruct
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/roles/{id} [get]
func (h *handlerV1) GetRole(c *gin.Context) {
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

	response, err := postgres.NewRoleRepo(h.db).GetRole(guid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get role", l.Error(err))
		return
	}

	respSection, err := postgres.NewSectionRepo(h.db).GetSection(response.Section.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get section by id for role", l.Error(err))
		return
	}
	response.Section = respSection

	c.JSON(http.StatusOK, response)
}

// GetListRoles ...
// @Summary ListRoles
// @Description This API for getting list of roles
// @Tags role
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query string false "Page"
// @Param limit query string false "Limit"
// @Success 200 {object} []structs.RoleStruct
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/roles [get]
func (h *handlerV1) GetListRoles(c *gin.Context) {
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

	response, count, err := postgres.NewRoleRepo(h.db).GetListRoles(int(params.Page), int(params.Limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to list roles", l.Error(err))
		return
	}

	for i, _ := range response {
		respSection, err := postgres.NewSectionRepo(h.db).GetSection(response[i].Section.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to get section by id for role", l.Error(err))
			return
		}
		response[i].Section = respSection
	}

	c.JSON(http.StatusOK, gin.H{
		"roles": response,
		"count": count,
	})
}

// UpdateRole ...
// @Summary UpdateRole
// @Description This API for updating role
// @Tags role
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "RoleId"
// @Param User request body structs.CreateRoleStruct true "RoleUpdateRequest"
// @Success 200 {object} structs.RoleStruct
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/roles/{id} [put]
func (h *handlerV1) UpdateRole(c *gin.Context) {
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
	var (
		req  structs.CreateRoleStruct
		body structs.RoleStruct
	)

	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	body.Id = c.Param("id")
	body.Section.Id = req.SectionId
	response, err := postgres.NewRoleRepo(h.db).UpdateRole(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update role", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteRole ...
// @Summary DeleteRole
// @Description This API for deleting the role
// @Tags role
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "RoleId"
// @Success 200
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/roles/{id} [delete]
func (h *handlerV1) DeleteRole(c *gin.Context) {
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

	err = postgres.NewRoleRepo(h.db).DeleteRole(guid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete role", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "Deleted",
	})
}
