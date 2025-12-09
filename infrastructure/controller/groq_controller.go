package controller

import (
	"miservicegolang/core/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AiController struct {
	usecase *usecase.AiUsecase
}

func NewAiController(a *usecase.AiUsecase) *AiController {
	return &AiController{
		usecase: a}
}

func (a *AiController) Generate(c *gin.Context) {
	var body struct {
		Prompt string `json:"prompt"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid Json body",
		})
		return
	}

	text, log := a.usecase.Generate(body.Prompt)
	if log.Error {
		c.JSON(http.StatusInternalServerError, log.Body)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": text,
	})
}
