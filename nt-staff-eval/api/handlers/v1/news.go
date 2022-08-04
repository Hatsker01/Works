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

// CreateNews ...
// @Summary Create news
// @Description Create news
// @Tags news
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param news body structs.CreateNewsReq true "newsCreateRequest"
// @Success 200 {object} structs.News
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/news [post]
func (h *handlerV1) CreateNews(c *gin.Context) {
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
	var body structs.CreateNews
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json in news", l.Error(err))
		return
	}

	id, err := uuid.NewV4()
	if err != nil {
		h.log.Error("failed while generating uuid in news", l.Error(err))
	}

	body.Id = id.String()

	response, err := postgres.NewNewsRepo(h.db).CreateNews(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create news", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetNews ...
// @Summary Get news
// @Description Get news
// @Tags news
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "newsId"
// @Success 200 {object} structs.News
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/news/{id} [get]
func (h *handlerV1) GetNews(c *gin.Context) {
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
	response, err := postgres.NewNewsRepo(h.db).GetNews(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get news", l.Error(err))
		return
	}
	if response.Author.Id != "" {
		response.Author, err = postgres.NewUserRepo(h.db).GetUser(response.Author.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to get news", l.Error(err))
			return
		}
		response.Author, err = h.FullUserInform(response.Author, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to get news to full user inform", l.Error(err))
			return
		}
	}
	category, err := postgres.NewNewsCategoryRepo(h.db).GetNewsCategory(response.Category.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get news category in news", l.Error(err))
		return
	}
	response.Category.Name = category.Name

	c.JSON(http.StatusOK, response)
}

// GetAllNews ...
// @Summary Get all news
// @Description Get all news
// @Tags news
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Param title query string false "title"
// @Param category_id query string false "category_id"
// @Param author_id query string false "author_id"
// @Success 200 {object} structs.News
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/news [get]
func (h *handlerV1) GetAllNews(c *gin.Context) {
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

	response, count, err := postgres.NewNewsRepo(h.db).GetListNews(params.Filters, int(params.Page), int(params.Limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get news", l.Error(err))
		return
	}

	for i := range response {
		response[i].Author, err = postgres.NewUserRepo(h.db).GetUser(response[i].Author.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to get news", l.Error(err))
			return
		}
		response[i].Author, err = h.FullUserInform(response[i].Author, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to get news to full user inform", l.Error(err))
			return
		}
		category, err := postgres.NewNewsCategoryRepo(h.db).GetNewsCategory(response[i].Category.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to get news category in news list", l.Error(err))
			return
		}
		response[i].Category.Name = category.Name
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  response,
		"count": count,
	})
}

// UpdateNews ...
// @Summary Update news
// @Description Update news
// @Tags news
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "newsId"
// @Param body body structs.CreateNewsReq true "news"
// @Success 200 {object} structs.Data
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/news/{id} [put]
func (h *handlerV1) UpdateNews(c *gin.Context) {
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
	var news structs.CreateNews
	if err := c.ShouldBindJSON(&news); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	news.Id = id
	response, err := postgres.NewNewsRepo(h.db).UpdateNews(news)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update news", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, structs.Data{Data: response, Message: "Success Updated"})
}

// DeleteNews ...
// @Summary Delete news
// @Description Delete news
// @Tags news
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "newsId"
// @Success 200 {object} structs.News
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/news/{id} [delete]
func (h *handlerV1) DeleteNews(c *gin.Context) {
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
	err = postgres.NewNewsRepo(h.db).DeleteNews(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete news", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, "Successful 'deleted'")
}
