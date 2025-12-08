package controller

import (
	"miservicegolang/core/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UploadController struct {
	UploadUC *usecase.UploadUsecase
}

func NewUploadController(u *usecase.UploadUsecase) *UploadController {
	return &UploadController{
		UploadUC: u,
	}
}

func (c *UploadController) UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid file"})
		return
	}

	openedFile, nil := file.Open()
	err, result := c.UploadUC.Execute(openedFile)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, result)

}
