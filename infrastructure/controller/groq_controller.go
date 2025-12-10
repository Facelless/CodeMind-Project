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
		`
		Gere um desafio de programação totalmente original que exija raciocínio logico, gere um desafio extramamente facil que nao precise apis, banco de dados.
		`,
	)

	if log.Error {
		c.JSON(http.StatusInternalServerError, log.Body)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        result.ID.Hex(),
		"challenge": result.Answer,
		"completed": result.Completed,
	})
}

func (a *AiController) Verify(c *gin.Context) {
	var body struct {
		Code string `json:"code"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid JSON body",
		})
		return
	}
	data, log := a.usecase.Verify("6938bc26c1628de6cc606fd2", body.Code)
	if log.Error {
		c.JSON(http.StatusBadRequest, log.Body)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        data.ID.Hex(),
		"verify":    data.Verify,
		"completed": data.Completed,
	})
}
