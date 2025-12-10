package controller

import (
	"miservicegolang/core/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProgressController struct {
	usecase usecase.GrowthUsecase
}

func NewProgressController(g usecase.GrowthUsecase) *ProgressController {
	return &ProgressController{usecase: g}
}

func (p *ProgressController) Rank(c *gin.Context) {
	log := p.usecase.Rank()

	if log.Error {
		c.JSON(http.StatusInternalServerError, log.Body)
		return
	}

	c.JSON(http.StatusOK, log.Body)
}
