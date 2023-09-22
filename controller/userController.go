package controller

import (
	"Patronus/model"
	"Patronus/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return UserController{userService}
}

func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*model.UserDBResponseModel)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": model.FilteredUserResponse(currentUser)}})
}
