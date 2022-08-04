package v1

//
//import (
//	"net/http"
//
//	"github.com/gin-gonic/gin"
//	"github.com/gofrs/uuid"
//
//	l "github.com/Hatsker01/nt-staff-eval/pkg/logger"
//	"github.com/Hatsker01/nt-staff-eval/pkg/structs"
//	"github.com/Hatsker01/nt-staff-eval/pkg/utils"
//	"github.com/Hatsker01/nt-staff-eval/storage/postgres"
//)
//
//// CreateBehavior ...
//// @Summary CreateBehavior
//// @Description This API for creating a new behavior
//// @Tags behavior
//// @Accept json
//// @Produce json
//// @Param behavior request body structs.BehaviorStruct true "behaviorCreateRequest"
//// @Success 200 {object} structs.BehaviorStruct
//// @Failure 400 {object} structs.StandardErrorModel
//// @Failure 500 {object} structs.StandardErrorModel
//// @Router /v1/behaviors/ [post]
//func (h *handlerV1) CreateBehavior(c *gin.Context) {
//	var body structs.BehaviorStruct
//
//	err := c.ShouldBindJSON(&body)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"error": err.Error(),
//		})
//		h.log.Error("failed to bind json", l.Error(err))
//		return
//	}
//
//	id, err := uuid.NewV4()
//	if err != nil {
//		h.log.Error("failed while generating uuid", l.Error(err))
//	}
//
//	body.Id = id.String()
//
//	response, err := postgres.NewBehaviorRepo(h.db).CreateBehavior(body)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"error": err.Error(),
//		})
//		h.log.Error("failed to create behavior", l.Error(err))
//		return
//	}
//
//	c.JSON(http.StatusCreated, response)
//}
//
//// GetBehavior ...
//// @Summary GetBehavior
//// @Description This API for getting behavior detail
//// @Tags behavior
//// @Accept json
//// @Produce json
//// @Param id path string true "BehaviorId"
//// @Success 200 {object} structs.BehaviorStruct
//// @Failure 400 {object} structs.StandardErrorModel
//// @Failure 500 {object} structs.StandardErrorModel
//// @Router /v1/behaviors/{id} [get]
//func (h *handlerV1) GetBehavior(c *gin.Context) {
//	guid := c.Param("id")
//
//	response, err := postgres.NewBehaviorRepo(h.db).GetBehavior(guid)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"error": err.Error(),
//		})
//		h.log.Error("failed to get behavior", l.Error(err))
//		return
//	}
//
//	c.JSON(http.StatusOK, response)
//}
//
//// GetListBehaviors ...
//// @Summary ListBehaviors
//// @Description This API for getting list of behaviors
//// @Tags behavior
//// @Accept json
//// @Produce json
//// @Param page query string false "Page"
//// @Param limit query string false "Limit"
//// @Success 200 {object} []structs.BehaviorStruct
//// @Failure 400 {object} structs.StandardErrorModel
//// @Failure 500 {object} structs.StandardErrorModel
//// @Router /v1/behaviors [get]
//func (h *handlerV1) GetListBehaviors(c *gin.Context) {
//	queryParams := c.Request.URL.Query()
//
//	params, errStr := utils.ParseQueryParams(queryParams)
//	if errStr != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"err": errStr[0],
//		})
//		h.log.Error("failed to parse query params json" + errStr[0])
//		return
//	}
//
//	response, count, err := postgres.NewBehaviorRepo(h.db).GetListBehaviors(int(params.Page), int(params.Limit))
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"error": err.Error(),
//		})
//		h.log.Error("failed to list behaviors", l.Error(err))
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"behaviors": response,
//		"count":     count,
//	})
//}
//
//// UpdateBehavior ...
//// @Summary UpdateBehavior
//// @Description This API for updating behavior
//// @Tags behavior
//// @Accept json
//// @Produce json
//// @Param id path string true "BehaivorId"
//// @Param User request body structs.BehaviorStruct true "BehaviorUpdateRequest"
//// @Success 200 {object} structs.BehaviorStruct
//// @Failure 400 {object} structs.StandardErrorModel
//// @Failure 500 {object} structs.StandardErrorModel
//// @Router /v1/behaviors/{id} [put]
//func (h *handlerV1) UpdateBehavior(c *gin.Context) {
//	var body structs.BehaviorStruct
//
//	err := c.ShouldBindJSON(&body)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"error": err.Error(),
//		})
//		h.log.Error("failed to bind json", l.Error(err))
//		return
//	}
//	body.Id = c.Param("id")
//
//	response, err := postgres.NewBehaviorRepo(h.db).UpdateBehavior(body)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"error": err.Error(),
//		})
//		h.log.Error("failed to update behavior", l.Error(err))
//		return
//	}
//
//	c.JSON(http.StatusOK, response)
//}
//
//// DeleteBehavior ...
//// @Summary DeleteBehavior
//// @Description This API for deleting the behavior
//// @Tags behavior
//// @Accept json
//// @Produce json
//// @Param id path string true "BehaviorId"
//// @Success 200
//// @Failure 400 {object} structs.StandardErrorModel
//// @Failure 500 {object} structs.StandardErrorModel
//// @Router /v1/behaviors/{id} [delete]
//func (h *handlerV1) DeleteBehavior(c *gin.Context) {
//	guid := c.Param("id")
//
//	err := postgres.NewBehaviorRepo(h.db).DeleteBehavior(guid)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"error": err.Error(),
//		})
//		h.log.Error("failed to delete behavior", l.Error(err))
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"response": "Deleted",
//	})
//}
