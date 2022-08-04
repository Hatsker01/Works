package v1

//// LoginAdmin
//// @Summary LoginAdmin
//// @Tags admin
//// @Accept json
//// @Produce json
//// @Param loginData body structs.AdminStruct true "login data"
//// @Success 200 {object} structs.AdminStruct
//// @Failure 400 {object} structs.StandardErrorModel
//// @Router /v1/admin [post]
//func (h *handlerV1) LoginAdmin(c *gin.Context) {
//	var admin structs.AdminLogin
//	err := c.ShouldBindJSON(&admin)
//	if err != nil {
//		h.log.Error("error while ShouldBindJSON into LoginAdmin", l.Error(err))
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"error": "error while ShouldBindJSON into LoginAdmin" + err.Error(),
//		})
//		return
//	}
//
//	adminToken, err := postgres.NewAdminRepo(h.db).Login(admin)
//	if err != nil {
//		h.log.Error("error while logging into LoginAdmin", l.Error(err))
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"error": "error while logging into LoginAdmin" + err.Error(),
//		})
//		return
//	}
//
//	//err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(lUser.Password))
//	//if err != nil {
//	//	h.log.Error("error while hashing password", logger.Error(err))
//	//	c.JSON(http.StatusInternalServerError, models.Error{
//	//		Message: "error while hashing password " + err.Error(),
//	//	})
//	//	return
//	//}
//
//	h.jwtHandler.Sub = strconv.Itoa(adminToken.ID)
//	h.jwtHandler.Role = "admin"
//	h.jwtHandler.Aud = []string{"admin_profile"}
//
//	access, refresh, err := h.jwtHandler.GenerateAuthJWT()
//	if err != nil {
//		h.log.Error("error while generating tokens at LoginAdmin", l.Error(err))
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"error": "error while generating tokens at LoginAdmin" + err.Error(),
//		})
//		return
//	}
//
//	adminToken.RefreshToken = refresh
//	adminToken.AccessToken = access
//
//	err = postgres.NewAdminRepo(h.db).Update(adminToken)
//	if err != nil {
//		h.log.Error("error while updating tokens at LoginAdmin", l.Error(err))
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"error": "error while updating tokens LoginAdmin" + err.Error(),
//		})
//		return
//	}
//	c.JSON(http.StatusOK, adminToken)
//}
