package controller

import (
	"miservicegolang/core/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AiController struct {
	usecase *usecase.AiUsecase
}

func NewAiController(u *usecase.AiUsecase) *AiController {
	return &AiController{
		usecase: u,
	}
}

func (c *AiController) Generate(ctx *gin.Context) {
	var body struct {
		Prompt string `json:"prompt"`
	}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid Json body",
		})
		return
	}

	text, log := c.usecase.Generate(body.Prompt)
	if log.Error {
		ctx.JSON(http.StatusInternalServerError, log.Body)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": text,
	})
}
