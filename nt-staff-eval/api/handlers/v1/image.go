package v1

import (
	"net/http"
	"strconv"

	l "github.com/Hatsker01/nt-staff-eval/pkg/logger"
	"github.com/Hatsker01/nt-staff-eval/pkg/structs"
	"github.com/Hatsker01/nt-staff-eval/storage/postgres"
	"github.com/gin-gonic/gin"
)

func (h *handlerV1) UploadUserImage(c *gin.Context) {
	file, err := c.FormFile("file")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
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

	newFileName := response.SpecId + response.FirstName + response.LastName + file.Filename[len(file.Filename)-4:]

	if err := c.SaveUploadedFile(file, h.cfg.PathUserImage+newFileName); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file",
		})
		return
	}

	err = postgres.NewImageRepo(h.db).LoadImage(structs.TypeUserImage, h.cfg.PathUserImage+newFileName, guid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to lead user image", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "success",
		"message": "file uploaded",
	})
}

func (h *handlerV1) UploadSectionImage(c *gin.Context) {
	file, err := c.FormFile("file")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}

	guid := c.Param("id")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.log.Error("Failed to parse string to int", l.Error(err))
	}
	response, err := postgres.NewSectionRepo(h.db).GetSection(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get section", l.Error(err))
		return
	}

	newFileName := response.Name + file.Filename[len(file.Filename)-4:]

	if err := c.SaveUploadedFile(file, h.cfg.PathSectionImage+newFileName); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file",
		})
		return
	}

	err = postgres.NewImageRepo(h.db).LoadImage(structs.TypeSectionImage, h.cfg.PathSectionImage+newFileName, guid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to lead section image", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"type":    "success",
		"message": "file uploaded",
	})
}
