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
	return &AiController{usecase: a}
}

func (a *AiController) Generate(c *gin.Context) {
	result, log := a.usecase.Generate(
		"Gere um desafio simples envolvendo programação. Não gere código.",
	)

	if log.Error {
		c.JSON(http.StatusInternalServerError, log.Body)
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

func (a *AiController) Verify(c *gin.Context) {
	var body struct {
		Code string `json:"code"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Invalid JSON"})
		return
	}

	result, log := a.usecase.Verify("693849deb204cc30c0240bdf", body.Code)

	if log.Error {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "log": log.Body})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "data": result})
}
