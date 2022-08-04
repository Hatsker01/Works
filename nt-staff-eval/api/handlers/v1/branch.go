package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	l "github.com/Hatsker01/nt-staff-eval/pkg/logger"
	"github.com/Hatsker01/nt-staff-eval/pkg/structs"
	"github.com/Hatsker01/nt-staff-eval/storage/postgres"
)

// CreateBranch ...
// @Summary CreateBranch
// @Description This API for creating a new branch
// @Tags branch
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param branch request body structs.CreateBranch true "BranchCreateRequest"
// @Success 200 {object} structs.BranchStruct
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/branches/ [post]
func (h *handlerV1) CreateBranch(c *gin.Context) {
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
	var body structs.CreateBranch

	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	response, err := postgres.NewBranchRepo(h.db).CreateBranch(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create branch", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetBranch ...
// @Summary GetBranch
// @Description This API for getting branch detail
// @Tags branch
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "BranchId"
// @Success 200 {object} structs.BranchStruct
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/branches/{id} [get]
func (h *handlerV1) GetBranch(c *gin.Context) {
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
	guid64 := int64(guid)
	if err != nil {
		h.log.Error("Failed to parse string to int", l.Error(err))
	}

	response, err := postgres.NewBranchRepo(h.db).GetBranch(guid64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get branch", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetListBranches ...
// @Summary ListBranches
// @Description This API for getting list of branches
// @Tags branch
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query string false "Page"
// @Param limit query string false "Limit"
// @Success 200 {object} []structs.BranchStruct
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/branches [get]
func (h *handlerV1) GetListBranches(c *gin.Context) {
	//queryParams := c.Request.URL.Query()
	//
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
	response, count, err := postgres.NewBranchRepo(h.db).GetListBranch()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to list branch", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"branches": response,
		"count":    count,
	})
}

// UpdateBranch ...
// @Summary UpdateBranch
// @Description This API for updating branch
// @Tags branch
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "BranchId"
// @Param User request body structs.UpdateBranch true "BranchUpdateRequest"
// @Success 200 {object} structs.BranchStruct
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/branches/{id} [put]
func (h *handlerV1) UpdateBranch(c *gin.Context) {
	var body structs.BranchStruct

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

	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	id64, err := strconv.Atoi(c.Param("id"))
	body.Id = int64(id64)
	if err != nil {
		h.log.Error("Failed to parse string to int", l.Error(err))
	}

	response, err := postgres.NewBranchRepo(h.db).UpdateBranch(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update branch", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteBranch ...
// @Summary DeleteBranch
// @Description This API for deleting the branch
// @Tags branch
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "BranchId"
// @Success 200
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/branches/{id} [delete]
func (h *handlerV1) DeleteBranch(c *gin.Context) {
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
	guid64 := int64(guid)
	if err != nil {
		h.log.Error("Failed to parse string to int", l.Error(err))
	}

	err = postgres.NewBranchRepo(h.db).DeleteBranch(guid64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete branch", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "Deleted",
	})
}
