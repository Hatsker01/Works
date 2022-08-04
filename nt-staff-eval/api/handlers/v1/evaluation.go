package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"

	l "github.com/Hatsker01/nt-staff-eval/pkg/logger"
	"github.com/Hatsker01/nt-staff-eval/pkg/structs"
	"github.com/Hatsker01/nt-staff-eval/storage/postgres"
)

// CreateEvaluation ...
// @Summary CreateEvaluation
// @Description This API for creating a new evaluation
// @Tags evaluation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param evaluation request body structs.UpdateEvaluationReq true "evaluationCreateRequest"
// @Success 200 {object} structs.UpdateEvaluation
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/evaluations/ [post]
func (h *handlerV1) CreateEvaluation(c *gin.Context) {
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
	var body structs.EvaluationStruct

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

	response, err := postgres.NewEvaluationRepo(h.db).CreateEvaluation(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create evaluation", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetEvaluation ...
// @Summary GetEvaluation
// @Description This API for getting evaluation detail
// @Tags evaluation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "EvaluationId"
// @Success 200 {object} structs.EvaluationStruct
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/evaluations/{id} [get]
func (h *handlerV1) GetEvaluation(c *gin.Context) {
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

	response, err := postgres.NewEvaluationRepo(h.db).GetEvaluation(guid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get evaluation", l.Error(err))
		return
	}
	if response.Section.Id != 0 {
		respSect, err := postgres.NewSectionRepo(h.db).GetSection(response.Section.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to get Section by Id for evaluation", l.Error(err))
			return
		}
		response.Section = respSect
	}
	c.JSON(http.StatusOK, response)
}

// GetListEvaluations ...
// @Summary ListEvaluations
// @Description This API for getting list of evaluations
// @Tags evaluation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} []structs.EvaluationStruct
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/evaluations [get]
func (h *handlerV1) GetListEvaluations(c *gin.Context) {
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

	response, count, err := postgres.NewEvaluationRepo(h.db).GetListEvaluations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to list evaluations", l.Error(err))
		return
	}

	for i, evaluation := range response {
		if evaluation.Section.Id != 0 {
			respSect, err := postgres.NewSectionRepo(h.db).GetSection(evaluation.Section.Id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				h.log.Error("failed to get Section by Id for evaluation", l.Error(err))
				return
			}
			response[i].Section = respSect
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"evaluations": response,
		"count":       count,
	})
}

// UpdateEvaluation ...
// @Summary UpdateEvaluation
// @Description This API for updating evaluation
// @Tags evaluation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "EvaluationId"
// @Param Evaluation request body structs.UpdateEvaluationReq true "EvaluationUpdateRequest"
// @Success 200 {object} structs.EvaluationStruct
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/evaluations/{id} [put]
func (h *handlerV1) UpdateEvaluation(c *gin.Context) {
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
	var body structs.UpdateEvaluation

	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	body.Id = c.Param("id")

	response, err := postgres.NewEvaluationRepo(h.db).UpdateEvaluation(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update evaluation", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteEvaluation ...
// @Summary DeleteEvaluation
// @Description This API for deleting the evaluation
// @Tags evaluation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "EvaluationId"
// @Success 200
// @Failure 400 {object} structs.StandardErrorModel
// @Failure 500 {object} structs.StandardErrorModel
// @Router /v1/evaluations/{id} [delete]
func (h *handlerV1) DeleteEvaluation(c *gin.Context) {
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

	err = postgres.NewEvaluationRepo(h.db).DeleteEvaluation(guid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete evaluation", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "Deleted",
	})
}
