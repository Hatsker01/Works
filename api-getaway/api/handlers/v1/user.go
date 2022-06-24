package v1

import (
	"context"
	"net/http"
	"time"

	pb "github.com/Hatsker01/Works/api-getaway/genproto"
	l "github.com/Hatsker01/Works/api-getaway/pkg/logger"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// CreateUser creates user
// route /v1/users [post]
func (h *handlerV1) Create(c *gin.Context) {
	var (
		body        pb.Post
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	//fmt.Println(&body)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.PostService().Create(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create user", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetUser gets user by id
// route /v1/users/{id} [get]
func (h *handlerV1) Update(c *gin.Context) {
	var (
		body        pb.Post
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.PostService().Update(
		ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *handlerV1) Delete(c *gin.Context) {
	guid := c.Param("id")
	//queryParams:=c.Request.URL.Query()
	//	params:=utils.ParseQueryParams2(queryParams)
	// if err!=nil{
	// 	c.JSON(http.StatusBadRequest,gin.H{
	// 		"error":err[0],
	// 	})
	// }
	//fmt.Println(guid)

	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	responce, err := h.serviceManager.PostService().Delete(ctx, &pb.DeleteRequest{
		Is: guid,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to list users", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, responce)

}
func (h *handlerV1) GetPostbyId(c *gin.Context) {
	guid := c.Param("id")

	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames=true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response,err:=h.serviceManager.PostService().GetPostbyId(ctx,&pb.GetPostRequest{
		Id: guid,
	})

	if err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to list users", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, response)
}
func (h *handlerV1) GetAll(c *gin.Context){
	ctx,cancel:=context.WithTimeout(context.Background(),time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	responce,err:=h.serviceManager.PostService().GetAll(ctx,&pb.Empty{})
	if err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get all posts", l.Error(err))
		return
	}
	c.JSON(http.StatusOK,responce)
}

// // ListUsers returns list of users
// // route /v1/users/ [get]
// func (h *handlerV1) ListUsers(c *gin.Context) {
// 	queryParams := c.Request.URL.Query()

// 	params, errStr := utils.ParseQueryParams(queryParams)
// 	if errStr != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": errStr[0],
// 		})
// 		h.log.Error("failed to parse query params json" + errStr[0])
// 		return
// 	}

// 	var jspbMarshal protojson.MarshalOptions
// 	jspbMarshal.UseProtoNames = true

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
// 	defer cancel()

// 	response, err := h.serviceManager.UserService().List(
// 		ctx, &pb.ListReq{
// 			Limit: params.Limit,
// 			Page:  params.Page,
// 		})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed to list users", l.Error(err))
// 		return
// 	}

// 	c.JSON(http.StatusOK, response)
// }

// // UpdateUser updates user by id
// // route /v1/users/{id} [put]
// func (h *handlerV1) UpdateUser(c *gin.Context) {
// 	var (
// 		body        pb.User
// 		jspbMarshal protojson.MarshalOptions
// 	)
// 	jspbMarshal.UseProtoNames = true

// 	err := c.ShouldBindJSON(&body)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed to bind json", l.Error(err))
// 		return
// 	}
// 	body.Id = c.Param("id")

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
// 	defer cancel()

// 	response, err := h.serviceManager.UserService().Update(ctx, &body)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed to update user", l.Error(err))
// 		return
// 	}

// 	c.JSON(http.StatusOK, response)
// }

// // DeleteUser deletes user by id
// // route /v1/users/{id} [delete]
// func (h *handlerV1) DeleteUser(c *gin.Context) {
// 	var jspbMarshal protojson.MarshalOptions
// 	jspbMarshal.UseProtoNames = true

// 	guid := c.Param("id")
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
// 	defer cancel()

// 	response, err := h.serviceManager.UserService().Delete(
// 		ctx, &pb.ByIdReq{
// 			Id: guid,
// 		})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed to delete user", l.Error(err))
// 		return
// 	}

// 	c.JSON(http.StatusOK, response)
// }
