package controller

import (
	"miservicegolang/core/domain/user"
	"miservicegolang/core/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	usecase *usecase.UserUsecase
}

func NewUserController(u *usecase.UserUsecase) *UserController {
	return &UserController{usecase: u}
}

func (h *UserController) Register(c *gin.Context) {
	var u user.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := h.usecase.Register(u); err.Error {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registered user"})
}
