package v1

import (
	"net/http"

	l "github.com/Hatsker01/nt-staff-eval/pkg/logger"
	"github.com/Hatsker01/nt-staff-eval/pkg/structs"
	"github.com/Hatsker01/nt-staff-eval/storage/postgres"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// CreateNewsCategory ...
// @Summary Create news category
// @Description Create news category
// @Tags category news
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param news_category body structs.Category true "newsCategoryCreateRequest"
// @Success 200 {object} structs.Category
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/categories [post]
func (h *handlerV1) CreateNewsCategory(c *gin.Context) {
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

	var body structs.Category
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json in news category", l.Error(err))
		return
	}

	id, err := uuid.NewV4()
	if err != nil {
		h.log.Error("failed while generating uuid in news category", l.Error(err))
	}

	body.Id = id.String()

	response, err := postgres.NewNewsCategoryRepo(h.db).CreateNewsCategory(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create news category", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetAllNewsCategory ...
// @Summary Get all news category
// @Description Get all news category
// @Tags category news
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} structs.Category
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/categories [get]
func (h *handlerV1) GetAllNewsCategory(c *gin.Context) {
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

	response, count, err := postgres.NewNewsCategoryRepo(h.db).GetListNewsCategory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get all news category", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  response,
		"count": count,
	})
}

// DeleteCategory ...
// @Summary Delete news category
// @Description Delete news category
// @Tags category news
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "id"
// @Success 200 {object} structs.StandardErrorModel
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/categories/{id} [delete]
func (h *handlerV1) DeleteCategory(c *gin.Context) {
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

	id := c.Param("id")

	err = postgres.NewNewsCategoryRepo(h.db).DeleteNewsCategory(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete news category", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, "Deleted")
}
